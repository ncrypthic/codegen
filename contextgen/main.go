package main

import (
	"bytes"
	"flag"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fatih/astrewrite"
)

func isContextPackage(imports []string, pkg string) bool {
	for _, str := range imports {
		if str == pkg {
			return true
		}
	}
	return false
}

func FunctionVisitor(imports []*ast.ImportSpec) func(ast.Node) (ast.Node, bool) {
	contextPkgNames := []string{"context"}
	for _, imp := range imports {
		if imp.Path.Value == `"context"` && imp.Name != nil {
			contextPkgNames = append(contextPkgNames, imp.Name.Name)
		}
	}
	return func(node ast.Node) (ast.Node, bool) {
		funcNode, ok := node.(*ast.FuncType)
		if !ok {
			return node, true
		}
		if len(funcNode.Params.List) > 0 {
			for _, field := range funcNode.Params.List {
				typeIdent, ok := field.Type.(*ast.SelectorExpr)
				if !ok {
					continue
				}
				xIdent, ok := typeIdent.X.(*ast.Ident)
				if !ok {
					continue
				}
				if isContextPackage(contextPkgNames, xIdent.Name) && typeIdent.Sel.Name == "Context" {
					return node, true
				}
			}
		}
		startingPos := funcNode.Params.Opening
		ctxFieldObj := ast.NewObj(ast.Var, "ctx")
		ctxFieldName := ast.Ident{
			NamePos: startingPos,
			Name:    "ctx",
			Obj:     ctxFieldObj,
		}
		ctxFieldType := ast.Ident{
			NamePos: startingPos,
			Name:    contextPkgNames[0] + ".Context",
			Obj:     ctxFieldObj,
		}
		ctxField := ast.Field{
			Doc:     nil,
			Names:   []*ast.Ident{&ctxFieldName},
			Type:    &ctxFieldType,
			Tag:     nil,
			Comment: nil,
		}
		fieldList := make([]*ast.Field, 0)
		fieldList = append(fieldList, &ctxField)
		fieldList = append(fieldList, funcNode.Params.List...)
		funcNode.Params.List = fieldList
		return node, true
	}
}

func ImportContextPackage(astFile *ast.File) *ast.File {
	importSpec := &ast.ImportSpec{Path: &ast.BasicLit{Value: `"context"`}}
	var importDecl *ast.GenDecl
	for _, decl := range astFile.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok == token.IMPORT {
			importDecl = genDecl
			break
		}
	}
	newSpecs := make([]ast.Spec, len(importDecl.Specs)+1)
	for i, spec := range importDecl.Specs {
		specType, ok := spec.(*ast.ImportSpec)
		if !ok {
			continue
		}
		log.Print("import path:", specType.Path.Value)
		if specType.Path.Value == importSpec.Path.Value {
			// import "context" found
			return astFile
		}
		newSpecs[i] = spec
	}
	newSpecs[len(importDecl.Specs)] = importSpec
	importDecl.Specs = newSpecs
	if !importDecl.Lparen.IsValid() {
		importDecl.Lparen = token.Pos(1)
		importDecl.Rparen = token.Pos(2)
	}
	log.Printf("import decl: %#v", importDecl)
	log.Printf("import specs: %d", len(importDecl.Specs))
	return astFile
}

func ApplyContextToFunctions(workdir string, fileInfo os.FileInfo, rewrite, dryRun bool) error {
	fset := token.NewFileSet()
	bSrc, err := ioutil.ReadFile(strings.TrimRight(workdir, "/") + "/" + fileInfo.Name())
	if err != nil {
		return err
	}
	astFile, err := parser.ParseFile(fset, "", string(bSrc), parser.ParseComments)
	if err != nil {
		return err
	}
	astFile = ImportContextPackage(astFile)
	ast.SortImports(fset, astFile)
	res := astrewrite.Walk(astFile, FunctionVisitor(astFile.Imports))
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, res)
	if dryRun {
		log.Print(buf.String())
		return nil
	}
	outputFileName := strings.TrimRight(workdir, "/") + "/gen_" + fileInfo.Name()
	if rewrite {
		outputFileName = strings.TrimRight(workdir, "/") + "/" + fileInfo.Name()
	}

	return ioutil.WriteFile(outputFileName, buf.Bytes(), 0644)
}

func ProcessDirectory(dir string, recursive, rewrite, dryRun bool) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, fileInfo := range files {
		if fileInfo.IsDir() && recursive {
			ProcessDirectory(strings.TrimRight(dir, "/")+"/"+fileInfo.Name(), recursive, rewrite, dryRun)
		} else if err := ApplyContextToFunctions(dir, fileInfo, rewrite, dryRun); err != nil {
			log.Print("Error: ", err)
		}
	}
	return nil
}

func main() {
	/*
		1. Find all *.go file
		2. Foreach go file in files
		2.1 Parse the file ast
		2.2 Modify imports to inculde context if not exist
		2.3 Find all go function defs
		2.4 Modify all func defs
		2.5 Re-write the file
	*/
	if len(os.Args) < 2 {
		log.Fatal("Usage: " + os.Args[0] + " <directory>")
	}

	var (
		recursive bool
		rewrite   bool
		dryRun    bool
	)

	flag.BoolVar(&recursive, "R", false, "recursive")
	flag.BoolVar(&recursive, "recursive", false, "recursive")
	flag.BoolVar(&rewrite, "u", false, "rewrite existing file")
	flag.BoolVar(&rewrite, "rewrite", false, "rewrite existing file")
	flag.BoolVar(&dryRun, "n", false, "dry run, not writing to any file")
	flag.BoolVar(&dryRun, "dry", false, "dry run, not writing to any file")
	flag.Parse()
	var wg sync.WaitGroup
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	workdir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, dir := range flag.Args() {
		wg.Add(1)
		go func(dir string) {
			defer wg.Done()
			if err := ProcessDirectory(dir, recursive, rewrite, dryRun); err != nil {
				log.Printf("Error: %v", err)
			}
		}(workdir + "/" + dir)
	}
	wg.Wait()
}
