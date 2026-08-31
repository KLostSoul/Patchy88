package main

import (
  "os"
  "path/filepath"
  "testing"
  "strings"
)

func TestManifestAndPatches(t *testing.T){
  targets,err:=loadRuntimeTargets(); if err!=nil{t.Fatal(err)}
  if len(targets)!=8{t.Fatalf("targets=%d",len(targets))}
}

func TestScanAndBatchPatch(t *testing.T){
  targets,err:=loadRuntimeTargets(); if err!=nil{t.Fatal(err)}
  srcRoot:=os.Getenv("PATCHY88_TEST_INPUT")
  if srcRoot==""{t.Skip("PATCHY88_TEST_INPUT not set")}
  dir:=t.TempDir()
  names:=[]string{
    "Mugen Senshi Valis II (Disk A)(1).d88","Mugen Senshi Valis II (Disk B).d88","Mugen Senshi Valis II (Disk C).d88","Mugen Senshi Valis II (Disk D).d88","Mugen Senshi Valis II (Disk E).d88","Mugen Senshi Valis II (Disk F).d88","Mugen Senshi Valis II (Disk G).d88","KANJI1(5).ROM",
  }
  for _,n:=range names{if err:=copyFile(filepath.Join(srcRoot,n),filepath.Join(dir,n));err!=nil{t.Fatal(err)}}
  scan,err:=scanFolder(dir,targets); if err!=nil{t.Fatal(err)}
  sel,problems:=preflight(scan,targets); if len(problems)!=0||len(sel)!=8{t.Fatalf("preflight problems=%v selected=%d",problems,len(sel))}
  logs,err:=batchPatch(dir,targets); if err!=nil{t.Fatal(err)}
  if len(logs)!=8{t.Fatalf("logs=%d",len(logs))}
  scan2,err:=scanFolder(dir,targets); if err!=nil{t.Fatal(err)}
  sel2,problems2:=preflight(scan2,targets); if len(problems2)!=0||len(sel2)!=8{t.Fatalf("post problems=%v selected=%d",problems2,len(sel2))}
  for _,m:=range sel2{if m.State!=StatePatched{t.Fatalf("%s state=%s",m.Target.Meta.ID,m.State)}}
  // All originals moved to same-folder .bak and outputs carry (K).
  for _,n:=range names{
    if _,err:=os.Stat(filepath.Join(dir,n+".bak"));err!=nil{t.Fatalf("missing backup for %s: %v",n,err)}
    ext:=filepath.Ext(n); out:=n[:len(n)-len(ext)]+"(K)"+ext
    if _,err:=os.Stat(filepath.Join(dir,out));err!=nil{t.Fatalf("missing output %s: %v",out,err)}
  }
}

func TestFilenameIndependentAndAtomicReject(t *testing.T){
  targets,err:=loadRuntimeTargets(); if err!=nil{t.Fatal(err)}
  srcRoot:=os.Getenv("PATCHY88_TEST_INPUT"); if srcRoot==""{t.Skip("PATCHY88_TEST_INPUT not set")}
  dir:=t.TempDir()
  srcNames:=[]string{
    "Mugen Senshi Valis II (Disk A)(1).d88","Mugen Senshi Valis II (Disk B).d88","Mugen Senshi Valis II (Disk C).d88","Mugen Senshi Valis II (Disk D).d88","Mugen Senshi Valis II (Disk E).d88","Mugen Senshi Valis II (Disk F).d88","Mugen Senshi Valis II (Disk G).d88","KANJI1(5).ROM",
  }
  dstNames:=[]string{"one.d88","two.d88","three.d88","four.d88","five.d88","six.d88","seven.d88","font.ROM"}
  for i,n:=range srcNames{if err:=copyFile(filepath.Join(srcRoot,n),filepath.Join(dir,dstNames[i]));err!=nil{t.Fatal(err)}}
  scan,err:=scanFolder(dir,targets); if err!=nil{t.Fatal(err)}
  sel,problems:=preflight(scan,targets); if len(problems)!=0||len(sel)!=8{t.Fatalf("filename-independent scan failed: %v %d",problems,len(sel))}
  // Turn only the first Disk A IPS record into its patched value -> PARTIAL.
  var diskA *RuntimeTarget
  for i:=range targets{if targets[i].Meta.ID=="disk_a"{diskA=&targets[i];break}}
  if diskA==nil{t.Fatal("disk_a missing")}
  p:=filepath.Join(dir,"one.d88"); b,err:=os.ReadFile(p);if err!=nil{t.Fatal(err)}
  r:=diskA.PatchRecords[0];copy(b[r.Offset:r.Offset+len(r.Data)],r.Data);if err:=os.WriteFile(p,b,0644);err!=nil{t.Fatal(err)}
  if _,err:=batchPatch(dir,targets);err==nil{t.Fatal("partial patch should have been rejected")}
  // No batch mutation is allowed after a preflight failure.
  entries,_:=os.ReadDir(dir)
  for _,e:=range entries{if strings.Contains(e.Name(),".bak")||strings.Contains(e.Name(),"(K)"){t.Fatalf("unexpected mutation after reject: %s",e.Name())}}
}
