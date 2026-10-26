//go:build windows

package main

import (
	"context"
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
	pathFromIDList   = shell32.NewProc("SHGetPathFromIDListW")
	coTaskMemFree    = ole32.NewProc("CoTaskMemFree")
	coInitializeEx   = ole32.NewProc("CoInitializeEx")
	coUninitialize   = ole32.NewProc("CoUninitialize")
)

const (
	wmDestroy    = 0x0002
	wmSize       = 0x0005
	wmClose      = 0x0010
	wmCommand    = 0x0111
	wmSetFont    = 0x0030
	wmLog        = 0x8001
	wmDone       = 0x8002
	emSetSel     = 0x00B1
	emReplaceSel = 0x00C2
	wsWindow     = 0x00CF0000
	wsVisible    = 0x10000000
	wsChild      = 0x40000000
	wsTabstop    = 0x00010000
	wsVScroll    = 0x00200000
	esReadonly   = 0x0800
	esMultiline  = 0x0004
	esAutoScroll = 0x0040
	wsClientEdge = 0x00000200
	btnDefault   = 0x1
	mbOk         = 0
	mbInfo       = 0x40
	mbError      = 0x10
	mbQuestion   = 0x20
	mbYesNo      = 4
	idYes        = 6
	idBrowse     = 1001
	idPatch      = 1002
	idExit       = 1003
	idCancel     = 1004
)

type wndClass struct {
	Size, Style                        uint32
	WndProc                            uintptr
	ClsExtra, WndExtra                 int32
	Instance, Icon, Cursor, Background uintptr
	MenuName, Name                     *uint16
	IconSmall                          uintptr
}
type point struct{ X, Y int32 }
type msgStruct struct {
	HWND           uintptr
	Message        uint32
	WParam, LParam uintptr
	Time           uint32
	Pt             point
	Private        uint32
}
type rect struct{ Left, Top, Right, Bottom int32 }
type browseInfo struct {
	Owner, Root, DisplayName, Title uintptr
	Flags                           uint32
	Callback, LParam                uintptr
	Image                           int32
}

var uiHWND, uiPath, uiBrowse, uiPatch, uiCancel, uiExit, uiLog, uiFooter, fontNormal, fontTitle uintptr
var uiEngine *Engine
var sourceFolder string
var busy bool
var cancelPatch context.CancelFunc
var pendingMu sync.Mutex
var pendingLogs []string
var pendingErr error
var pendingBIN, pendingCUE string

func u16(s string) *uint16 { return syscall.StringToUTF16Ptr(s) }
func send(h uintptr, m uint32, w, l uintptr) uintptr {
	r, _, _ := sendMessage.Call(h, uintptr(m), w, l)
	return r
}
func setFont(h, f uintptr) {
	if h != 0 && f != 0 {
		send(h, wmSetFont, f, 1)
	}
}
func textOn(h uintptr, s string) { setWindowText.Call(h, uintptr(unsafe.Pointer(u16(s)))) }
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
	lines := pendingLogs
	pendingLogs = nil
	pendingMu.Unlock()
	for _, line := range lines {
		logLine(line)
	}
}
func dialog(s string, flags uintptr) int {
	r, _, _ := msgBox.Call(uiHWND, uintptr(unsafe.Pointer(u16(s))), uintptr(unsafe.Pointer(u16(programName))), flags)
	return int(r)
}
func enable(h uintptr, on bool) {
	v := uintptr(0)
	if on {
		v = 1
	}
	enableWindow.Call(h, v)
}

// The shell folder picker returns a PIDL. Always free it with CoTaskMemFree.
// Selecting a *folder*, not an individual CUE/IMG/CCD, is intentional.
func pickFolder() string {
	// The new-style shell dialog needs COM initialized on the caller's STA.
	hr, _, _ := coInitializeEx.Call(0, 0x2)     // COINIT_APARTMENTTHREADED
	comOK := uint32(hr) == 0 || uint32(hr) == 1 // S_OK / S_FALSE
	if comOK {
		defer coUninitialize.Call()
	}
	flags := uint32(0x1) // BIF_RETURNONLYFSDIRS
	if comOK {
		flags |= 0x40 | 0x10
	} // NEWDIALOGSTYLE | EDITBOX
	var display [260]uint16
	title := u16("Dune II 원본 CD 이미지가 있는 폴더를 선택하세요")
	info := browseInfo{
		Owner:       uiHWND,
		DisplayName: uintptr(unsafe.Pointer(&display[0])),
		Title:       uintptr(unsafe.Pointer(title)),
		Flags:       flags,
	}
	pidl, _, _ := browseForFolder.Call(uintptr(unsafe.Pointer(&info)))
	runtime.KeepAlive(title)
	if pidl == 0 {
		return ""
	}
	defer coTaskMemFree.Call(pidl)
	var path [32768]uint16
	ok, _, _ := pathFromIDList.Call(pidl, uintptr(unsafe.Pointer(&path[0])))
	if ok == 0 {
		return ""
	}
	return syscall.UTF16ToString(path[:])
}
func setBusy(v bool) {
	busy = v
	enable(uiBrowse, !v)
	enable(uiPatch, !v && sourceFolder != "")
	enable(uiExit, !v)
	enable(uiCancel, v)
}
func chooseSource() {
	if busy {
		return
	}
	p := pickFolder()
	if p == "" {
		return
	}
	sourceFolder = ""
	enable(uiPatch, false)
	textOn(uiPath, p)
	logLine("원본 폴더 자동 검색: " + p)
	src, err := uiEngine.FindSource(p)
	if err != nil {
		logLine("[검사 실패] " + err.Error())
		dialog("원본 이미지를 찾지 못했습니다.\n\n"+err.Error(), mbOk|mbError)
		return
	}
	sourceFolder = p
	logLine("[검사 완료] 지원 원본 자동 인식: " + src)
	enable(uiPatch, true)
}
func startPatch() {
	if busy || sourceFolder == "" {
		return
	}
	if dialog("지정한 폴더에서 자동 인식한 원본 IMG/BIN으로 한글판 BIN/CUE를 생성합니다.\n\n원본 IMG/CCD/SUB/CUE는 수정하지 않습니다.\n게임 원본과 완성 BIN은 배포물에 포함되지 않습니다.\n계속하시겠습니까?", mbYesNo|mbQuestion) != idYes {
		return
	}
	p := sourceFolder
	ctx, cancel := context.WithCancel(context.Background())
	cancelPatch = cancel
	setBusy(true)
	go func() {
		bin, cue, err := uiEngine.ApplyContext(ctx, p, queuedLog)
		pendingMu.Lock()
		pendingBIN = bin
		pendingCUE = cue
		pendingErr = err
		pendingMu.Unlock()
		postMsg.Call(uiHWND, wmDone, 0, 0)
	}()
}
func finishPatch() {
	flushLogs()
	pendingMu.Lock()
	bin, cue, err := pendingBIN, pendingCUE, pendingErr
	pendingBIN = ""
	pendingCUE = ""
	if cancelPatch != nil {
		cancelPatch()
		cancelPatch = nil
	}
	pendingErr = nil
	pendingMu.Unlock()
	if err != nil {
		logLine("[패치 실패] " + err.Error())
		dialog("패치 실패:\n"+err.Error(), mbOk|mbError)
	} else {
		logLine("[완료] BIN: " + bin)
		logLine("[완료] CUE: " + cue)
		dialog("한글판 BIN/CUE 생성 및 전체 해시 검증 완료.\n\n"+bin+"\n"+cue, mbOk|mbInfo)
	}
	setBusy(false)
}
func createControl(ex uintptr, klass, label string, style uintptr, x, y, w, h, id int) uintptr {
	r, _, _ := createWindow.Call(ex, uintptr(unsafe.Pointer(u16(klass))), uintptr(unsafe.Pointer(u16(label))), style, uintptr(x), uintptr(y), uintptr(w), uintptr(h), uiHWND, uintptr(id), 0, 0)
	if r != 0 {
		setFont(r, fontNormal)
	}
	return r
}
func layout() {
	if uiLog == 0 {
		return
	}
	var r rect
	getClientRect.Call(uiHWND, uintptr(unsafe.Pointer(&r)))
	w, h := int(r.Right), int(r.Bottom)
	m, bw := 22, 140
	if w < 570 {
		return
	}
	moveWindow.Call(uiPath, uintptr(m), 124, uintptr(w-m*2-bw-8), 31, 1)
	moveWindow.Call(uiBrowse, uintptr(w-m-bw), 124, uintptr(bw), 31, 1)
	moveWindow.Call(uiPatch, uintptr(m), 170, 180, 34, 1)
	moveWindow.Call(uiCancel, uintptr(m+195), 170, 140, 34, 1)
	moveWindow.Call(uiExit, uintptr(w-m-90), 170, 90, 34, 1)
	top, fh := 220, 72
	lh := h - top - fh - 25
	if lh < 150 {
		lh = 150
	}
	moveWindow.Call(uiLog, uintptr(m), uintptr(top), uintptr(w-m*2), uintptr(lh), 1)
	moveWindow.Call(uiFooter, uintptr(m), uintptr(top+lh+8), uintptr(w-m*2), uintptr(fh), 1)
}
func wndProc(h uintptr, m uint32, w, l uintptr) uintptr {
	switch m {
	case wmCommand:
		switch int(w & 0xffff) {
		case idBrowse:
			chooseSource()
			return 0
		case idPatch:
			startPatch()
			return 0
		case idCancel:
			if busy && cancelPatch != nil {
				logLine("[사용자 중단] xdelta3 종료와 임시 파일 정리를 요청했습니다.")
				cancelPatch()
				enable(uiCancel, false)
			}
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
	case wmDone:
		finishPatch()
		return 0
	case wmClose:
		if busy {
			dialog("패치 작업 중에는 종료할 수 없습니다.", mbOk|mbInfo)
			return 0
		}
		destroyWindow.Call(h)
		return 0
	case wmDestroy:
		postQuit.Call(0)
		return 0
	}
	r, _, _ := defaultWndProc.Call(h, uintptr(m), w, l)
	return r
}
func runGUI() int {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
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
	name := u16("PatchyDOSDune2MainWindow")
	wc := wndClass{Size: uint32(unsafe.Sizeof(wndClass{})), Style: 3, WndProc: syscall.NewCallback(wndProc), Instance: hinst, Cursor: cursor, Background: 6, Name: name}
	if r, _, _ := registerClass.Call(uintptr(unsafe.Pointer(&wc))); r == 0 {
		return 2
	}
	uiHWND, _, _ = createWindow.Call(0, uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(u16(programName+" — MS-DOS CD 한글패치"))), wsWindow|wsVisible, 120, 80, 900, 700, 0, 0, hinst, 0)
	if uiHWND == 0 {
		return 2
	}
	h1 := int32(-17)
	fontNormal, _, _ = createFont.Call(uintptr(h1), 0, 0, 0, 400, 0, 0, 0, 1, 0, 0, 0, 0, uintptr(unsafe.Pointer(u16("맑은 고딕"))))
	h2 := int32(-29)
	fontTitle, _, _ = createFont.Call(uintptr(h2), 0, 0, 0, 600, 0, 0, 0, 1, 0, 0, 0, 0, uintptr(unsafe.Pointer(u16("맑은 고딕"))))
	title := createControl(0, "STATIC", "PatchyDOS — Dune II: Battle for Arrakis", wsChild|wsVisible, 22, 18, 840, 41, 0)
	setFont(title, fontTitle)
	createControl(0, "STATIC", "MS-DOS 《Dune II: Battle for Arrakis》 — CD 원본 자동 검색 → 한글판 BIN/CUE 생성", wsChild|wsVisible, 22, 69, 840, 27, 0)
	createControl(0, "STATIC", "원본 이미지 폴더 (하위 폴더 검색 포함)", wsChild|wsVisible, 22, 102, 240, 20, 0)
	uiPath = createControl(wsClientEdge, "EDIT", "", wsChild|wsVisible|esReadonly, 22, 124, 692, 31, 0)
	uiBrowse = createControl(0, "BUTTON", "폴더 선택...", wsChild|wsVisible|wsTabstop, 730, 124, 140, 31, idBrowse)
	uiPatch = createControl(0, "BUTTON", "한글패치 적용", wsChild|wsVisible|wsTabstop|btnDefault, 22, 170, 180, 34, idPatch)
	uiCancel = createControl(0, "BUTTON", "작업 중단", wsChild|wsVisible|wsTabstop, 220, 170, 140, 34, idCancel)
	uiExit = createControl(0, "BUTTON", "종료", wsChild|wsVisible|wsTabstop, 780, 170, 90, 34, idExit)
	uiLog = createControl(wsClientEdge, "EDIT", "", wsChild|wsVisible|esReadonly|esMultiline|esAutoScroll|wsVScroll, 22, 220, 844, 370, 0)
	uiFooter = createControl(0, "STATIC", "원본 IMG/CCD/SUB/CUE 및 완성 한글판 BIN은 배포물에 포함되지 않습니다.\r\n선택 폴더 및 하위 폴더에서 원본 IMG/BIN을 해시로 자동 인식합니다.\r\n원본은 변경하지 않으며, 기존 결과가 다르면 덮어쓰지 않습니다.", wsChild|wsVisible, 22, 605, 844, 72, 0)
	enable(uiPatch, false)
	enable(uiCancel, false)
	logLine(programName)
	logLine("xdelta, CUE 템플릿, xdelta3 무결성 검사 완료")
	layout()
	showWindow.Call(uiHWND, 10)
	updateWindow.Call(uiHWND)
	var ms msgStruct
	for {
		r, _, _ := getMessage.Call(uintptr(unsafe.Pointer(&ms)), 0, 0, 0)
		if int32(r) <= 0 {
			break
		}
		translateMessage.Call(uintptr(unsafe.Pointer(&ms)))
		dispatchMessage.Call(uintptr(unsafe.Pointer(&ms)))
	}
	return int(ms.WParam)
}
func main() { os.Exit(runGUI()) }
