package main

import (
	"fmt"
	"strings"
)

// Rets describes function return parameters.
type Rets struct {
	Name          string
	Type          string
	ReturnsError  bool
	FailCond      string
	fnMaybeAbsent bool
}

// ErrorVarName returns error variable name for r.
func (r *Rets) ErrorVarName() string {
	if r.ReturnsError {
		return "err"
	}
	if r.Type == "error" {
		return r.Name
	}
	return ""
}

// ToParams converts r into slice of *Param.
func (r *Rets) ToParams() []*Param {
	ps := make([]*Param, 0)
	if len(r.Name) > 0 {
		ps = append(ps, &Param{Name: r.Name, Type: r.Type})
	}
	if r.ReturnsError {
		ps = append(ps, &Param{Name: "err", Type: "error"})
	}
	return ps
}

// List returns source code of syscall return parameters.
func (r *Rets) List() string {
	params := r.ToParams()
	if len(params) == 0 {
		return "(err error)"
	}

	if p := params[len(params)-1]; p.Type != "error" && p.Name != "err" {
		params = append(params, &Param{ Name: "err", Type: "error" })
	}
	
	s := join(params, func(p *Param) string { return p.Name + " " + p.Type }, ", ")
	return "(" + s + ")"
}

// FuncList returns source code of helper return parameters.
func (r *Rets) HelperList() string {
	params := r.ToParams()
	if len(params) == 0 {
		return ""
	}
	
	s := join(params, func(p *Param) string { return p.Name + " " + p.Type }, ", ")
	return "(" + s + ")"
}

// PrintList returns source code of trace printing part correspondent
// to syscall return values.
func (r *Rets) PrintList() string {
	return join(r.ToParams(), func(p *Param) string { return fmt.Sprintf(`"%s=", %s, `, p.Name, p.Name) }, `", ", `)
}

// SetReturnValuesCode returns source code that accepts syscall return values.
func (r *Rets) SetReturnValuesCode() string {
	if r.Name == "" && !r.ReturnsError {
		return "_, _, errno := "
	}
	retvar := "r0"
	if r.Name == "" {
		retvar = "r1"
	}
	
	return fmt.Sprintf("%s, _, errno := ", retvar)
}

func (r *Rets) useLongHandleErrorCode(retvar string) string {
	const code = `if %s {
		err = errnoErr(e1)
	}`
	cond := retvar + " == 0"
	if r.FailCond != "" {
		cond = strings.Replace(r.FailCond, "failretval", retvar, 1)
	}
	return fmt.Sprintf(code, cond)
}

// SetErrorCode returns source code that sets return parameters.
func (r *Rets) SetErrorCode() string {
	const code = `if r0 != 0 {
		%s = %sErrno(r0)
	}`
	const ntstatus = `if r0 != 0 {
		ntstatus = %sNTStatus(r0)
	}`
	const checkerrno = `if errno != windows.NOERROR {
		err = errno
	}`
	
	if r.Name == "" && !r.ReturnsError {
		return checkerrno
	}
	if r.Name == "" {
		return r.useLongHandleErrorCode("r1")
	}
	if r.Type == "error" && r.Name == "ntstatus" {
		return fmt.Sprintf(ntstatus, windowsdot())
	}
	if r.Type == "error" {
		return fmt.Sprintf(code, r.Name, syscalldot())
	}
	
	s := ""
	switch {
	case r.Type[0] == '*':
		s = fmt.Sprintf("%s = (%s)(unsafe.Pointer(r0))", r.Name, r.Type)
	case r.Type == "bool":
		s = fmt.Sprintf("%s = r0 != 0", r.Name)
	default:
		s = fmt.Sprintf("%s = %s(r0)", r.Name, r.Type)
	}

	return s + "\n" + checkerrno
	/* if !r.ReturnsError {
		return s
	}
	return s + "\n\t" + r.useLongHandleErrorCode(r.Name) */
}
