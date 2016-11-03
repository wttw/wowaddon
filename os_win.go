// +build windows

package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

var installDirs = []string{
	`C:\Program Files\World of Warcraft`,
	`C:\Program Files (x86)\World of Warcraft`,
	`C:\Users\Public\Games\World of Warcraft`,
	`D:\Program Files\World of Warcraft`,
}

const (
	enableProcessedOutput           uintptr = 0x01
	enableVirtualTerminalProcessing         = 0x04
)

const (
	errorInvalidHandle uintptr = 0x06
)

func EnableColor() error {
	// This obtains a handle for the kernel32 DLL and then the SetConsoleMode function from that DLL. This function
	// is described here: https://msdn.microsoft.com/en-us/library/windows/desktop/ms686033(v=vs.85).aspx
	dll := syscall.MustLoadDLL("kernel32")
	setConsoleMode := dll.MustFindProc("SetConsoleMode")
	if os.Stdout == nil {
		return errors.New("stdout is nil")
	}
	handle := syscall.Handle(os.Stdout.Fd())

	// Call the obtained SetConsoleMode, setting ENABLE_PROCESSED_OUTPUT (was set anyway) and
	// ENABLE_VIRTUAL_TERMINAL_PROCESSING (the important one)
	r1, _, err := setConsoleMode.Call(uintptr(handle), uintptr(enableProcessedOutput|enableVirtualTerminalProcessing))
	if r1 == 0 {
		errNo := "(unknown)"
		if en, ok := err.(syscall.Errno); ok {
			if uintptr(en) == errorInvalidHandle {
				return nil
			}
			errNo = fmt.Sprintf("%x", int(en))
		}
		return fmt.Errorf("setting console mode: %s -- %s", errNo, err)
	}

	return nil
}
