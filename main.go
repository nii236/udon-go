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
	f, err := os.Open("./assembly/udon_funcs_data.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	defer f.Close()
	uasm, err := assembly.NewUdonAssembly(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	uc := &UdonCompiler{
		UASM: uasm,
	}
	err = uc.MakeUASMCode(os.Stdout, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type UdonCompiler struct {
	UASM                 *assembly.UdonAssembly
	Node                 ast.Node
	CurrentFuncRetType   []*assembly.UdonTypeName
	CurrentBreakLabel    []assembly.LabelName
	CurrentContinueLabel []assembly.LabelName
}

func (uc *UdonCompiler) MakeUASMCode(w io.Writer, in io.Reader) error {

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "main.go", in, 0)
	if err != nil {
		return fmt.Errorf("failed to parse file: %v", err)
	}

	// return address
	uc.UASM.VarTable.AddVar(assembly.VarName("ret_addr"), "UInt32", "0xFFFFFFFF")

	// this
	uc.UASM.VarTable.AddVar(assembly.VarName("this_trans"), "Transform", "this")
	uc.UASM.VarTable.AddVar(assembly.VarName("this_gameObj"), "GameObject", "this")

	// parse and eval AST
	// FORCE CAST, NO CHECK
	err = handleDecls(uc.UASM, w, f.Decls)
	if err != nil {
		return fmt.Errorf("handle decls: %w", err)
	}

	ret_code := ""
	dataSegment, err := uc.UASM.VarTable.MakeDataSeg()
	if err != nil {
		return fmt.Errorf("make data seg: %w", err)
	}
	ret_code += dataSegment
	ret_code += uc.UASM.MakeCodeSeg()
	// ret_code = uc.UASM.ReplaceTmpAdrr(ret_code)
	fmt.Println(ret_code)
	return nil
}

func (uc *UdonCompiler) PreCheckFuncDefs(body string) {}
