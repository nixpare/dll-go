package main

import (
	"fmt"
	"strings"
)

type ret struct {
	Name 	string
	Type 	string
}

// Rets describes function return parameters.
type Rets []ret

func (r *Rets) ToParams() []ret {
	params := make([]ret, len(*r))
	copy(params, *r)
	return params
}

func (r *Rets) RetStructList() string {
	a := make([]string, 0, len(*r))
	for i, x := range r.ToParams() {
		a = append(a, fmt.Sprintf("r%d %s", i, x.Type))
	}
	return strings.Join(a, "\n")
}

// List returns source code of syscall return parameters.
func (r *Rets) List() string {
	params := r.ToParams()

	if len(params) == 0 {
		return "(err error)"
	}

	if p := params[len(params)-1]; p.Type != "error" && p.Name != "err" {
		params = append(params, ret{ Name: "err", Type: "error" })
	}
	
	s := join(params, func(x ret) string { return x.Name + " " + x.Type }, ", ")
	return "(" + s + ")"
}

// FuncList returns source code of helper return parameters.
func (r *Rets) HelperList() string {
	params := r.ToParams()

	if len(params) == 0 {
		return ""
	}
	
	s := join(params, func(x ret) string { return x.Name + " " + x.Type }, ", ")
	return "(" + s + ")"
}

// SetReturnValuesCode returns source code that accepts syscall return values.
func (r *Rets) SetReturnValuesCode() string {
	params := r.ToParams()

	switch len(params) {
	case 0:
		return "_, _, errno := "
	default:
		return "r0, _, errno := "
	}
}

// SetErrorCode returns source code that sets return parameters.
func (r *Rets) SetErrorCode() string {
	const checkerrno = `if errno != windows.NOERROR {
		err = errno
		return
	}
	if r0 == 0 {
		err = syscall.Errno(1)
		return
	}`
	params := r.ToParams()

	if len(params) == 0 {
		return checkerrno
	}

	a := make([]string, 0, len(*r))
	for i, p := range params {
		a = append(a, fmt.Sprintf("%s = _r.r%d", p.Name, i))
	}

	return checkerrno + "\n" + strings.Join(a, "\n")
}
