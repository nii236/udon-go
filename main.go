package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"udon-go/assembly"
)

func main() {
	// src is the input for which we want to print the AST.

	f, err := os.Open("./assembly/udon_funcs_data.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer f.Close()
	uasm, err := assembly.NewUdonAssembly(f)
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
	UASM                 *assembly.UdonAssembly
	Node                 ast.Node
	CurrentFuncRetType   []*assembly.UdonTypeName
	CurrentBreakLabel    []assembly.LabelName
	CurrentContinueLabel []assembly.LabelName
}

func (uc *UdonCompiler) MakeUASMCode(w io.Writer, rdr io.Reader) (string, error) {
	fset := token.NewFileSet() // positions are relative to fset

	f, err := parser.ParseFile(fset, "main.go", rdr, 0)
	if err != nil {
		return "", nil
	}

	uc.UASM.VarTable.AddVar(assembly.VarName("ret_addr"), assembly.UdonTypeUInt32, "0xFFFFFFFF")
	uc.UASM.VarTable.AddVar(assembly.VarName("this_trans"), assembly.UdonTypeTransform, "this")
	uc.UASM.VarTable.AddVar(assembly.VarName("this_gameObj"), assembly.UdonTypeGameObject, "this")

	err = handleDecls(uc.UASM, w, f)
	if err != nil {
		return "", fmt.Errorf("handle decls: %w", err)
	}

	ret_code := ""
	dataSegment, err := uc.UASM.VarTable.MakeDataSeg()
	if err != nil {
		return "", fmt.Errorf("make data seg: %w", err)
	}
	ret_code += dataSegment
	ret_code += uc.UASM.MakeCodeSeg()
	// ret_code = uc.UASM.ReplaceTmpAdrr(ret_code)
	return ret_code, nil
}

func walk(v ast.Visitor, node ast.Node) {
	ast.Walk(v, node)
}

type Visitor struct {
	UASM *assembly.UdonAssembly
}

func Str(in interface{}) string {
	return fmt.Sprintf("%s", in)
}

// func (v *Visitor) Visit(node ast.Node) ast.Visitor {
// 	switch nt := node.(type) {
// 	case *ast.IfStmt:
// 	case *ast.BlockStmt:
// 	case *ast.AssignStmt:
// 		v.UASM.Assign(assembly.VarName(Str(nt.Lhs[0])), assembly.VarName(Str(nt.Rhs[0])))
// 	case *ast.BasicLit:
// 		v.UASM.VarTable.AddVarGlobal(assembly.VarName(nt.Value))
// 		v.UASM.VarTable.AddVar(assembly.VarName(nt.Value), assembly.UdonTypeUInt32, nt.Value)
// 	case *ast.ReturnStmt:
// 	case *ast.GenDecl:
// 		switch nt.Tok {
// 		case token.TYPE:
// 			for _, s := range nt.Specs {
// 				ts := s.(*ast.TypeSpec)
// 				switch ts.Type.(type) {
// 				case *ast.StructType:
// 				}
// 			}
// 		case token.VAR, token.CONST:
// 			for _, spec := range nt.Specs {
// 				vs := spec.(*ast.ValueSpec)
// 				if len(vs.Names) > 1 {
// 					fmt.Println("more than one assignment not supported")
// 					return v
// 				}
// 				switch l := vs.Values[0].(type) {
// 				case *ast.BasicLit:
// 					v.UASM.VarTable.AddVarGlobal(assembly.VarName(vs.Names[0].Name))
// 					v.UASM.VarTable.AddVar(assembly.VarName(vs.Names[0].Name), assembly.UdonTypeUInt32, l.Value)
// 				}
// 			}
// 		}
// 	case *ast.FuncDecl:
// 		argTypes := []assembly.UdonTypeName{}
// 		retTypes := []assembly.UdonTypeName{}
// 		argNames := []assembly.VarName{}
// 		for _, arg := range nt.Type.Params.List {
// 			if len(arg.Names) > 1 {
// 				fmt.Println("multiple args to type not supported")
// 				return v
// 			}
// 			argTypes = append(argTypes, assembly.UdonTypeName(fmt.Sprintf("%s", arg.Type)))
// 			argNames = append(argNames, assembly.VarName(fmt.Sprintf("%s", arg.Names[0])))
// 		}
// 		for _, ret := range nt.Type.Results.List {
// 			retTypes = append(retTypes, assembly.UdonTypeName(fmt.Sprintf("%s", ret.Type)))
// 		}

// 		if len(retTypes) > 1 {
// 			fmt.Println("multiple returns not supported")
// 			return v
// 		}
// 		v.UASM.DefFuncTable.AddFunc(assembly.FuncName(nt.Name.Name), argTypes, retTypes[0], argNames)
// 	}
// 	return v
// }
