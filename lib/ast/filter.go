// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import (
	"lib/token"
	"sort"
)

// ----------------------------------------------------------------------------
// Export filtering

// exportFilter is a special filter function to extract exported nodes.
func exportFilter(name string) bool {
	return IsExported(name)
}

// FileExports trims the AST for a Go source file in place such that
// only exported nodes remain: all top-level identifiers which are not exported
// and their associated information (such as type, initial value, or function
// body) are removed. Non-exported fields and methods of exported types are
// stripped. The File.Comments list is not changed.
//
// FileExports reports whether there are exported declarations.
//
func FileExports(src *File) bool {
	return filterFile(src, exportFilter, true)
}

// PackageExports trims the AST for a Go package in place such that
// only exported nodes remain. The pkg.Files list is not changed, so that
// file names and top-level package comments don't get lost.
//
// PackageExports reports whether there are exported declarations;
// it returns false otherwise.
//
func PackageExports(pkg *Package) bool {
	return filterPackage(pkg, exportFilter, true)
}

// ----------------------------------------------------------------------------
// General filtering

type Filter func(string) bool

func filterIdentList(list []*Ident, f Filter) []*Ident {
	j := 0
	for _, x := range list {
		if f(x.Name) {
			list[j] = x
			j++
		}
	}
	return list[0:j]
}

// fieldName assumes that x is the type of an anonymous field and
// returns the corresponding field name. If x is not an acceptable
// anonymous field, the result is nil.
//
func fieldName(x Expr) *Ident {
	switch t := x.(type) {
	case *Ident:
		return t
	case *SelectorExpr:
		if _, ok := t.X.(*Ident); ok {
			return t.Sel
		}
	}
	return nil
}

func filterFieldList(fields *FieldList, filter Filter, export bool) (removedFields bool) {
	if fields == nil {
		return false
	}
	list := fields.List
	j := 0
	for _, f := range list {
		keepField := false
		if len(f.Names) == 0 {
			// anonymous field
			name := fieldName(f.Type)
			keepField = name != nil && filter(name.Name)
		} else {
			n := len(f.Names)
			f.Names = filterIdentList(f.Names, filter)
			if len(f.Names) < n {
				removedFields = true
			}
			keepField = len(f.Names) > 0
		}
		if keepField {
			if export {
				filterType(f.Type, filter, export)
			}
			list[j] = f
			j++
		}
	}
	if j < len(list) {
		removedFields = true
	}
	fields.List = list[0:j]
	return
}

func filterParamList(fields *FieldList, filter Filter, export bool) bool {
	if fields == nil {
		return false
	}
	var b bool
	for _, f := range fields.List {
		if filterType(f.Type, filter, export) {
			b = true
		}
	}
	return b
}

func filterType(typ Expr, f Filter, export bool) bool {
	switch t := typ.(type) {
	case *Ident:
		return f(t.Name)
	case *ParenExpr:
		return filterType(t.X, f, export)
	case *ArrayType:
		return filterType(t.Elt, f, export)
	case *StructType:
		if filterFieldList(t.Fields, f, export) {
			t.Incomplete = true
		}
		return len(t.Fields.List) > 0
	case *FuncType:
		b1 := filterParamList(t.Params, f, export)
		/*b2 := filterParamList(t.Results, f, export)*/
		return b1 /*|| b2*/
	case *MapType:
		b1 := filterType(t.Key, f, export)
		b2 := filterType(t.Value, f, export)
		return b1 || b2
	}
	return false
}





// FilterFile trims the AST for a Go file in place by removing all
// names from top-level declarations (including struct field and
// interface method names, but not from parameter lists) that don't
// pass through the filter f. If the declaration is empty afterwards,
// the declaration is removed from the AST. Import declarations are
// always removed. The File.Comments list is not changed.
//
// FilterFile reports whether there are any top-level declarations
// left after filtering.
//
func FilterFile(src *File, f Filter) bool {
	return filterFile(src, f, false)
}

func filterFile(src *File, f Filter, export bool) bool {
	j := 0
	return j > 0
}

// FilterPackage trims the AST for a Go package in place by removing
// all names from top-level declarations (including struct field and
// interface method names, but not from parameter lists) that don't
// pass through the filter f. If the declaration is empty afterwards,
// the declaration is removed from the AST. The pkg.Files list is not
// changed, so that file names and top-level package comments don't get
// lost.
//
// FilterPackage reports whether there are any top-level declarations
// left after filtering.
//
func FilterPackage(pkg *Package, f Filter) bool {
	return filterPackage(pkg, f, false)
}

func filterPackage(pkg *Package, f Filter, export bool) bool {
	hasDecls := false
	for _, src := range pkg.Files {
		if filterFile(src, f, export) {
			hasDecls = true
		}
	}
	return hasDecls
}

// ----------------------------------------------------------------------------
// Merging of package files

// The MergeMode flags control the behavior of MergePackageFiles.
type MergeMode uint

const (
	// If set, duplicate function declarations are excluded.
	FilterFuncDuplicates MergeMode = 1 << iota
	// If set, comments that are not associated with a specific
	// AST node (as Doc or Comment) are excluded.
	FilterUnassociatedComments
	// If set, duplicate import declarations are excluded.
	FilterImportDuplicates
)


// MergePackageFiles creates a file AST by merging the ASTs of the
// files belonging to a package. The mode flags control merging behavior.
//
func MergePackageFiles(pkg *Package, mode MergeMode) *File {
	// Count the number of package docs, comments and declarations across
	// all package files. Also, compute sorted list of filenames, so that
	// subsequent iterations can always iterate in the same order.
	ndecls := 0
	filenames := make([]string, len(pkg.Files))
	i := 0
	for filename, f := range pkg.Files {
		filenames[i] = filename
		i++
		ndecls += len(f.Decls)
	}
	sort.Strings(filenames)

	var pos token.Pos

	// Collect declarations from all package files.
	var decls []Decl

	// Collect import specs from all package files.
	var imports []*ImportSpec
	if mode&FilterImportDuplicates != 0 {
		seen := make(map[string]bool)
		for _, filename := range filenames {
			f := pkg.Files[filename]
			for _, imp := range f.Imports {
				if path := imp.Path.Value; !seen[path] {
					// TODO: consider handling cases where:
					// - 2 imports exist with the same import path but
					//   have different local names (one should probably
					//   keep both of them)
					// - 2 imports exist but only one has a comment
					// - 2 imports exist and they both have (possibly
					//   different) comments
					imports = append(imports, imp)
					seen[path] = true
				}
			}
		}
	} else {
		for _, f := range pkg.Files {
			imports = append(imports, f.Imports...)
		}
	}

	// TODO(gri) need to compute unresolved identifiers!
	return &File{pos, NewIdent(pkg.Name), decls, pkg.Scope, imports, nil, }
}
