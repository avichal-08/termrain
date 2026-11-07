//go:build windows

package main

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
	procGetCSBI        = kernel32.NewProc("GetConsoleScreenBufferInfo")
)

const (
	enableProcessedInput  = 0x0001
	enableLineInput       = 0x0002
	enableEchoInput       = 0x0004
	enableVTInput         = 0x0200
	enableProcessedOutput = 0x0001
	enableVTProcessing    = 0x0004
)

func setConsoleMode(h syscall.Handle, mode uint32) error {
	r, _, err := procSetConsoleMode.Call(uintptr(h), uintptr(mode))
	if r == 0 {
		return err
	}
	return nil
}

func enableRaw() (func(), error) {
	in := syscall.Handle(os.Stdin.Fd())
	out := syscall.Handle(os.Stdout.Fd())
	var inMode, outMode uint32
	if err := syscall.GetConsoleMode(in, &inMode); err != nil {
		return nil, err
	}
	if err := syscall.GetConsoleMode(out, &outMode); err != nil {
		return nil, err
	}
	raw := inMode&^(enableProcessedInput|enableLineInput|enableEchoInput) | enableVTInput
	if err := setConsoleMode(in, raw); err != nil {
		return nil, err
	}
	if err := setConsoleMode(out, outMode|enableProcessedOutput|enableVTProcessing); err != nil {
		setConsoleMode(in, inMode)
		return nil, err
	}
	return func() {
		setConsoleMode(in, inMode)
		setConsoleMode(out, outMode)
	}, nil
}

type coord struct{ X, Y int16 }
type smallRect struct{ L, T, R, B int16 }
type screenBufferInfo struct {
	Size    coord
	Cursor  coord
	Attr    uint16
	Window  smallRect
	MaxSize coord
}

func termSize() (int, int) {
	var info screenBufferInfo
	r, _, _ := procGetCSBI.Call(os.Stdout.Fd(), uintptr(unsafe.Pointer(&info)))
	if r == 0 {
		return 80, 24
	}
	return int(info.Window.R-info.Window.L) + 1, int(info.Window.B-info.Window.T) + 1
}
