package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"udon-go/assembly"
)

func main() {
	// src is the input for which we want to print the AST.

	b, err := ioutil.ReadFile("./sample/sample.go")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	uc := &UdonCompiler{}
	uc.MakeUASMCode(os.Stdout, bytes.NewReader(b))
}

type UdonCompiler struct {
	VarTable             assembly.VarTable
	UASM                 assembly.UdonAssembly
	DefFuncTable         assembly.DefFuncTable
	UdonMethodTable      assembly.UdonMethodTable
	Node                 ast.Node
	CurrentFuncRetType   []*assembly.UdonTypeName
	CurrentBreakLabel    []assembly.LabelName
	CurrentContinueLabel []assembly.LabelName
}

func (uc *UdonCompiler) MakeUASMCode(out io.Writer, in io.Reader) (string, error) {

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "main.go", in, 0)
	if err != nil {
		return "", fmt.Errorf("failed to parse file: %v", err)
	}

	// return address
	uc.VarTable.AddVar(assembly.VarName("ret_addr"), &assembly.UdonTypeName{Name: "UInt32", InitValue: "0"}, "0xFFFFFFFF")

	// this
	uc.VarTable.AddVar(assembly.VarName("this_trans"), &assembly.UdonTypeName{Name: "Transform", InitValue: "this"}, "this")
	uc.VarTable.AddVar(assembly.VarName("this_gameObj"), &assembly.UdonTypeName{Name: "GameObject", InitValue: "this"}, "this")

	// parse and eval AST
	// FORCE CAST, NO CHECK
	handleDecls(uc.UASM, out, f.Decls)

	ret_code := ""
	dataSegment, err := uc.VarTable.MakeDataSeg()
	if err != nil {
		return "", fmt.Errorf("make data seg: %w", err)
	}
	ret_code += dataSegment
	ret_code += uc.UASM.MakeCodeSeg()
	ret_code = uc.UASM.ReplaceTmpAdrr(ret_code)
	return ret_code, nil
}

func (uc *UdonCompiler) PreCheckFuncDefs(body string) {}
