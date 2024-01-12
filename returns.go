package main

import (
	"fmt"
	"log"
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
	
	s := join(params, func(x ret) string { return x.Name + " uintptr" }, ", ")
	return "(" + s + ")"
}

// SetReturnValuesCode returns source code that accepts syscall return values.
func (r *Rets) SetReturnValuesCode() string {
	params := r.ToParams()

	switch len(params) {
	case 0:
		return "_, _, errno := "
	case 1:
		return "r0, _, errno := "
	case 2:
		return "r0, r1, errno := "
	default:
		log.Fatalf(
			"can accept 2 return values max: got %d\n",
			len(params),
		)
		return "unreachable"
	}
}

// SetErrorCode returns source code that sets return parameters.
func (r *Rets) SetErrorCode() string {
	const checkerrno = `if errno != windows.NOERROR {
		err = errno
	}`
	params := r.ToParams()

	if len(params) == 0 {
		return checkerrno
	}

	a := make([]string, 0, len(params))
	for i, p := range params {
		s := fmt.Sprintf("%s)(unsafe.Pointer(r%d))", p.Type, i)
		if !strings.HasPrefix(p.Type, "*") {
			s = "*(*" + s
		} else {
			s = "(" + s
		}
		s = p.Name + " = " + s
		a = append(a, s)
	}

	return strings.Join(a, "\n") + "\n" + checkerrno
}
