package main

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strings"
)

const (
    appName = "Patchy88 - Mugen Senshi Valis II PC-88"
    appVersion = "1.0.0"
    manifestName = "Patchy88-Valis2.json"
)

type Manifest struct {
    Schema string `json:"schema"`
    Version int `json:"version"`
    Game string `json:"game"`
    Language string `json:"language"`
    PatcherVersion string `json:"patcher_version"`
    Targets []TargetMeta `json:"targets"`
}

type TargetMeta struct {
    ID string `json:"id"`
    DisplayName string `json:"display_name"`
    Kind string `json:"kind"`
    ExpectedSize int `json:"expected_size"`
    KnownOriginalSHA256 string `json:"known_original_sha256"`
    KnownPatchedSHA256 string `json:"known_patched_sha256"`
    Patch string `json:"patch"`
    PatchSHA256 string `json:"patch_sha256"`
    Records []RecordMeta `json:"records"`
}

type RecordMeta struct {
    Index int `json:"index"`
    Offset int `json:"offset"`
    Length int `json:"length"`
    BeforeSHA256 string `json:"before_sha256"`
    AfterSHA256 string `json:"after_sha256"`
}

type IPSRecord struct {
    Offset int
    Data []byte
}

type RuntimeTarget struct {
    Meta TargetMeta
    PatchRecords []IPSRecord
}

type FileState string
const (
    StateOriginal FileState = "ORIGINAL"
    StatePatched FileState = "ALREADY_PATCHED"
    StatePartial FileState = "PARTIAL"
    StateIncompatible FileState = "INCOMPATIBLE"
)

type Match struct {
    Target *RuntimeTarget
    Path string
    State FileState
    FullSHA256 string
}

type ScanResult struct {
    Folder string
    ByTarget map[string][]Match
    Ignored []string
}

type Prepared struct {
    Match Match
    Output string
    Backup string
    Temp string
}

func sha256Hex(b []byte) string {
    s := sha256.Sum256(b)
    return hex.EncodeToString(s[:])
}

func appDir() string {
    if v := os.Getenv("PATCHY88_APP_DIR"); v != "" { return v }
    exe, err := os.Executable()
    if err == nil { return filepath.Dir(exe) }
    wd, _ := os.Getwd()
    return wd
}

func loadManifest() (*Manifest, error) {
    p := filepath.Join(appDir(), manifestName)
    b, err := os.ReadFile(p)
    if err != nil { return nil, fmt.Errorf("매니페스트를 열 수 없습니다: %s: %w", p, err) }
    var m Manifest
    if err := json.Unmarshal(b, &m); err != nil { return nil, fmt.Errorf("매니페스트 JSON 오류: %w", err) }
    if m.Schema != "patchy88.pc88.valis2.ips_manifest" { return nil, errors.New("지원하지 않는 매니페스트 형식입니다") }
    if len(m.Targets) != 8 { return nil, fmt.Errorf("대상 수가 8개가 아닙니다: %d", len(m.Targets)) }
    return &m, nil
}

func parseIPS(blob []byte) ([]IPSRecord, error) {
    if len(blob) < 8 || string(blob[:5]) != "PATCH" { return nil, errors.New("IPS 헤더(PATCH)가 없습니다") }
    c := 5
    out := make([]IPSRecord, 0)
    for {
        if c+3 > len(blob) { return nil, errors.New("IPS가 EOF 전에 잘렸습니다") }
        if string(blob[c:c+3]) == "EOF" {
            c += 3
            if c != len(blob) { return nil, errors.New("IPS EOF 뒤 확장 데이터는 허용하지 않습니다") }
            if len(out) == 0 { return nil, errors.New("IPS 레코드가 없습니다") }
            return out, nil
        }
        if c+5 > len(blob) { return nil, errors.New("IPS 레코드 헤더가 잘렸습니다") }
        off := int(blob[c])<<16 | int(blob[c+1])<<8 | int(blob[c+2])
        n := int(blob[c+3])<<8 | int(blob[c+4])
        c += 5
        var data []byte
        if n != 0 {
            if c+n > len(blob) { return nil, fmt.Errorf("IPS 데이터가 잘렸습니다: 0x%06X", off) }
            data = append([]byte(nil), blob[c:c+n]...)
            c += n
        } else {
            if c+3 > len(blob) { return nil, fmt.Errorf("IPS RLE가 잘렸습니다: 0x%06X", off) }
            repeat := int(blob[c])<<8 | int(blob[c+1])
            if repeat == 0 { return nil, fmt.Errorf("IPS RLE 길이가 0입니다: 0x%06X", off) }
            value := blob[c+2]
            c += 3
            data = make([]byte, repeat)
            for i := range data { data[i] = value }
        }
        out = append(out, IPSRecord{Offset: off, Data: data})
    }
}

func loadRuntimeTargets() ([]RuntimeTarget, error) {
    m, err := loadManifest()
    if err != nil { return nil, err }
    out := make([]RuntimeTarget, 0, len(m.Targets))
    for _, tm := range m.Targets {
        p := filepath.Join(appDir(), filepath.FromSlash(tm.Patch))
        blob, err := os.ReadFile(p)
        if err != nil { return nil, fmt.Errorf("%s IPS를 열 수 없습니다: %w", tm.DisplayName, err) }
        if sha256Hex(blob) != strings.ToLower(tm.PatchSHA256) { return nil, fmt.Errorf("%s IPS SHA-256 불일치", tm.DisplayName) }
        recs, err := parseIPS(blob)
        if err != nil { return nil, fmt.Errorf("%s IPS 오류: %w", tm.DisplayName, err) }
        if len(recs) != len(tm.Records) { return nil, fmt.Errorf("%s IPS 레코드 수 불일치", tm.DisplayName) }
        for i := range recs {
            meta := tm.Records[i]
            if recs[i].Offset != meta.Offset || len(recs[i].Data) != meta.Length { return nil, fmt.Errorf("%s IPS 구조 불일치 #%d", tm.DisplayName, meta.Index) }
            if sha256Hex(recs[i].Data) != strings.ToLower(meta.AfterSHA256) { return nil, fmt.Errorf("%s IPS 데이터 불일치 #%d", tm.DisplayName, meta.Index) }
        }
        out = append(out, RuntimeTarget{Meta: tm, PatchRecords: recs})
    }
    return out, nil
}

func d88StructureOK(data []byte) error {
    if len(data) < 0x2B0 { return errors.New("D88 헤더보다 파일이 짧습니다") }
    declared := int(data[0x1C]) | int(data[0x1D])<<8 | int(data[0x1E])<<16 | int(data[0x1F])<<24
    if declared != len(data) { return fmt.Errorf("D88 헤더 크기 %d / 실제 %d", declared, len(data)) }
    pointers := make([]int, 0, 164)
    for i:=0; i<164; i++ {
        o := 0x20+i*4
        p := int(data[o]) | int(data[o+1])<<8 | int(data[o+2])<<16 | int(data[o+3])<<24
        if p != 0 { pointers = append(pointers, p) }
    }
    if len(pointers) != 80 { return fmt.Errorf("D88 트랙 포인터 수가 80개가 아닙니다: %d", len(pointers)) }
    if pointers[0] < 0x2B0 { return errors.New("첫 트랙 포인터가 비정상입니다") }
    prev := 0
    for _, p := range pointers {
        if p <= prev || p >= len(data) { return errors.New("D88 트랙 포인터 순서/범위가 비정상입니다") }
        prev = p
    }
    return nil
}

func classifyBytes(data []byte, t *RuntimeTarget) FileState {
    if len(data) != t.Meta.ExpectedSize { return StateIncompatible }
    before, after, unknown := 0,0,0
    for _, r := range t.Meta.Records {
        end := r.Offset + r.Length
        if r.Offset < 0 || end > len(data) { unknown++; continue }
        h := sha256Hex(data[r.Offset:end])
        if h == strings.ToLower(r.BeforeSHA256) { before++
        } else if h == strings.ToLower(r.AfterSHA256) { after++
        } else { unknown++ }
    }
    if unknown != 0 { return StateIncompatible }
    if before > 0 && after == 0 { return StateOriginal }
    if after > 0 && before == 0 { return StatePatched }
    if before > 0 && after > 0 { return StatePartial }
    return StateIncompatible
}

func scanFolder(folder string, targets []RuntimeTarget) (*ScanResult, error) {
    info, err := os.Stat(folder)
    if err != nil || !info.IsDir() { return nil, fmt.Errorf("폴더를 열 수 없습니다: %s", folder) }
    entries, err := os.ReadDir(folder)
    if err != nil { return nil, err }
    result := &ScanResult{Folder: folder, ByTarget: map[string][]Match{}}
    for i := range targets { result.ByTarget[targets[i].Meta.ID] = []Match{} }
    for _, e := range entries {
        if e.IsDir() { continue }
        ext := strings.ToLower(filepath.Ext(e.Name()))
        if ext != ".d88" && ext != ".rom" { continue }
        p := filepath.Join(folder, e.Name())
        data, err := os.ReadFile(p)
        if err != nil { result.Ignored = append(result.Ignored, e.Name()+" (읽기 실패)"); continue }
        matched := 0
        for i := range targets {
            t := &targets[i]
            if (t.Meta.Kind == "d88" && ext != ".d88") || (t.Meta.Kind == "rom" && ext != ".rom") { continue }
            state := classifyBytes(data, t)
            if state == StateIncompatible { continue }
            if t.Meta.Kind == "d88" {
                if err := d88StructureOK(data); err != nil { continue }
            }
            result.ByTarget[t.Meta.ID] = append(result.ByTarget[t.Meta.ID], Match{Target:t, Path:p, State:state, FullSHA256:sha256Hex(data)})
            matched++
        }
        if matched == 0 { result.Ignored = append(result.Ignored, e.Name()) }
    }
    return result, nil
}

func stateKorean(s FileState) string {
    switch s {
    case StateOriginal: return "정상 원본"
    case StatePatched: return "이미 패치됨"
    case StatePartial: return "부분 패치"
    default: return "호환되지 않음"
    }
}

func preflight(scan *ScanResult, targets []RuntimeTarget) ([]Match, []string) {
    selected := make([]Match,0,len(targets))
    problems := []string{}
    for i := range targets {
        t := &targets[i]
        list := scan.ByTarget[t.Meta.ID]
        if len(list)==0 {
            problems = append(problems, fmt.Sprintf("%s: 찾지 못함", t.Meta.DisplayName))
            continue
        }
        if len(list)>1 {
            names := make([]string,0,len(list)); for _,m := range list { names=append(names,filepath.Base(m.Path)) }
            problems = append(problems, fmt.Sprintf("%s: 중복/모호함 (%s)", t.Meta.DisplayName, strings.Join(names,", ")))
            continue
        }
        m := list[0]
        if m.State==StatePartial {
            problems = append(problems, fmt.Sprintf("%s: 부분 패치 상태 (%s)", t.Meta.DisplayName, filepath.Base(m.Path)))
            continue
        }
        selected = append(selected,m)
    }
    return selected,problems
}

func outputPathFor(path string) string {
    ext := filepath.Ext(path)
    base := strings.TrimSuffix(path, ext)
    return base+"(K)"+ext
}

func uniqueBackupPath(path string) string {
    c := path+".bak"
    if _,err:=os.Stat(c); os.IsNotExist(err) { return c }
    for n:=1;;n++ {
        c=fmt.Sprintf("%s.%d.bak",path,n)
        if _,err:=os.Stat(c); os.IsNotExist(err) { return c }
    }
}

func applyRecords(base []byte, recs []IPSRecord) ([]byte,error) {
    out:=append([]byte(nil),base...)
    for _,r:=range recs {
        end:=r.Offset+len(r.Data)
        if r.Offset<0 || end>len(out) { return nil,fmt.Errorf("IPS가 파일 끝을 넘습니다: 0x%06X",r.Offset) }
        copy(out[r.Offset:end],r.Data)
    }
    return out,nil
}

func writeTempVerified(m Match) (string,error) {
    src,err:=os.ReadFile(m.Path); if err!=nil{return "",err}
    patched,err:=applyRecords(src,m.Target.PatchRecords); if err!=nil{return "",err}
    if classifyBytes(patched,m.Target)!=StatePatched{return "",errors.New("메모리상 패치 결과 사후검증 실패")}
    if m.Target.Meta.Kind=="d88" { if err:=d88StructureOK(patched);err!=nil{return "",fmt.Errorf("패치 결과 D88 구조검증 실패: %w",err)} }
    f,err:=os.CreateTemp(filepath.Dir(m.Path),".patchy88-valis2-*.tmp"); if err!=nil{return "",err}
    tmp:=f.Name(); ok:=false
    defer func(){ if !ok { f.Close(); os.Remove(tmp) } }()
    if _,err=f.Write(patched);err!=nil{return "",err}; if err=f.Sync();err!=nil{return "",err}; if err=f.Close();err!=nil{return "",err}
    check,err:=os.ReadFile(tmp); if err!=nil{return "",err}
    if classifyBytes(check,m.Target)!=StatePatched{return "",errors.New("임시파일 사후검증 실패")}
    ok=true; return tmp,nil
}

func batchPatch(folder string, targets []RuntimeTarget) ([]string,error) {
    scan,err:=scanFolder(folder,targets); if err!=nil{return nil,err}
    selected,problems:=preflight(scan,targets)
    if len(problems)>0{return nil,fmt.Errorf("사전검증 실패:\n%s",strings.Join(problems,"\n"))}
    if len(selected)!=len(targets){return nil,errors.New("8개 대상이 모두 확인되지 않았습니다")}

    prepared:=[]Prepared{}
    logs:=[]string{}
    for _,m:=range selected {
        if m.State==StatePatched { logs=append(logs,fmt.Sprintf("%s: 이미 패치됨 - %s",m.Target.Meta.DisplayName,filepath.Base(m.Path))); continue }
        if m.State!=StateOriginal{return nil,fmt.Errorf("%s 상태가 안전하지 않습니다: %s",m.Target.Meta.DisplayName,m.State)}
        out:=outputPathFor(m.Path)
        if _,err:=os.Stat(out); err==nil{return nil,fmt.Errorf("출력 파일이 이미 존재합니다: %s",filepath.Base(out))}
        tmp,err:=writeTempVerified(m); if err!=nil {
            for _,p:=range prepared{os.Remove(p.Temp)}
            return nil,fmt.Errorf("%s 임시 패치 생성 실패: %w",m.Target.Meta.DisplayName,err)
        }
        prepared=append(prepared,Prepared{Match:m,Output:out,Backup:uniqueBackupPath(m.Path),Temp:tmp})
    }
    if len(prepared)==0 { logs=append(logs,"모든 대상이 이미 패치되어 있습니다."); return logs,nil }

    committed:=[]Prepared{}
    rollback:=func(){
        for i:=len(committed)-1;i>=0;i--{
            p:=committed[i]
            os.Remove(p.Output)
            if _,err:=os.Stat(p.Backup);err==nil{_ = os.Rename(p.Backup,p.Match.Path)}
        }
        for _,p:=range prepared{_ = os.Remove(p.Temp)}
    }
    for _,p:=range prepared {
        if err:=os.Rename(p.Match.Path,p.Backup);err!=nil{rollback();return nil,fmt.Errorf("%s 백업 생성 실패: %w",p.Match.Target.Meta.DisplayName,err)}
        if err:=os.Rename(p.Temp,p.Output);err!=nil{
            _=os.Rename(p.Backup,p.Match.Path)
            rollback();return nil,fmt.Errorf("%s 출력 확정 실패: %w",p.Match.Target.Meta.DisplayName,err)
        }
        committed=append(committed,p)
    }
    // Final verification of every committed output before declaring success.
    for _,p:=range committed {
        b,err:=os.ReadFile(p.Output); if err!=nil{rollback();return nil,fmt.Errorf("%s 최종 파일 읽기 실패: %w",p.Match.Target.Meta.DisplayName,err)}
        if classifyBytes(b,p.Match.Target)!=StatePatched{rollback();return nil,fmt.Errorf("%s 최종 사후검증 실패",p.Match.Target.Meta.DisplayName)}
        if p.Match.Target.Meta.Kind=="d88" { if err:=d88StructureOK(b);err!=nil{rollback();return nil,fmt.Errorf("%s 최종 D88 구조검증 실패: %w",p.Match.Target.Meta.DisplayName,err)} }
    }
    sort.Slice(committed,func(i,j int)bool{return committed[i].Match.Target.Meta.ID<committed[j].Match.Target.Meta.ID})
    for _,p:=range committed {
        logs=append(logs,fmt.Sprintf("%s: 완료 | 백업 %s | 출력 %s",p.Match.Target.Meta.DisplayName,filepath.Base(p.Backup),filepath.Base(p.Output)))
    }
    return logs,nil
}

func copyFile(src,dst string) error {
    in,err:=os.Open(src); if err!=nil{return err}; defer in.Close()
    out,err:=os.Create(dst); if err!=nil{return err}
    _,cpErr:=io.Copy(out,in); syncErr:=out.Sync(); closeErr:=out.Close()
    if cpErr!=nil{return cpErr}; if syncErr!=nil{return syncErr}; return closeErr
}
