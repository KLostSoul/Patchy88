//go:build windows

package main

import (
    "fmt"
    "path/filepath"
    "runtime"
    "strings"
    "syscall"
    "unsafe"
)

var (
    user32 = syscall.NewLazyDLL("user32.dll")
    kernel32 = syscall.NewLazyDLL("kernel32.dll")
    gdi32 = syscall.NewLazyDLL("gdi32.dll")
    shell32 = syscall.NewLazyDLL("shell32.dll")
    ole32 = syscall.NewLazyDLL("ole32.dll")

    pRegisterClassExW = user32.NewProc("RegisterClassExW")
    pCreateWindowExW = user32.NewProc("CreateWindowExW")
    pDefWindowProcW = user32.NewProc("DefWindowProcW")
    pShowWindow = user32.NewProc("ShowWindow")
    pUpdateWindow = user32.NewProc("UpdateWindow")
    pGetMessageW = user32.NewProc("GetMessageW")
    pTranslateMessage = user32.NewProc("TranslateMessage")
    pDispatchMessageW = user32.NewProc("DispatchMessageW")
    pPostQuitMessage = user32.NewProc("PostQuitMessage")
    pDestroyWindow = user32.NewProc("DestroyWindow")
    pSendMessageW = user32.NewProc("SendMessageW")
    pSetWindowTextW = user32.NewProc("SetWindowTextW")
    pMessageBoxW = user32.NewProc("MessageBoxW")
    pEnableWindow = user32.NewProc("EnableWindow")
    pMoveWindow = user32.NewProc("MoveWindow")
    pGetClientRect = user32.NewProc("GetClientRect")
    pLoadCursorW = user32.NewProc("LoadCursorW")
    pGetModuleHandleW = kernel32.NewProc("GetModuleHandleW")
    pCreateFontW = gdi32.NewProc("CreateFontW")
    pSHBrowseForFolderW = shell32.NewProc("SHBrowseForFolderW")
    pSHGetPathFromIDListW = shell32.NewProc("SHGetPathFromIDListW")
    pCoTaskMemFree = ole32.NewProc("CoTaskMemFree")
    pCoInitializeEx = ole32.NewProc("CoInitializeEx")
    pCoUninitialize = ole32.NewProc("CoUninitialize")
)

const (
    WM_DESTROY = 0x0002
    WM_SIZE = 0x0005
    WM_COMMAND = 0x0111
    WM_SETFONT = 0x0030
    EM_SETSEL = 0x00B1
    EM_REPLACESEL = 0x00C2

    WS_OVERLAPPEDWINDOW = 0x00CF0000
    WS_VISIBLE = 0x10000000
    WS_CHILD = 0x40000000
    WS_TABSTOP = 0x00010000
    WS_BORDER = 0x00800000
    WS_VSCROLL = 0x00200000
    ES_READONLY = 0x0800
    ES_MULTILINE = 0x0004
    ES_AUTOVSCROLL = 0x0040
    WS_EX_CLIENTEDGE = 0x00000200
    BS_DEFPUSHBUTTON = 0x00000001

    MB_OK = 0x00000000
    MB_ICONINFORMATION = 0x00000040
    MB_ICONERROR = 0x00000010
    MB_ICONQUESTION = 0x00000020
    MB_YESNO = 0x00000004
    IDYES = 6

    BIF_RETURNONLYFSDIRS = 0x0001
    BIF_NEWDIALOGSTYLE = 0x0040
    COINIT_APARTMENTTHREADED = 0x2

    ID_FOLDER = 1001
    ID_SCAN = 1002
    ID_PATCH = 1003
    ID_EXIT = 1004
)

type WNDCLASSEX struct {
    CbSize uint32
    Style uint32
    LpfnWndProc uintptr
    CbClsExtra int32
    CbWndExtra int32
    HInstance uintptr
    HIcon uintptr
    HCursor uintptr
    HbrBackground uintptr
    LpszMenuName *uint16
    LpszClassName *uint16
    HIconSm uintptr
}

type POINT struct { X, Y int32 }
type MSG struct {
    Hwnd uintptr
    Message uint32
    WParam uintptr
    LParam uintptr
    Time uint32
    Pt POINT
    LPrivate uint32
}
type RECT struct { Left, Top, Right, Bottom int32 }
type BROWSEINFO struct {
    HwndOwner uintptr
    PidlRoot uintptr
    PszDisplayName *uint16
    LpszTitle *uint16
    UlFlags uint32
    Lpfn uintptr
    LParam uintptr
    IImage int32
}

var (
    hwndMain uintptr
    hwndFolderEdit uintptr
    hwndFolderButton uintptr
    hwndScanButton uintptr
    hwndPatchButton uintptr
    hwndExitButton uintptr
    hwndLog uintptr
    hwndFooter uintptr
    normalFont uintptr
    titleFont uintptr
    runtimeTargets []RuntimeTarget
    currentFolder string
)

func u16(s string) *uint16 { return syscall.StringToUTF16Ptr(s) }

func setText(hwnd uintptr, s string) { pSetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(u16(s)))) }
func send(hwnd uintptr, msg uint32, w, l uintptr) uintptr { r,_,_:=pSendMessageW.Call(hwnd,uintptr(msg),w,l); return r }
func setFont(hwnd, font uintptr) { send(hwnd,WM_SETFONT,font,1) }
func enable(hwnd uintptr, on bool) { v:=uintptr(0); if on{v=1}; pEnableWindow.Call(hwnd,v) }

func appendLog(s string) {
    s = strings.ReplaceAll(s, "\n", "\r\n") + "\r\n"
    send(hwndLog, EM_SETSEL, ^uintptr(0), ^uintptr(0))
    send(hwndLog, EM_REPLACESEL, 0, uintptr(unsafe.Pointer(u16(s))))
}

func message(text string, flags uintptr) int {
    r,_,_:=pMessageBoxW.Call(hwndMain,uintptr(unsafe.Pointer(u16(text))),uintptr(unsafe.Pointer(u16("Patchy88 - 바리스 II"))),flags)
    return int(r)
}

func browseFolder() string {
    var display [260]uint16
    bi:=BROWSEINFO{HwndOwner:hwndMain,PszDisplayName:&display[0],LpszTitle:u16("바리스 II Disk A~G와 KANJI1 ROM이 있는 폴더를 선택하세요."),UlFlags:BIF_RETURNONLYFSDIRS|BIF_NEWDIALOGSTYLE}
    pidl,_,_:=pSHBrowseForFolderW.Call(uintptr(unsafe.Pointer(&bi)))
    if pidl==0{return ""}
    defer pCoTaskMemFree.Call(pidl)
    var path [32768]uint16
    ok,_,_:=pSHGetPathFromIDListW.Call(pidl,uintptr(unsafe.Pointer(&path[0])))
    if ok==0{return ""}
    return syscall.UTF16ToString(path[:])
}

func scanAndDisplay() bool {
    if currentFolder==""{message("먼저 폴더를 선택해 주세요.",MB_OK|MB_ICONINFORMATION);return false}
    scan,err:=scanFolder(currentFolder,runtimeTargets)
    if err!=nil{appendLog("[오류] "+err.Error());message(err.Error(),MB_OK|MB_ICONERROR);enable(hwndPatchButton,false);return false}
    appendLog("----------------------------------------")
    appendLog("폴더 검사: "+currentFolder)
    selected,problems:=preflight(scan,runtimeTargets)
    for i:=range runtimeTargets{
        t:=&runtimeTargets[i]
        list:=scan.ByTarget[t.Meta.ID]
        if len(list)==1{
            appendLog(fmt.Sprintf("%-11s : %-11s | %s",t.Meta.DisplayName,stateKorean(list[0].State),filepath.Base(list[0].Path)))
        }else if len(list)==0{
            appendLog(fmt.Sprintf("%-11s : 찾지 못함",t.Meta.DisplayName))
        }else{
            names:=[]string{};for _,m:=range list{names=append(names,filepath.Base(m.Path))}
            appendLog(fmt.Sprintf("%-11s : 중복/모호함 | %s",t.Meta.DisplayName,strings.Join(names,", ")))
        }
    }
    if len(scan.Ignored)>0{appendLog(fmt.Sprintf("기타/비호환 파일: %d개 (패치하지 않음)",len(scan.Ignored)))}
    ok:=len(problems)==0 && len(selected)==len(runtimeTargets)
    if ok{appendLog("[검사 완료] A~G + KANJI1 8개 대상이 모두 확인되었습니다.")}else{appendLog("[검사 실패] "+strings.Join(problems," | "))}
    enable(hwndPatchButton,ok)
    return ok
}

func doPatch() {
    if !scanAndDisplay(){return}
    if message("확인된 Disk A~G와 KANJI1 ROM을 모두 한글패치합니다.\n\n원본은 같은 폴더의 .bak 파일로 보존되고, 결과는 파일명 뒤에 (K)가 붙습니다.\n계속하시겠습니까?",MB_YESNO|MB_ICONQUESTION)!=IDYES{return}
    appendLog("[전체 패치] 8개 대상의 임시 패치 및 사후검증을 시작합니다.")
    logs,err:=batchPatch(currentFolder,runtimeTargets)
    if err!=nil{appendLog("[패치 실패] "+err.Error());message("패치하지 못했습니다.\n\n"+err.Error(),MB_OK|MB_ICONERROR);scanAndDisplay();return}
    for _,line:=range logs{appendLog(line)}
    appendLog("[완료] 바리스 II 한글패치 처리가 끝났습니다.")
    message("바리스 II Disk A~G와 KANJI1 ROM 처리가 완료되었습니다.",MB_OK|MB_ICONINFORMATION)
    scanAndDisplay()
}

func loword(v uintptr) int { return int(v & 0xffff) }

func wndProc(hwnd uintptr,msg uint32,wparam,lparam uintptr) uintptr {
    switch msg{
    case WM_COMMAND:
        switch loword(wparam){
        case ID_FOLDER:
            if p:=browseFolder();p!=""{currentFolder=p;setText(hwndFolderEdit,p);appendLog("폴더 선택: "+p);scanAndDisplay()}
            return 0
        case ID_SCAN: scanAndDisplay(); return 0
        case ID_PATCH: doPatch(); return 0
        case ID_EXIT: pDestroyWindow.Call(hwnd); return 0
        }
    case WM_SIZE:
        layout(hwnd); return 0
    case WM_DESTROY:
        pPostQuitMessage.Call(0); return 0
    }
    r,_,_:=pDefWindowProcW.Call(hwnd,uintptr(msg),wparam,lparam);return r
}

func createControl(exStyle uintptr,class,text string,style uintptr,x,y,w,h int,id int) uintptr {
    hwnd,_,_:=pCreateWindowExW.Call(exStyle,uintptr(unsafe.Pointer(u16(class))),uintptr(unsafe.Pointer(u16(text))),style,uintptr(x),uintptr(y),uintptr(w),uintptr(h),hwndMain,uintptr(id),0,0)
    if hwnd!=0{setFont(hwnd,normalFont)}
    return hwnd
}

func layout(hwnd uintptr){
    if hwndLog==0{return}
    var r RECT;pGetClientRect.Call(hwnd,uintptr(unsafe.Pointer(&r)))
    width:=int(r.Right-r.Left);height:=int(r.Bottom-r.Top)
    margin:=22; buttonW:=116; gap:=8
    pMoveWindow.Call(hwndFolderEdit,uintptr(margin),uintptr(120),uintptr(width-margin*2-buttonW-gap),uintptr(30),1)
    pMoveWindow.Call(hwndFolderButton,uintptr(width-margin-buttonW),uintptr(120),uintptr(buttonW),uintptr(30),1)
    pMoveWindow.Call(hwndScanButton,uintptr(margin),uintptr(166),uintptr(116),uintptr(34),1)
    pMoveWindow.Call(hwndPatchButton,uintptr(margin+124),uintptr(166),uintptr(174),uintptr(34),1)
    pMoveWindow.Call(hwndExitButton,uintptr(width-margin-90),uintptr(166),uintptr(90),uintptr(34),1)
    footerH:=52; logTop:=216; logH:=height-logTop-footerH-22
    if logH<180{logH=180}
    pMoveWindow.Call(hwndLog,uintptr(margin),uintptr(logTop),uintptr(width-margin*2),uintptr(logH),1)
    pMoveWindow.Call(hwndFooter,uintptr(margin),uintptr(logTop+logH+10),uintptr(width-margin*2),uintptr(footerH),1)
}

func createFonts(){
    normalFont,_,_=pCreateFontW.Call(^uintptr(17),0,0,0,400,0,0,0,1,0,0,0,0,uintptr(unsafe.Pointer(u16("Segoe UI"))))
    titleFont,_,_=pCreateFontW.Call(^uintptr(27),0,0,0,600,0,0,0,1,0,0,0,0,uintptr(unsafe.Pointer(u16("Segoe UI"))))
}

func runGUI() int {
    runtime.LockOSThread()
    defer runtime.UnlockOSThread()
    pCoInitializeEx.Call(0,COINIT_APARTMENTTHREADED);defer pCoUninitialize.Call()
    var err error
    runtimeTargets,err=loadRuntimeTargets()
    if err!=nil{pMessageBoxW.Call(0,uintptr(unsafe.Pointer(u16(err.Error()))),uintptr(unsafe.Pointer(u16("Patchy88 시작 오류"))),MB_OK|MB_ICONERROR);return 2}

    hinst,_,_:=pGetModuleHandleW.Call(0)
    cursor,_,_:=pLoadCursorW.Call(0,32512)
    className:=u16("Patchy88Valis2Window")
    wc:=WNDCLASSEX{CbSize:uint32(unsafe.Sizeof(WNDCLASSEX{})),Style:3,LpfnWndProc:syscall.NewCallback(wndProc),HInstance:hinst,HCursor:cursor,HbrBackground:6,LpszClassName:className}
    if r,_,_:=pRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)));r==0{return 2}
    hwnd,_,_:=pCreateWindowExW.Call(0,uintptr(unsafe.Pointer(className)),uintptr(unsafe.Pointer(u16("Patchy88 - 몽환전사 바리스 II 한국어 패치"))),WS_OVERLAPPEDWINDOW|WS_VISIBLE,100,80,920,720,0,0,hinst,0)
    if hwnd==0{return 2};hwndMain=hwnd
    createFonts()
    title:=createControl(0,"STATIC","몽환전사 바리스 II 한국어 패치",WS_CHILD|WS_VISIBLE,22,18,840,36,0);setFont(title,titleFont)
    createControl(0,"STATIC","폴더 하나만 지정하면 Disk A~G와 KANJI1 ROM을 내용으로 자동 식별하고 안전하게 일괄 패치합니다.",WS_CHILD|WS_VISIBLE,22,62,850,24,0)
    createControl(0,"STATIC","대상 폴더",WS_CHILD|WS_VISIBLE,22,96,100,20,0)
    hwndFolderEdit=createControl(WS_EX_CLIENTEDGE,"EDIT","",WS_CHILD|WS_VISIBLE|ES_READONLY,22,120,720,30,0)
    hwndFolderButton=createControl(0,"BUTTON","폴더 선택...",WS_CHILD|WS_VISIBLE|WS_TABSTOP,760,120,116,30,ID_FOLDER)
    hwndScanButton=createControl(0,"BUTTON","다시 검사",WS_CHILD|WS_VISIBLE|WS_TABSTOP,22,166,116,34,ID_SCAN)
    hwndPatchButton=createControl(0,"BUTTON","전체 한글패치 적용",WS_CHILD|WS_VISIBLE|WS_TABSTOP|BS_DEFPUSHBUTTON,146,166,174,34,ID_PATCH);enable(hwndPatchButton,false)
    hwndExitButton=createControl(0,"BUTTON","종료",WS_CHILD|WS_VISIBLE|WS_TABSTOP,808,166,90,34,ID_EXIT)
    hwndLog=createControl(WS_EX_CLIENTEDGE,"EDIT","",WS_CHILD|WS_VISIBLE|ES_READONLY|ES_MULTILINE|ES_AUTOVSCROLL|WS_VSCROLL,22,216,854,400,0)
    hwndFooter=createControl(0,"STATIC","정상 적용 시 원본은 같은 폴더의 <원본파일명>.bak으로 보존되고 결과는 <원본명>(K)<확장자>로 생성됩니다. 기존 백업은 덮어쓰지 않습니다.",WS_CHILD|WS_VISIBLE,22,628,854,44,0)
    appendLog("Patchy88 Valis II v"+appVersion)
    appendLog("폴더를 선택하면 A~G와 KANJI1을 자동 검사합니다.")
    layout(hwndMain)
    pShowWindow.Call(hwndMain,10);pUpdateWindow.Call(hwndMain)
    var msg MSG
    for {r,_,_:=pGetMessageW.Call(uintptr(unsafe.Pointer(&msg)),0,0,0);if int32(r)<=0{break};pTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)));pDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))}
    return int(msg.WParam)
}

func main(){runGUI()}
