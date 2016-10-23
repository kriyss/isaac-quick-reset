package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"

	"github.com/AllenDang/w32"
	"github.com/simulatedsimian/joystick"
)

const (
	xbox360StartButton = 128
	keyboardRButton    = 82
	keyEventDown       = 0
	keyEventUp         = 0x0002
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procSetForegroundWindow = user32.NewProc("SetForegroundWindow")
	procEnumWindows         = user32.NewProc("EnumWindows")
	procGetWindowTextW      = user32.NewProc("GetWindowTextW")
)

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	hwnd, err := findWindow("Binding of Isaac: Afterbirth")
	exitOnError(err)
	setForegroundWindow(hwnd)
	time.Sleep(time.Second * 3)

	js, err := joystick.Open(0)
	exitOnError(err)

	ticker := time.NewTicker(time.Duration(100) * time.Millisecond)

	for {
		<-ticker.C
		state, err := js.Read()
		exitOnError(err)
		if state.Buttons == xbox360StartButton {
			pressKey(keyboardRButton, keyEventDown)
		} else {
			pressKey(keyboardRButton, keyEventUp)
		}
	}
}

func pressKey(vk uint16, flags int) {
	w32.SendInput([]w32.INPUT{
		w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			Ki: w32.KEYBDINPUT{
				WVk:         vk,
				WScan:       0,
				DwFlags:     flags,
				Time:        0,
				DwExtraInfo: 0,
			},
		},
	})
}

func setForegroundWindow(hwnd uintptr) bool {
	ret, _, _ := procSetForegroundWindow.Call(hwnd)
	return ret != 0
}

func findWindow(title string) (uintptr, error) {
	var hwnd syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		b := make([]uint16, 200)
		_, err := getWindowText(h, &b[0], int32(len(b)))
		if err != nil {
			return 1 // continue enumeration
		}
		if syscall.UTF16ToString(b) == title {
			// note the window
			hwnd = h
			return 0 // stop enumeration
		}
		return 1 // continue enumeration
	})
	enumWindows(cb, 0)
	if hwnd == 0 {
		return 0, fmt.Errorf("No window with title '%s' found", title)
	}
	return hwnd, nil
}

func enumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func getWindowText(hwnd syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
