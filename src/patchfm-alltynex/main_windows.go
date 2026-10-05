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
	comdlg32         = syscall.NewLazyDLL("comdlg32.dll")
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
	getOpenFileName  = comdlg32.NewProc("GetOpenFileNameW")
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
type openFileName struct {
	StructSize                                uint32
	Owner, Instance                           uintptr
	Filter, CustomFilter                      *uint16
	MaxCustomFilter, FilterIndex              uint32
	File                                      *uint16
	MaxFile                                   uint32
	FileTitle                                 *uint16
	MaxFileTitle                              uint32
	InitialDir, Title                         *uint16
	Flags                                     uint32
	FileOffset, FileExtension                 uint16
	DefExt                                    *uint16
	CustData, Hook, TemplateName, ReservedPtr uintptr
	Reserved, DwFlagsEx                       uint32
}

var uiHWND, uiPath, uiBrowse, uiPatch, uiExit, uiLog, uiFooter, fontNormal, fontTitle uintptr
var uiEngine *Engine
var sourcePath string
var busy bool
var pendingMu sync.Mutex
var pendingLogs []string
var pendingErr error
var pendingOut string

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
	x := pendingLogs
	pendingLogs = nil
	pendingMu.Unlock()
	for _, s := range x {
		logLine(s)
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
func pickZip() string {
	var buf [32768]uint16
	// OPENFILENAMEW filters are NUL-separated and double-NUL terminated.
	// Build each UTF-16 field separately because StringToUTF16 rejects embedded NUL.
	filter := append(syscall.StringToUTF16("ZIP 파일 (*.zip)"), syscall.StringToUTF16("*.zip")...)
	filter = append(filter, syscall.StringToUTF16("모든 파일 (*.*)")...)
	filter = append(filter, syscall.StringToUTF16("*.*")...)
	filter = append(filter, 0)
	of := openFileName{StructSize: uint32(unsafe.Sizeof(openFileName{})), Owner: uiHWND, Filter: &filter[0], File: &buf[0], MaxFile: uint32(len(buf)), Title: u16("Alltynex 원본 ZIP 선택"), Flags: 0x00001000 | 0x00000800}
	r, _, _ := getOpenFileName.Call(uintptr(unsafe.Pointer(&of)))
	if r == 0 {
		return ""
	}
	return syscall.UTF16ToString(buf[:])
}
func setBusy(v bool) {
	busy = v
	enable(uiBrowse, !v)
	enable(uiPatch, !v && sourcePath != "")
	enable(uiExit, !v)
}
func chooseSource() {
	if busy {
		return
	}
	p := pickZip()
	if p == "" {
		return
	}
	sourcePath = p
	textOn(uiPath, p)
	logLine("원본 검사: " + p)
	if err := uiEngine.ValidateSource(p); err != nil {
		sourcePath = ""
		enable(uiPatch, false)
		logLine("[검사 실패] " + err.Error())
		dialog("지원 원본이 아닙니다.\n\n"+err.Error(), mbOk|mbError)
		return
	}
	logLine("[검사 완료] 지원되는 Alltynex 원본 ZIP입니다.")
	enable(uiPatch, true)
}
func startPatch() {
	if busy || sourcePath == "" {
		return
	}
	if dialog("선택한 원본 ZIP에서 한글판 ISO를 생성합니다.\n\n원본 ZIP은 수정하지 않습니다.\n완성 ISO는 배포물에 포함되어 있지 않습니다.\n계속하시겠습니까?", mbYesNo|mbQuestion) != idYes {
		return
	}
	p := sourcePath
	setBusy(true)
	go func() {
		out, err := uiEngine.Apply(p, queuedLog)
		pendingMu.Lock()
		pendingOut = out
		pendingErr = err
		pendingMu.Unlock()
		postMsg.Call(uiHWND, wmDone, 0, 0)
	}()
}
func finishPatch() {
	flushLogs()
	pendingMu.Lock()
	out, err := pendingOut, pendingErr
	pendingOut = ""
	pendingErr = nil
	pendingMu.Unlock()
	if err != nil {
		logLine("[패치 실패] " + err.Error())
		dialog("패치 실패:\n"+err.Error(), mbOk|mbError)
	} else {
		logLine("[완료] " + out)
		dialog("한글판 ISO 생성 및 전체 해시 검증을 완료했습니다.\n\n"+out, mbOk|mbInfo)
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
	w := int(r.Right)
	h := int(r.Bottom)
	m := 22
	bw := 120
	moveWindow.Call(uiPath, uintptr(m), 124, uintptr(w-m*2-bw-8), 31, 1)
	moveWindow.Call(uiBrowse, uintptr(w-m-bw), 124, uintptr(bw), 31, 1)
	moveWindow.Call(uiPatch, uintptr(m), 170, 180, 34, 1)
	moveWindow.Call(uiExit, uintptr(w-m-90), 170, 90, 34, 1)
	top := 220
	fh := 60
	lh := h - top - fh - 25
	if lh < 160 {
		lh = 160
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
	name := u16("PatchFMAlltynexMainWindow")
	wc := wndClass{Size: uint32(unsafe.Sizeof(wndClass{})), Style: 3, WndProc: syscall.NewCallback(wndProc), Instance: hinst, Cursor: cursor, Background: 6, Name: name}
	if r, _, _ := registerClass.Call(uintptr(unsafe.Pointer(&wc))); r == 0 {
		return 2
	}
	uiHWND, _, _ = createWindow.Call(0, uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(u16(programName+" — FM TOWNS 한글패치"))), wsWindow|wsVisible, 120, 80, 900, 700, 0, 0, hinst, 0)
	if uiHWND == 0 {
		return 2
	}
	h1 := int32(-17)
	fontNormal, _, _ = createFont.Call(uintptr(h1), 0, 0, 0, 400, 0, 0, 0, 1, 0, 0, 0, 0, uintptr(unsafe.Pointer(u16("맑은 고딕"))))
	h2 := int32(-29)
	fontTitle, _, _ = createFont.Call(uintptr(h2), 0, 0, 0, 600, 0, 0, 0, 1, 0, 0, 0, 0, uintptr(unsafe.Pointer(u16("맑은 고딕"))))
	title := createControl(0, "STATIC", "PatchFM — Alltynex", wsChild|wsVisible, 22, 18, 820, 41, 0)
	setFont(title, fontTitle)
	createControl(0, "STATIC", "FM TOWNS 《Alltynex》 — 원본 ZIP을 검증하고 한글판 ISO를 생성합니다.", wsChild|wsVisible, 22, 69, 820, 27, 0)
	createControl(0, "STATIC", "원본 ZIP", wsChild|wsVisible, 22, 102, 160, 20, 0)
	uiPath = createControl(wsClientEdge, "EDIT", "", wsChild|wsVisible|esReadonly, 22, 124, 700, 31, 0)
	uiBrowse = createControl(0, "BUTTON", "ZIP 선택...", wsChild|wsVisible|wsTabstop, 735, 124, 120, 31, idBrowse)
	uiPatch = createControl(0, "BUTTON", "한글패치 적용", wsChild|wsVisible|wsTabstop|btnDefault, 22, 170, 180, 34, idPatch)
	uiExit = createControl(0, "BUTTON", "종료", wsChild|wsVisible|wsTabstop, 765, 170, 90, 34, idExit)
	uiLog = createControl(wsClientEdge, "EDIT", "", wsChild|wsVisible|esReadonly|esMultiline|esAutoScroll|wsVScroll, 22, 220, 833, 380, 0)
	uiFooter = createControl(0, "STATIC", "원본 ZIP과 완성 ISO는 PatchFM 배포물에 포함되지 않습니다. 원본 ZIP은 수정하지 않습니다.\r\n결과 ISO는 전체 크기·MD5·SHA-256 검증을 통과한 경우에만 확정합니다.", wsChild|wsVisible, 22, 608, 833, 60, 0)
	enable(uiPatch, false)
	logLine(programName)
	logLine("xdelta 및 xdelta3 무결성 검증 완료")
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

var _ = fmt.Sprintf
