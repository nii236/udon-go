package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"udon-go/asm"
)

func main() {
	// src is the input for which we want to print the AST.

	f, err := os.Open("./assembly/udon_funcs_data.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer f.Close()
	uasm, err := asm.NewUdonAssembly(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	uc := &UdonCompiler{
		UASM: uasm,
	}

	srcFile, err := os.Open("./sample/sample.go")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer srcFile.Close()
	result, err := uc.MakeUASMCode(os.Stdout, srcFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(result)
}

type UdonCompiler struct {
	UASM                 *asm.UdonAssembly
	Node                 ast.Node
	CurrentFuncRetType   []*asm.UdonTypeName
	CurrentBreakLabel    []asm.LabelName
	CurrentContinueLabel []asm.LabelName
}

func (uc *UdonCompiler) MakeUASMCode(w io.Writer, rdr io.Reader) (string, error) {
	fset := token.NewFileSet() // positions are relative to fset

	f, err := parser.ParseFile(fset, "main.go", rdr, 0)
	if err != nil {
		return "", nil
	}

	collectFuncs(uc.UASM, f)
	uc.UASM.VarTable.AddVar(asm.VarName("ret_addr"), asm.UdonTypeUInt32, "0xFFFFFFFF")
	uc.UASM.VarTable.AddVar(asm.VarName("this_trans"), asm.UdonTypeTransform, "this")
	uc.UASM.VarTable.AddVar(asm.VarName("this_gameObj"), asm.UdonTypeGameObject, "this")

	err = handleDecls(uc.UASM, w, f)
	if err != nil {
		return "", fmt.Errorf("handle decls: %w", err)
	}

	retCode := ""
	dataSegment, err := uc.UASM.VarTable.MakeDataSeg()
	if err != nil {
		return "", fmt.Errorf("make data seg: %w", err)
	}
	retCode += dataSegment
	retCode += uc.UASM.MakeCodeSeg()
	// retCode = uc.UASM.ReplaceTmpAdrr(retCode)
	return retCode, nil
}

type Visitor struct {
	UASM *asm.UdonAssembly
}

func Str(in interface{}) string {
	return fmt.Sprintf("%s", in)
}

func collectFuncs(uasm *asm.UdonAssembly, fileNode *ast.File) {
	v := &Visitor{uasm}
	ast.Walk(v, fileNode)
}
func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	switch nt := node.(type) {
	case *ast.FuncDecl:
		argTypes := []asm.UdonTypeName{}
		retTypes := []asm.UdonTypeName{}
		argNames := []asm.VarName{}
		for _, arg := range nt.Type.Params.List {
			if len(arg.Names) > 1 {
				fmt.Println("multiple args to type not supported")
				os.Exit(1)
				return v
			}
			argTypes = append(argTypes, asm.UdonTypeName(fmt.Sprintf("%s", arg.Type)))
			argNames = append(argNames, asm.VarName(fmt.Sprintf("%s", arg.Names[0])))
		}
		for _, ret := range nt.Type.Results.List {
			retTypes = append(retTypes, asm.UdonTypeName(fmt.Sprintf("%s", ret.Type)))
		}

		if len(retTypes) > 1 {
			fmt.Println("multiple returns not supported")
			os.Exit(1)
			return v
		}
		v.UASM.FuncTable.Put(asm.FuncName(nt.Name.Name), argTypes, retTypes[0], argNames)
	}
	return v
}
