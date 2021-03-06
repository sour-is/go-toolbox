package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
	"unicode"

	"sour.is/x/toolbox/log"
)

type model struct {
	Types []structMap
}

type structMap struct {
	Name   string
	ID     string
	Fields []structField
	Table  bool
	View   bool
	ROnly  bool
	HasID  bool
}

type structField struct {
	Name      string
	Container string
	Auto      bool
	ROnly     bool
}

func init() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: gogen CONFIG MODEL")
	}
}

func main() {
	fset := token.NewFileSet() // positions are relative to fset

	// Parse src but stop after processing the imports.
	f, err := parser.ParseFile(fset, os.Args[2], nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	var m model

	for _, d := range f.Decls {
		var ok bool
		var gen *ast.GenDecl
		if gen, ok = d.(*ast.GenDecl); !ok {
			continue
		}

		for _, spec := range gen.Specs {

			var ts *ast.TypeSpec
			if ts, ok = spec.(*ast.TypeSpec); !ok {
				continue
			}

			smap := structMap{}
			smap.Name = ts.Name.Name
			log.Debugf("type %s\n", ts.Name.Name)

			var st *ast.StructType
			if st, ok = ts.Type.(*ast.StructType); !ok {
				continue
			}

			smap.Table, smap.View, smap.Fields, smap.HasID, smap.ID = procStruct(st)

			if !smap.Table && !smap.View {
				// No Table or View == Do not generate.
				continue
			}
			if !smap.Table && smap.View {
				// View only == Only generate read items.
				smap.ROnly = true
			}

			m.Types = append(m.Types, smap)
		}
	}

	fh, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	basepath := filepath.Dir(os.Args[1])
	err = os.Chdir(basepath)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		fn := scanner.Text()
		log.Print(fn)
		tpl, err := template.New("tpl").Funcs(fnMap).ParseFiles(fn)
		if err != nil {
			log.Fatal(err)
		}

		w, err := os.Create(fn + ".go")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprint(w, "// Code generated by sour.is/x/toolbox/cmd/modelgen, DO NOT EDIT.\n\n")
		err = tpl.ExecuteTemplate(w, tpl.Templates()[0].ParseName, m)
		if err != nil {
			log.Error(err)
		}
	}
}

func procStruct(st *ast.StructType) (table, view bool, fields []structField, hasID bool, ID string) {
	for _, field := range st.Fields.List {
		for _, name := range field.Names {
			if name.Name == "ID" {
				ID = name.Name
				hasID = true
			}

			tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))

			auto := false
			if db := tag.Get("db"); strings.Contains(db, ",AUTO") {
				auto = true
			}

			readonly := false
			if db := tag.Get("db"); strings.Contains(db, ",RO") {
				readonly = true
			}

			if db := tag.Get("db"); strings.Contains(db, ",ID") {
				hasID = true
				ID = name.Name
			}

			if tab := tag.Get("table"); tab != "" {
				table = true
			}
			if tab := tag.Get("view"); tab != "" {
				view = true
			}

			if arr, ok := field.Type.(*ast.ArrayType); ok {
				if _, ok = arr.Elt.(*ast.Ident); !ok {
					log.Debug("non-ident array ", name, field)
					continue
				}

				fields = append(
					fields,
					structField{
						Name:      name.Name,
						Container: tag.Get("cont"),
						Auto:      auto,
						ROnly:     readonly,
					})
				continue
			}

			fields = append(
				fields,
				structField{
					Name:      name.Name,
					Container: "",
					Auto:      auto,
					ROnly:     readonly,
				})
		}
	}

	return
}

var fnMap = template.FuncMap{
	"snakeCase": SnakeCase,
	"spineCase": SpineCase,
}

// SnakeCase converts the given string to snake case following the Golang format:
// acronyms are converted to lower-case and preceded by an underscore.
func SnakeCase(s string) string {
	in := []rune(s)
	isLower := func(idx int) bool {
		return idx >= 0 && idx < len(in) && unicode.IsLower(in[idx])
	}

	out := make([]rune, 0, len(in)+len(in)/2)
	for i, r := range in {
		if unicode.IsUpper(r) {
			r = unicode.ToLower(r)
			if i > 0 && in[i-1] != '_' && (isLower(i-1) || isLower(i+1)) {
				out = append(out, '_')
			}
		}
		out = append(out, r)
	}

	return string(out)
}

// SpineCase converts the given string to snake case following the Golang format:
// acronyms are converted to lower-case and preceded by an underscore.
func SpineCase(s string) string {
	in := []rune(s)
	isLower := func(idx int) bool {
		return idx >= 0 && idx < len(in) && unicode.IsLower(in[idx])
	}

	out := make([]rune, 0, len(in)+len(in)/2)
	for i, r := range in {
		if unicode.IsUpper(r) {
			r = unicode.ToLower(r)
			if i > 0 && in[i-1] != '-' && (isLower(i-1) || isLower(i+1)) {
				out = append(out, '-')
			}
		}
		out = append(out, r)
	}

	return string(out)
}
