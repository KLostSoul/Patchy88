//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	gdi32            = syscall.NewLazyDLL("gdi32.dll")
	shell32          = syscall.NewLazyDLL("shell32.dll")
	ole32            = syscall.NewLazyDLL("ole32.dll")
	registerClass    = user32.NewProc("RegisterClassExW")
	createWindow     = user32.NewProc("CreateWindowExW")
	defaultWndProc   = user32.NewProc("DefWindowProcW")
	postMsg          = user32.NewProc("PostMessageW")
	showWindow       = user32.NewProc("ShowWindow")
	updateWindow     = user32.NewProc("UpdateWindow")
	getMessage       = user32.NewProc("GetMessageW")
	translateMessage = user32.NewProc("TranslateMessage")
	dispatchMessage  = user32.NewProc("DispatchMessageW")
	postQuit         = user32.NewProc("PostQuitMessage")
	destroyWindow    = user32.NewProc("DestroyWindow")
	sendMessage      = user32.NewProc("SendMessageW")
	setWindowText    = user32.NewProc("SetWindowTextW")
	msgBox           = user32.NewProc("MessageBoxW")
	enableWindow     = user32.NewProc("EnableWindow")
	moveWindow       = user32.NewProc("MoveWindow")
	getClientRect    = user32.NewProc("GetClientRect")
	loadCursor       = user32.NewProc("LoadCursorW")
	getModuleHandle  = kernel32.NewProc("GetModuleHandleW")
	createFont       = gdi32.NewProc("CreateFontW")
	browseForFolder  = shell32.NewProc("SHBrowseForFolderW")
	folderPath       = shell32.NewProc("SHGetPathFromIDListW")
	coTaskFree       = ole32.NewProc("CoTaskMemFree")
	coInit           = ole32.NewProc("CoInitializeEx")
	coUninit         = ole32.NewProc("CoUninitialize")
)

const (
	wmDestroy    = 0x0002
	wmSize       = 0x0005
	wmClose      = 0x0010
	wmCommand    = 0x0111
	wmSetFont    = 0x0030
	wmLog        = 0x8001
	wmScanDone   = 0x8002
	wmPatchDone  = 0x8003
	emSetSel     = 0x00B1
	emReplaceSel = 0x00C2
	wsWindow     = 0x00CF0000
	wsVisible    = 0x10000000
	wsChild      = 0x40000000
	wsTabstop    = 0x00010000
	wsBorder     = 0x00800000
	wsVScroll    = 0x00200000
	esReadonly   = 0x0800
	esMultiline  = 0x0004
	esAutoScroll = 0x0040
	wsClientEdge = 0x00000200
	btnDefault   = 0x00000001
	mbOk         = 0
	mbInfo       = 0x40
	mbError      = 0x10
	mbQuestion   = 0x20
	mbYesNo      = 4
	mbYesNoCancel = 3
	idYes        = 6
	idNo         = 7
	folderFlags  = 0x41
	idBrowse     = 1001
	idScan       = 1002
	idPatch      = 1003
	idExit       = 1004
)

type wndClass struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   uintptr
	Icon       uintptr
	Cursor     uintptr
	Background uintptr
	MenuName   *uint16
	Name       *uint16
	IconSmall  uintptr
}
type point struct{ X, Y int32 }
type msgStruct struct {
	HWND    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      point
	Private uint32
}
type rect struct{ Left, Top, Right, Bottom int32 }
type browseInfo struct {
	Owner    uintptr
	Root     uintptr
	Display  *uint16
	Title    *uint16
	Flags    uint32
	Callback uintptr
	Data     uintptr
	Image    int32
}

var (
	uiHWND       uintptr
	uiPath       uintptr
	uiBrowse     uintptr
	uiScan       uintptr
	uiPatch      uintptr
	uiExit       uintptr
	uiLog        uintptr
	uiFooter     uintptr
	fontNormal   uintptr
	fontTitle    uintptr
	uiEngine     *Engine
	chosenFolder string
	latestScan   *ScanResult
	busy         bool
	pendingMu    sync.Mutex
	pendingLogs  []string
	pendingScan  *ScanResult
	pendingErr   error
)

func u16(s string) *uint16       { return syscall.StringToUTF16Ptr(s) }
func textOn(h uintptr, s string) { setWindowText.Call(h, uintptr(unsafe.Pointer(u16(s)))) }
func send(h uintptr, msg uint32, w, l uintptr) uintptr {
	r, _, _ := sendMessage.Call(h, uintptr(msg), w, l)
	return r
}
func setFont(h, f uintptr) {
	if h != 0 && f != 0 {
		send(h, wmSetFont, f, 1)
	}
}
func enable(h uintptr, on bool) {
	v := uintptr(0)
	if on {
		v = 1
	}
	enableWindow.Call(h, v)
}
func logLine(s string) {
	if uiLog == 0 {
		return
	}
	s = strings.ReplaceAll(s, "\n", "\r\n") + "\r\n"
	send(uiLog, emSetSel, ^uintptr(0), ^uintptr(0))
	send(uiLog, emReplaceSel, 0, uintptr(unsafe.Pointer(u16(s))))
}
func queuedLog(s string) {
	pendingMu.Lock()
	pendingLogs = append(pendingLogs, s)
	pendingMu.Unlock()
	postMsg.Call(uiHWND, wmLog, 0, 0)
}
func flushLogs() {
	pendingMu.Lock()
	logs := pendingLogs
	pendingLogs = nil
	pendingMu.Unlock()
	for _, s := range logs {
		logLine(s)
	}
}
func dialog(s string, flags uintptr) int {
	r, _, _ := msgBox.Call(uiHWND, uintptr(unsafe.Pointer(u16(s))), uintptr(unsafe.Pointer(u16(programName))), flags)
	return int(r)
}
func pickFolder() string {
	var display [260]uint16
	data := browseInfo{Owner: uiHWND, Display: &display[0], Title: u16("일본판 또는 영문판 CCD / IMG / SUB가 있는 폴더를 선택하세요."), Flags: folderFlags}
	pidl, _, _ := browseForFolder.Call(uintptr(unsafe.Pointer(&data)))
	if pidl == 0 {
		return ""
	}
	defer coTaskFree.Call(pidl)
	var path [32768]uint16
	ok, _, _ := folderPath.Call(pidl, uintptr(unsafe.Pointer(&path[0])))
	if ok == 0 {
		return ""
	}
	return syscall.UTF16ToString(path[:])
}
func setBusy(isBusy bool) {
	busy = isBusy
	enable(uiBrowse, !busy)
	enable(uiScan, !busy)
	enable(uiPatch, !busy && latestScan != nil)
	enable(uiExit, !busy)
}
func startScan() {
	if busy {
		return
	}
	if chosenFolder == "" {
		dialog("폴더를 선택하세요.", mbOk|mbInfo)
		return
	}
	latestScan = nil
	setBusy(true)
	logLine("────────────────────────────────────────────")
	logLine("검사: " + chosenFolder)
	folder := chosenFolder
	go func() {
		r, err := uiEngine.Scan(folder)
		pendingMu.Lock()
		pendingScan = r
		pendingErr = err
		pendingMu.Unlock()
		postMsg.Call(uiHWND, wmScanDone, 0, 0)
	}()
}
func finishedScan() {
	flushLogs()
	pendingMu.Lock()
	r, err := pendingScan, pendingErr
	pendingScan = nil
	pendingErr = nil
	pendingMu.Unlock()
	if err != nil {
		logLine("[검사 실패] " + err.Error())
		latestScan = nil
		setBusy(false)
		return
	}
	if r.Edition == "" && len(r.Options) == 2 {
		logLine("일본판과 영문판 모두 발견됨 — 사용할 원본 선택 필요")
		choice := dialog("일본판과 영문판이 모두 발견됐습니다.\n\n예(Y): 일본판 패치\n아니오(N): 영문판 패치\n취소: 아무것도 적용하지 않음", mbYesNoCancel|mbQuestion)
		selected := ""
		switch choice {
		case idYes:
			selected = "Japanese"
		case idNo:
			selected = "English"
		default:
			logLine("선택 취소 — 다시 검사하면 다시 선택할 수 있습니다")
			latestScan = nil
			setBusy(false)
			return
		}
		var err error
		r, err = uiEngine.ScanWithEdition(r.Folder, selected)
		if err != nil {
			logLine("[원본 선택 실패] " + err.Error())
			latestScan = nil
			setBusy(false)
			return
		}
		logLine("사용자 선택: " + selected)
	}
	latestScan = r
	logLine("판별: " + r.Edition)
	for _, ext := range exts {
		if r.Already[ext] {
			logLine(strings.ToUpper(ext) + ": 한글판 결과 검증됨")
		} else {
			logLine(strings.ToUpper(ext) + " 원본: " + filepath.Base(r.Inputs[ext]))
		}
	}
	for _, n := range r.Notes {
		logLine(n)
	}
	extras := 0
	for _, x := range uiEngine.Manifest.Extras {
		if r.ExtrasPresent[x.Filename] {
			extras++
		}
	}
	logLine(fmt.Sprintf("추가 파일: %d / 3개 이미 존재", extras))
	logLine("[검사 완료] 패치를 적용할 수 있습니다.")
	setBusy(false)
}
func startPatch() {
	if busy || latestScan == nil {
		return
	}
	chosen := map[string]string{"Japanese": "일본판", "English": "영문판", "AlreadyPatched": "이미 패치됨"}[latestScan.Edition]
	if dialog("선택한 원본: "+chosen+"\n\n한글판 CCD/IMG/SUB를 생성하고 Mirrors_Kor1.00.cue와 D88 2개를 추가합니다.\n원본은 수정하지 않고 기존 파일도 덮어쓰지 않습니다.\n계속하시겠습니까?", mbYesNo|mbQuestion) != idYes {
		return
	}
	s := latestScan
	latestScan = nil
	setBusy(true)
	go func() {
		err := uiEngine.Apply(s, queuedLog)
		pendingMu.Lock()
		pendingErr = err
		pendingMu.Unlock()
		postMsg.Call(uiHWND, wmPatchDone, 0, 0)
	}()
}
func finishedPatch() {
	flushLogs()
	pendingMu.Lock()
	err := pendingErr
	pendingErr = nil
	pendingMu.Unlock()
	if err != nil {
		logLine("[패치 실패] " + err.Error())
		dialog("패치 실패:\n"+err.Error()+"\n\n원본 파일은 변경되지 않았습니다.", mbOk|mbError)
	} else {
		logLine("[완료] Mirrors_Kor1.00 CCD/IMG/SUB + CUE + D88 2개")
		dialog("한글판과 추가 파일이 같은 폴더에 생성됐습니다.\n결과 MD5 및 SHA-256 검증을 통과했습니다.", mbOk|mbInfo)
	}
	setBusy(false)
	// Refresh asynchronously to support re-runs and skip valid preexisting output files.
	startScan()
}
func createControl(ex uintptr, klass, label string, style uintptr, x, y, w, h, id int) uintptr {
	handle, _, _ := createWindow.Call(ex, uintptr(unsafe.Pointer(u16(klass))), uintptr(unsafe.Pointer(u16(label))), style, uintptr(x), uintptr(y), uintptr(w), uintptr(h), uiHWND, uintptr(id), 0, 0)
	if handle != 0 {
		setFont(handle, fontNormal)
	}
	return handle
}
func layout() {
	if uiLog == 0 {
		return
	}
	var r rect
	getClientRect.Call(uiHWND, uintptr(unsafe.Pointer(&r)))
	width := int(r.Right - r.Left)
	height := int(r.Bottom - r.Top)
	margin := 22
	bw := 120
	gap := 8
	if width < 650 {
		width = 650
	}
	moveWindow.Call(uiPath, uintptr(margin), 124, uintptr(width-margin*2-bw-gap), 31, 1)
	moveWindow.Call(uiBrowse, uintptr(width-margin-bw), 124, uintptr(bw), 31, 1)
	moveWindow.Call(uiScan, uintptr(margin), 170, 124, 34, 1)
	moveWindow.Call(uiPatch, uintptr(margin+132), 170, 186, 34, 1)
	moveWindow.Call(uiExit, uintptr(width-margin-90), 170, 90, 34, 1)
	logTop := 220
	footerH := 58
	logH := height - logTop - footerH - 25
	if logH < 160 {
		logH = 160
	}
	moveWindow.Call(uiLog, uintptr(margin), uintptr(logTop), uintptr(width-margin*2), uintptr(logH), 1)
	moveWindow.Call(uiFooter, uintptr(margin), uintptr(logTop+logH+8), uintptr(width-margin*2), uintptr(footerH), 1)
}
func wndProc(h uintptr, m uint32, w, l uintptr) uintptr {
	switch m {
	case wmCommand:
		switch int(w & 0xffff) {
		case idBrowse:
			if !busy {
				if p := pickFolder(); p != "" {
					chosenFolder = p
					textOn(uiPath, p)
					startScan()
				}
			}
			return 0
		case idScan:
			startScan()
			return 0
		case idPatch:
			startPatch()
			return 0
		case idExit:
			if !busy {
				destroyWindow.Call(h)
			}
			return 0
		}
	case wmSize:
		layout()
		return 0
	case wmLog:
		flushLogs()
		return 0
	case wmScanDone:
		finishedScan()
		return 0
	case wmPatchDone:
		finishedPatch()
		return 0
	case wmClose:
		if busy {
			dialog("검사 또는 패치 작업이 진행 중입니다. 완료 후 종료하세요.", mbOk|mbInfo)
			return 0
		}
		destroyWindow.Call(h)
		return 0
	case wmDestroy:
		postQuit.Call(0)
		return 0
	}
	ret, _, _ := defaultWndProc.Call(h, uintptr(m), w, l)
	return ret
}
func runGUI() int {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	coInit.Call(0, 2)
	defer coUninit.Call()
	exe, err := os.Executable()
	if err != nil {
		return 2
	}
	uiEngine, err = NewEngine(filepath.Join(filepath.Dir(exe), "assets"))
	if err != nil {
		msgBox.Call(0, uintptr(unsafe.Pointer(u16("필수 파일 검증 실패:\n"+err.Error()))), uintptr(unsafe.Pointer(u16(programName))), mbOk|mbError)
		return 2
	}
	hinst, _, _ := getModuleHandle.Call(0)
	cursor, _, _ := loadCursor.Call(0, 32512)
	name := u16("MirrorsKor100MainWindow")
	wc := wndClass{Size: uint32(unsafe.Sizeof(wndClass{})), Style: 3, WndProc: syscall.NewCallback(wndProc), Instance: hinst, Cursor: cursor, Background: 6, Name: name}
	if r, _, _ := registerClass.Call(uintptr(unsafe.Pointer(&wc))); r == 0 {
		return 2
	}
	uiHWND, _, _ = createWindow.Call(0, uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(u16(programName+" — xdelta 한글패치"))), wsWindow|wsVisible, 100, 70, 945, 735, 0, 0, hinst, 0)
	if uiHWND == 0 {
		return 2
	}
	// Negative font height gives consistent physical text sizing under Win32.
	h1 := int32(-17)
	fontNormal, _, _ = createFont.Call(uintptr(h1), 0, 0, 0, 400, 0, 0, 0, 1, 0, 0, 0, 0, uintptr(unsafe.Pointer(u16("맑은 고딕"))))
	h2 := int32(-29)
	fontTitle, _, _ = createFont.Call(uintptr(h2), 0, 0, 0, 600, 0, 0, 0, 1, 0, 0, 0, 0, uintptr(unsafe.Pointer(u16("맑은 고딕"))))
	title := createControl(0, "STATIC", "Mirrors_Kor1.00", wsChild|wsVisible, 22, 18, 880, 41, 0)
	setFont(title, fontTitle)
	createControl(0, "STATIC", "PC-8801 《Mirrors》 — 일본판과 영문판 모두 있으면 적용할 원본을 선택합니다.", wsChild|wsVisible, 22, 69, 870, 27, 0)
	createControl(0, "STATIC", "원본 파일 폴더", wsChild|wsVisible, 22, 102, 160, 20, 0)
	uiPath = createControl(wsClientEdge, "EDIT", "", wsChild|wsVisible|esReadonly, 22, 124, 730, 31, 0)
	uiBrowse = createControl(0, "BUTTON", "폴더 선택...", wsChild|wsVisible|wsTabstop, 765, 124, 120, 31, idBrowse)
	uiScan = createControl(0, "BUTTON", "다시 검사", wsChild|wsVisible|wsTabstop, 22, 170, 124, 34, idScan)
	uiPatch = createControl(0, "BUTTON", "한글패치 적용", wsChild|wsVisible|wsTabstop|btnDefault, 154, 170, 186, 34, idPatch)
	uiExit = createControl(0, "BUTTON", "종료", wsChild|wsVisible|wsTabstop, 790, 170, 90, 34, idExit)
	uiLog = createControl(wsClientEdge, "EDIT", "", wsChild|wsVisible|esReadonly|esMultiline|esAutoScroll|wsVScroll, 22, 220, 860, 405, 0)
	uiFooter = createControl(0, "STATIC", "원본은 수정하지 않습니다. 결과 CCD/IMG/SUB, Mirrors_Kor1.00.cue, D88 2개를 같은 폴더에 생성합니다.\r\nMD5 + SHA-256 검증 실패 시 결과를 확정하지 않습니다.", wsChild|wsVisible, 22, 630, 870, 58, 0)
	enable(uiPatch, false)
	logLine(programName + " — 일본판/영문판 별도 xdelta 경로")
	logLine("필수 패치 6개, xdelta3.exe, CUE, D88 2개 무결성 검증 완료")
	layout()
	showWindow.Call(uiHWND, 10)
	updateWindow.Call(uiHWND)
	var m msgStruct
	for {
		ret, _, _ := getMessage.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(ret) <= 0 {
			break
		}
		translateMessage.Call(uintptr(unsafe.Pointer(&m)))
		dispatchMessage.Call(uintptr(unsafe.Pointer(&m)))
	}
	return int(m.WParam)
}
func main() { os.Exit(runGUI()) }
