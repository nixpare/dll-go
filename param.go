package main

import (
	"fmt"
	"strings"
)

// Param is function parameter
type Param struct {
	Name      string
	Type      string
	fn        *Fn
	tmpVarIdx int
}

// tmpVar returns temp variable name that will be used to represent p during syscall.
func (p *Param) tmpVar() string {
	if p.tmpVarIdx < 0 {
		p.tmpVarIdx = p.fn.curTmpVarIdx
		p.fn.curTmpVarIdx++
	}
	return fmt.Sprintf("_p%d", p.tmpVarIdx)
}

// TmpVarReadbackCode returns source code for reading back the temp variable into the original variable.
func (p *Param) TmpVarReadbackCode() string {
	switch {
	case p.Type == "*bool":
		return fmt.Sprintf("*%s = %s != 0", p.Name, p.tmpVar())
	default:
		return ""
	}
}

// TmpVarHelperCode returns source code for helper's temp variable.
func (p *Param) TmpVarHelperCode() string {
	return "var _ string"
}

// SyscallArgList returns source code fragments representing p parameter
// in syscall. Slices are translated into 2 syscall parameters: pointer to
// the first element and length.
func (p *Param) SyscallArg() string {
	var arg string
	if !strings.HasPrefix(p.Type, "*") {
		arg = "&"
	}
	arg += p.Name

	return "uintptr(unsafe.Pointer(" + arg + "))"
}

// IsError determines if p parameter is used to return error.
func (p *Param) IsError() bool {
	return p.Name == "err" && p.Type == "error"
}
