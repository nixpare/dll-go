// Code generated by 'go generate'; DO NOT EDIT.

package main

import "C"
import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

var (
	modprint = syscall.NewLazyDLL("./print.dll")

	procPrintDLL = modprint.NewProc("_PrintDLL")
)

func PrintDLL(s string) (n int, err error) {
	r0, _, errno := syscall.SyscallN(procPrintDLL.Addr(), uintptr(unsafe.Pointer(&s)))
	n = int(r0)
	if errno != windows.NOERROR {
		err = errno
	}
	return
}

//export _PrintDLL
func _PrintDLL(s *string) (n int) {
	return Print(*s)
}
