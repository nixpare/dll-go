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
	modhello_world = syscall.NewLazyDLL("./hello-world.dll")

	procPrintDLL = modhello_world.NewProc("_PrintDLL")
)

type ret struct {
	r0 int
	r1 error
}

func PrintDLL(msg *message) (n int, err error) {
	u := new(ret)
	r0, _, errno := procPrintDLL.Call(uintptr(unsafe.Pointer(u)), uintptr(unsafe.Pointer(msg)))
	
	if errno != windows.NOERROR {
		err = errno
		return
	}
	if r0 == 0 {
		err = syscall.Errno(1)
		return
	}

	// procPrintDLLres := (*ret)(unsafe.Pointer(r0))
	n = u.r0
	err = u.r1
	return
}

//export _PrintDLL
func _PrintDLL(msg uintptr) (n int, err error) {
	return Print((*message)(unsafe.Pointer(msg)))
}
