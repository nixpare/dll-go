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

// BoolTmpVarCode returns source code for bool temp variable.
func (p *Param) BoolTmpVarCode() string {
	const code = `var %[1]s uint32
	if %[2]s {
		%[1]s = 1
	}`
	return fmt.Sprintf(code, p.tmpVar(), p.Name)
}

// BoolPointerTmpVarCode returns source code for bool temp variable.
func (p *Param) BoolPointerTmpVarCode() string {
	const code = `var %[1]s uint32
	if *%[2]s {
		%[1]s = 1
	}`
	return fmt.Sprintf(code, p.tmpVar(), p.Name)
}

// SliceTmpVarCode returns source code for slice temp variable.
func (p *Param) SliceTmpVarCode() string {
	const code = `var %s *%s
	if len(%s) > 0 {
		%s = &%s[0]
	}`
	tmp := p.tmpVar()
	return fmt.Sprintf(code, tmp, p.Type[2:], p.Name, tmp, p.Name)
}

// StringTmpVarCode returns source code for string temp variable.
func (p *Param) StringTmpVarCode() string {
	errvar := p.fn.Rets.ErrorVarName()
	if errvar == "" {
		errvar = "err"
	}

	tmp := p.tmpVar()
	const code = `var %s %s
	%s, %s = %s(%s)`
	s := fmt.Sprintf(code, tmp, p.fn.StrconvType(), tmp, errvar, p.fn.StrconvFunc(), p.Name)
	if errvar == "-" {
		return s
	}
	const morecode = `
	if %s != nil {
		return
	}`
	return s + fmt.Sprintf(morecode, errvar)
}

// TmpVarCode returns source code for temp variable.
func (p *Param) TmpVarCode() string {
	switch {
	case p.Type == "bool":
		return p.BoolTmpVarCode()
	case p.Type == "*bool":
		return p.BoolPointerTmpVarCode()
	case strings.HasPrefix(p.Type, "[]"):
		return p.SliceTmpVarCode()
	default:
		return ""
	}
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
