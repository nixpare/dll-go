package main

import (
	"fmt"
	"log"
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
		s := ""
		switch {
		case p.Type[0] == '*':
			s = fmt.Sprintf("%s = (%s)(unsafe.Pointer(r%d))", p.Name, p.Type, i)
		case p.Type == "bool":
			s = fmt.Sprintf("%s = r%d != 0", p.Name, i)
		default:
			s = fmt.Sprintf("%s = %s(r%d)", p.Name, p.Type, i)
		}
		a = append(a, s)
	}

	return strings.Join(a, "\n") + "\n" + checkerrno
}
