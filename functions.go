package main

import (
	"errors"
	"fmt"
	"go/token"
	"strings"
)

// Fn describes syscall function.
type Fn struct {
	Name        string
	Params      []*Param
	Rets        *Rets
	PrintTrace  bool
	dllname     string
	dllfuncname string
	src         string
	// TODO: get rid of this field and just use parameter index instead
	curTmpVarIdx int // insure tmp variables have uniq names
}

// extractParams parses s to extract function parameters.
func extractParams(s string, f *Fn) ([]*Param, error) {
	s = trim(s)
	if s == "" {
		return nil, nil
	}
	a := strings.Split(s, ",")
	ps := make([]*Param, len(a))
	for i := range ps {
		s2 := trim(a[i])
		b := strings.Split(s2, " ")
		if len(b) != 2 {
			b = strings.Split(s2, "\t")
			if len(b) != 2 {
				return nil, errors.New("Could not extract function parameter from \"" + s2 + "\"")
			}
		}
		ps[i] = &Param{
			Name:      trim(b[0]),
			Type:      trim(b[1]),
			fn:        f,
			tmpVarIdx: -1,
		}
	}
	return ps, nil
}

// extractSection extracts text out of string s starting after start
// and ending just before end. found return value will indicate success,
// and prefix, body and suffix will contain correspondent parts of string s.
func extractSection(s string, start, end rune) (prefix, body, suffix string, found bool) {
	s = trim(s)
	if strings.HasPrefix(s, string(start)) {
		// no prefix
		body = s[1:]
	} else {
		a := strings.SplitN(s, string(start), 2)
		if len(a) != 2 {
			return "", "", s, false
		}
		prefix = a[0]
		body = a[1]
	}
	a := strings.SplitN(body, string(end), 2)
	if len(a) != 2 {
		return "", "", "", false
	}
	return prefix, a[0], a[1], true
}

// newFn parses string s and return created function Fn.
func newFn(s string) (*Fn, error) {
	s = trim(s)
	f := &Fn{
		Rets:       new(Rets),
		src:        s,
		PrintTrace: *printTraceFlag,
	}
	// function name and args
	prefix, body, s, found := extractSection(s, '(', ')')
	if !found || prefix == "" {
		return nil, errors.New("Could not extract function name and parameters from \"" + f.src + "\"")
	}
	f.Name = prefix
	var err error
	f.Params, err = extractParams(body, f)
	if err != nil {
		return nil, err
	}
	// return values
	_, body, s, found = extractSection(s, '(', ')')
	if found {
		r, err := extractParams(body, f)
		if err != nil {
			return nil, err
		}

		if len(r) > 2 {
			return nil, errors.New("Too many return values in \"" + f.src + "\"")
		}

		for _, x := range r {
			*f.Rets = append(*f.Rets, ret{ Name: x.Name, Type: x.Type })
		}
	}
	
	// dll and dll function names
	s = trim(s)
	if s == "" {
		return f, nil
	}
	if !strings.HasPrefix(s, "=") {
		return nil, errors.New("Could not extract dll name from \"" + f.src + "\"")
	}
	s = trim(s[1:])
	if i := strings.LastIndex(s, "."); i >= 0 {
		f.dllname = s[:i]
		f.dllfuncname = s[i+1:]
	} else {
		f.dllfuncname = s
	}
	if f.dllfuncname == "" {
		return nil, fmt.Errorf("function name is not specified in %q", s)
	}
	
	return f, nil
}

// DLLName returns DLL name for function f.
func (f *Fn) DLLName() string {
	if f.dllname == "" {
		return "kernel32"
	}
	return f.dllname
}

// DLLVar returns a valid Go identifier that represents DLLName.
func (f *Fn) DLLVar() string {
	id := strings.Map(func(r rune) rune {
		switch r {
		case '.', '-':
			return '_'
		default:
			return r
		}
	}, f.DLLName()[strings.LastIndex(f.DLLName(), "/")+1:])

	if !token.IsIdentifier(id) {
		panic(fmt.Errorf("could not create Go identifier for DLLName %q", f.DLLName()))
	}
	return id
}

// DLLFuncName returns DLL function name for function f.
func (f *Fn) DLLFuncName() string {
	return f.Name + "DLL"
}

// ParamList returns source code for function f parameters.
func (f *Fn) ParamList() string {
	return join(f.Params, func(p *Param) string { return p.Name + " " + p.Type }, ", ")
}

// HelperParamList returns source code for helper function f parameters.
func (f *Fn) HelperParamList() string {
	a := make([]string, 0, len(f.Params))
	for _, p := range f.Params {
		var s string
		s += p.Name + " "
		if !strings.HasPrefix(p.Type, "*") {
			s += "*"
		}
		s += p.Type

		a = append(a, s)
	}
	return strings.Join(a, ", ")
}

// ParamPrintList returns source code of trace printing part correspondent
// to syscall input parameters.
func (f *Fn) ParamPrintList() string {
	return join(f.Params, func(p *Param) string { return fmt.Sprintf(`"%s=", %s, `, p.Name, p.Name) }, `", ", `)
}

// SyscallParamList returns source code for SyscallX parameters for function f.
func (f *Fn) SyscallParamList() string {
	a := make([]string, 0)
	for _, p := range f.Params {
		a = append(a, p.SyscallArg())
	}
	
	return strings.Join(a, ", ")
}

// HelperCallParamList returns source code of call into function f helper.
func (f *Fn) HelperCallParamList() string {
	a := make([]string, 0, len(f.Params))
	for _, p := range f.Params {
		s := p.Name
		if !strings.HasPrefix(p.Type, "*") {
			s = "*" + s
		}

		a = append(a, s)
	}
	return strings.Join(a, ", ")
}
