package main

import (
    "fmt"
    "hash/adler32"
    "io"
    "os"
    "strings"
)

// The official v1.01 release checks full MD5 and SHA-256 of all outputs.
// Adler-32 fallback is retained only for deliberately incomplete test fixtures.
func verifyTarget(path string, got fileHashes, target TargetDef) (string,error) {
    if len(target.MD5)==32 && len(target.SHA256)==64 {
        if target.Size>0 && got.Size!=target.Size {
            return "",fmt.Errorf("검증 대상 크기 불일치: got=%d expected=%d",got.Size,target.Size)
        }
        if !hashEqual(got,target.Hashes) {
            return "",fmt.Errorf("MD5/SHA-256 불일치: MD5=%s SHA-256=%s",got.MD5,got.SHA256)
        }
        return "MD5/SHA-256",nil
    }
    if target.Size<=0 || got.Size!=target.Size || len(target.WindowAdler32)==0 || len(target.WindowAdler32)!=len(target.WindowSizes) {
        return "",fmt.Errorf("IMG 출력 크기/윈도우 검증표 불일치: got=%d expected=%d",got.Size,target.Size)
    }
    f,err:=os.Open(path)
    if err!=nil {return "",err}
    defer f.Close()
    var total int64
    for i,size:=range target.WindowSizes {
        if size<=0 || size>8*1024*1024{return "",fmt.Errorf("IMG 윈도우 %d 크기 오류",i)}
        h:=adler32.New()
        n,err:=io.CopyN(h,f,size)
        if err!=nil || n!=size {return "",fmt.Errorf("IMG 윈도우 %d 읽기 실패: %v",i,err)}
        actual:=fmt.Sprintf("%08x",h.Sum32())
        if !strings.EqualFold(actual,target.WindowAdler32[i]) {return "",fmt.Errorf("IMG 윈도우 %d Adler32 불일치: got=%s expected=%s",i,actual,target.WindowAdler32[i])}
        total+=size
    }
    if total!=target.Size{return "",fmt.Errorf("IMG 윈도우 길이 합계 오류: %d",total)}
    var final [1]byte
    n,err:=f.Read(final[:])
    if n!=0 || err!=io.EOF{return "",fmt.Errorf("IMG 출력 끝에 예상외 데이터")}
    return fmt.Sprintf("%d개 윈도우 Adler32 및 크기",len(target.WindowAdler32)),nil
}
