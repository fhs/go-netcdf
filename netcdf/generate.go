// Copyright 2014 The Go-NetCDF Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// +build ignore

// This tool generates nc_*.go files from nc_double.go
package main

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

// Selector represents a selector (e.g. "C.double" where X is "C" and Sel is "double)
type Selector struct {
	X, Sel string
}

type File struct {
	Name      string
	Idents    []string // identifiers
	DocIdents []string // documented identifiers
	Keys      []string
	Selectors []Selector
}

// TheFile serves as a template for generating the other files in OutFiles.
// The Selectors are filled in later from Idents.
var TheFile = File{
	Name: "nc_double.go",
	Idents: []string{
		"float64",
		"DOUBLE",
		"C.double",
		"C.nc_get_att_double",
		"C.nc_get_var_double",
		"C.nc_put_att_double",
		"C.nc_put_var_double",
	},
	DocIdents: []string{
		"testReadFloat64s",
		"testWriteFloat64s",
		"Float64sReader",
		"GetFloat64s",
		"ReadFloat64s",
		"WriteFloat64s",
	},
	Keys: []string{"float64", "Float64s", "DOUBLE", "C.double", "_double"},
}

// OutFiles are the files that needs to be generated from TheFile.
// Idents, DocIdents, and Selectors are filled in later based on TheFile.
//
// We return []byte (i.e. []uint8) for CHAR Type instead of a string because
// we want to be flexible: '\0' characters may or may not require trimming
// (e.g. if the author of the NetCDF file used CHAR type when UBYTE type
// should have been used).
var OutFiles = []File{
	{
		Name: "nc_uint64.go",
		Keys: []string{"uint64", "Uint64s", "UINT64", "C.ulonglong", "_ulonglong"},
	},
	{
		Name: "nc_int64.go",
		Keys: []string{"int64", "Int64s", "INT64", "C.longlong", "_longlong"},
	},
	{
		Name: "nc_uint.go",
		Keys: []string{"uint32", "Uint32s", "UINT", "C.uint", "_uint"},
	},
	{
		Name: "nc_int.go",
		Keys: []string{"int32", "Int32s", "INT", "C.int", "_int"},
	},
	{
		Name: "nc_float.go",
		Keys: []string{"float32", "Float32s", "FLOAT", "C.float", "_float"},
	},
	{
		Name: "nc_ushort.go",
		Keys: []string{"uint16", "Uint16s", "USHORT", "C.ushort", "_ushort"},
	},
	{
		Name: "nc_short.go",
		Keys: []string{"int16", "Int16s", "SHORT", "C.short", "_short"},
	},
	{
		Name: "nc_ubyte.go",
		Keys: []string{"uint8", "Uint8s", "UBYTE", "C.uchar", "_uchar"},
	},
	{
		Name: "nc_byte.go",
		Keys: []string{"int8", "Int8s", "BYTE", "C.schar", "_schar"},
	},
	{
		Name: "nc_char.go",
		Keys: []string{"byte", "Bytes", "CHAR", "C.char", "_text"},
	},
}

func init() {
	rename := func(old, kold, knew []string) []string {
		new := make([]string, len(old))
		for i, o := range old {
			for j := range kold {
				s := strings.Replace(o, kold[j], knew[j], 1)
				if s != o {
					new[i] = s
					break
				}
			}
		}
		return new
	}
	getSelectors := func(idents []string) []Selector {
		ss := make([]Selector, len(idents))
		k := 0
		for _, s := range idents {
			v := strings.Split(s, ".")
			if len(v) > 2 {
				log.Fatalf("invalid selector %s\n", s)
			}
			if len(v) == 2 {
				ss[k] = Selector{X: v[0], Sel: v[1]}
				k++
			}
		}
		return ss[:k]
	}

	TheFile.Selectors = getSelectors(TheFile.Idents)
	for i, of := range OutFiles {
		OutFiles[i].Idents = rename(TheFile.Idents, TheFile.Keys, of.Keys)
		OutFiles[i].DocIdents = rename(TheFile.DocIdents, TheFile.Keys, of.Keys)
		OutFiles[i].Selectors = getSelectors(OutFiles[i].Idents)
	}
}

func renameIdent(id *ast.Ident, old, new []string) {
	for i, o := range old {
		if id.Name == o {
			id.Name = new[i]
			break
		}
	}
}

func (f *File) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	// nc_put_att_text takes 5 arguments unlike nc_put_att_double which takes 6 arguments
	if f.Name == "nc_char.go" {
		if call, ok := n.(*ast.CallExpr); ok {
			if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
				if X, ok := fun.X.(*ast.Ident); ok && X.Name == "C" && fun.Sel.Name == "nc_put_att_double" {
					if len(call.Args) != 6 {
						log.Fatalf("C.nc_put_att_double call has wrong number of arguments")
					}
					call.Args[3] = call.Args[4]
					call.Args[4] = call.Args[5]
					call.Args = call.Args[:5]
				}
			}
		}
	}

	if id, ok := n.(*ast.Ident); ok {
		renameIdent(id, TheFile.Idents, f.Idents)
		renameIdent(id, TheFile.DocIdents, f.DocIdents)
	}
	if sel, ok := n.(*ast.SelectorExpr); ok {
		if X, ok := sel.X.(*ast.Ident); ok {
			sel := sel.Sel
			for i, o := range TheFile.Selectors {
				if X.Name == o.X && sel.Name == o.Sel {
					X.Name = f.Selectors[i].X
					sel.Name = f.Selectors[i].Sel
					break
				}
			}
		}
	}
	if c, ok := n.(*ast.Comment); ok {
		for i, o := range TheFile.DocIdents {
			s := strings.Replace(c.Text, o, f.DocIdents[i], 1)
			if s != c.Text {
				c.Text = s
				break
			}
		}
	}
	return f
}

func main() {
	for _, of := range OutFiles {
		// TODO: parse only once
		fset := token.NewFileSet()
		p, err := parser.ParseFile(fset, TheFile.Name, nil, parser.ParseComments)
		if err != nil {
			log.Fatalf("parsing %s failed: %v\n", TheFile.Name, err)
		}

		ast.Walk(&of, p)

		f, err := os.Create(of.Name)
		if err != nil {
			log.Fatalf("creating %s failed: %v\n", of.Name, err)
		}
		if err := format.Node(f, fset, p); err != nil {
			log.Fatalf("format.Node failed: %v\n", err)
		}
		f.Close()
	}
}
