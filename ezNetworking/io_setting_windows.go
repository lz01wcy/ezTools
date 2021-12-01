//go:build windows
// +build windows

package ezNetworking

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/windows"
)

func setIOBufferSize(c syscall.RawConn, bufferSize int, opt int) *EZNetError {
	rsErr := &EZNetError{}
	err := c.Control(func(fd uintptr) {
		if err := windows.SetsockoptInt(windows.Handle(fd), syscall.SOL_SOCKET, opt, bufferSize); err != nil {
			rsErr.IsSettingError = true
			rsErr.errorDesc = err.Error()
		}
	})
	if err != nil {
		rsErr.IsSettingError = true
		rsErr.errorDesc = err.Error()
		return rsErr
	}
	if rsErr.errorDesc != "" {
		return rsErr
	}
	err = c.Control(func(fd uintptr) {
		if actuallyValue, err := windows.GetsockoptInt(windows.Handle(fd), syscall.SOL_SOCKET, opt); err != nil {
			rsErr.IsLoadingError = true
			rsErr.errorDesc = err.Error()
			return
		} else if actuallyValue < bufferSize {
			rsErr.IsCheckingError = true
			rsErr.errorDesc = fmt.Sprintf("set buffer failed,wanted result:%d,actually:%d", bufferSize, actuallyValue)
			return
		}
	})
	if err != nil {
		rsErr.IsLoadingError = true
		rsErr.errorDesc = err.Error()
		return rsErr
	}
	if rsErr.errorDesc != "" {
		return rsErr
	}
	return nil
}
