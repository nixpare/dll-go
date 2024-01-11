package main

import (
	"bufio"
	//"errors"
	"go/parser"
	"go/token"
	"os"
	//"path/filepath"
	//"runtime"
	"sort"
	"strings"
)

// DLL is a DLL's filename and a string that is valid in a Go identifier that should be used when
// naming a variable that refers to the DLL.
type DLL struct {
	Name string
	Var  string
}

// Source files and functions.
type Source struct {
	Funcs           []*Fn
	DLLFuncNames    []*Fn
	Files           []string
	StdLibImports   []string
	ExternalImports []string
	HasCGO          bool
}

func (src *Source) Import(pkg string) {
	src.StdLibImports = append(src.StdLibImports, pkg)
	sort.Strings(src.StdLibImports)
}

func (src *Source) ExternalImport(pkg string) {
	if pkg == "golang.org/x/sys/windows" {
		return
	}

	src.ExternalImports = append(src.ExternalImports, pkg)
	sort.Strings(src.ExternalImports)
}

// ParseFiles parses files listed in fs and extracts all syscall
// functions listed in sys comments. It returns source files
// and functions collection *Source if successful.
func ParseFiles(fs []string) (*Source, error) {
	src := &Source{
		Funcs: make([]*Fn, 0),
		Files: make([]string, 0),
		StdLibImports: []string{ "syscall", "unsafe" },
		ExternalImports: []string{ "golang.org/x/sys/windows" },
	}
	for _, file := range fs {
		if err := src.ParseFile(file); err != nil {
			return nil, err
		}
	}
	src.DLLFuncNames = make([]*Fn, 0, len(src.Funcs))
	uniq := make(map[string]bool, len(src.Funcs))
	for _, fn := range src.Funcs {
		name := fn.DLLFuncName()
		if !uniq[name] {
			src.DLLFuncNames = append(src.DLLFuncNames, fn)
			uniq[name] = true
		}
	}	

	return src, nil
}

// DLLs return dll names for a source set src.
func (src *Source) DLLs() []DLL {
	uniq := make(map[string]bool)
	r := make([]DLL, 0)
	for _, f := range src.Funcs {
		id := f.DLLVar()
		if _, found := uniq[id]; !found {
			uniq[id] = true
			r = append(r, DLL{f.DLLName(), id})
		}
	}
	sort.Slice(r, func(i, j int) bool {
		return r[i].Var < r[j].Var
	})
	return r
}

// ParseFile adds additional file path to a source set src.
func (src *Source) ParseFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	for s.Scan() {
		t := trim(s.Text())
		if len(t) < 7 {
			continue
		}
		if !strings.HasPrefix(t, "//sys") {
			continue
		}
		t = t[5:]
		if !(t[0] == ' ' || t[0] == '\t') {
			continue
		}
		f, err := newFn(t[1:])
		if err != nil {
			return err
		}
		src.Funcs = append(src.Funcs, f)
	}
	if err := s.Err(); err != nil {
		return err
	}
	src.Files = append(src.Files, path)
	sort.Slice(src.Funcs, func(i, j int) bool {
		fi, fj := src.Funcs[i], src.Funcs[j]
		if fi.DLLName() == fj.DLLName() {
			return fi.DLLFuncName() < fj.DLLFuncName()
		}
		return fi.DLLName() < fj.DLLName()
	})

	// get package name
	fset := token.NewFileSet()
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	pkg, err := parser.ParseFile(fset, "", file, parser.PackageClauseOnly)
	if err != nil {
		return err
	}
	packageName = pkg.Name.Name

	return nil
}

// IsStdRepo reports whether src is part of standard library.
/* func (src *Source) IsStdRepo() (bool, error) {
	if len(src.Files) == 0 {
		return false, errors.New("no input files provided")
	}
	abspath, err := filepath.Abs(src.Files[0])
	if err != nil {
		return false, err
	}
	goroot := runtime.GOROOT()
	if runtime.GOOS == "windows" {
		abspath = strings.ToLower(abspath)
		goroot = strings.ToLower(goroot)
	}
	sep := string(os.PathSeparator)
	if !strings.HasSuffix(goroot, sep) {
		goroot += sep
	}
	return strings.HasPrefix(abspath, goroot), nil
} */