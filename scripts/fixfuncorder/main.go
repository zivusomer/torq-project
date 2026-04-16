// fixfuncorder moves exported top-level functions before the first unexported
// top-level function in each file, matching the function-order lint rule.
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	var changedFiles int
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		changed, err := processFile(path)
		if err != nil {
			return err
		}
		if changed {
			changedFiles++
			fmt.Fprintf(os.Stderr, "fixed function order: %s\n", path)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if changedFiles > 0 {
		fmt.Fprintf(os.Stderr, "fixfuncorder: updated %d file(s)\n", changedFiles)
	}
}

func needsReorder(decls []ast.Decl) bool {
	seenPrivate := false
	for _, d := range decls {
		fd, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fd.Name.IsExported() && seenPrivate {
			return true
		}
		if !fd.Name.IsExported() {
			seenPrivate = true
		}
	}
	return false
}

func reorderFuncDecls(decls []ast.Decl) []ast.Decl {
	if !needsReorder(decls) {
		return decls
	}
	firstPriv := -1
	for i, d := range decls {
		fd, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if !fd.Name.IsExported() {
			firstPriv = i
			break
		}
	}
	if firstPriv < 0 {
		return decls
	}
	var move []*ast.FuncDecl
	skip := make(map[int]bool)
	for i := firstPriv + 1; i < len(decls); i++ {
		fd, ok := decls[i].(*ast.FuncDecl)
		if ok && fd.Name.IsExported() {
			move = append(move, fd)
			skip[i] = true
		}
	}
	if len(move) == 0 {
		return decls
	}
	out := make([]ast.Decl, 0, len(decls))
	for i := 0; i < firstPriv; i++ {
		out = append(out, decls[i])
	}
	for _, fd := range move {
		out = append(out, fd)
	}
	for i := firstPriv; i < len(decls); i++ {
		if skip[i] {
			continue
		}
		out = append(out, decls[i])
	}
	return out
}

func processFile(path string) (changed bool, err error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return false, fmt.Errorf("%s: %w", path, err)
	}
	if !needsReorder(node.Decls) {
		return false, nil
	}
	node.Decls = reorderFuncDecls(node.Decls)
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return false, fmt.Errorf("%s: format: %w", path, err)
	}
	newSrc := buf.Bytes()
	if bytes.Equal(src, newSrc) {
		return false, nil
	}
	if err := os.WriteFile(path, newSrc, 0o644); err != nil {
		return false, err
	}
	return true, nil
}
