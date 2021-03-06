package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"
	"udon-go/asm"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	// src is the input for which we want to print the AST.

	f, err := os.Open("./asm/udon_funcs_data.txt")
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

	srcFile, err := os.Open("./sample/func.go")
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
	err = uc.UASM.VarTable.AddVar(asm.VarName("ret_addr"), asm.UdonTypeUInt32, "0xFFFFFFFF")
	if err != nil {
		return "", fmt.Errorf("add init vars: %w", err)
	}
	err = uc.UASM.VarTable.AddVar(asm.VarName("this_trans"), asm.UdonTypeTransform, "this")
	if err != nil {
		return "", fmt.Errorf("add init vars: %w", err)
	}
	err = uc.UASM.VarTable.AddVar(asm.VarName("this_gameObj"), asm.UdonTypeGameObject, "this")
	if err != nil {
		return "", fmt.Errorf("add init vars: %w", err)
	}

	c := &Compiler{}
	err = c.handleDecls(uc.UASM, w, f)
	if err != nil {
		return "", fmt.Errorf("handle decls: %w", err)
	}
	spew.Dump(uc.UASM.VarTable.VarDict)
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

// Visit each node for func collection
func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	var err error
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
		if nt.Type.Results != nil && len(nt.Type.Results.List) > 0 {
			for _, ret := range nt.Type.Results.List {
				if ret.Tag == nil {
					break
				}
				retTypes = append(retTypes, TokenToUnity(ret.Tag.Kind))
			}
		}

		if len(retTypes) > 1 {
			fmt.Println("multiple returns not supported")
			os.Exit(1)
			return v
		}

		if nt.Name.Name == "main" {
			err = v.UASM.AddEvent(asm.EventName("_start"), argNames, argTypes)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if strings.HasPrefix(nt.Name.Name, "_") {
			err = v.UASM.AddEvent(asm.EventName(nt.Name.Name), argNames, argTypes)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			udonReturnType := asm.GoNil
			if len(retTypes) > 0 {
				udonReturnType = retTypes[0]
			}
			v.UASM.FuncTable.Put(asm.FuncName(nt.Name.Name), argTypes, udonReturnType, argNames)
		}

	}
	return v
}

func TokenToUnity(kind token.Token) asm.UdonTypeName {
	switch kind {
	case token.INT:
		return asm.GoInt

	case token.FLOAT:
		return asm.GoFloat32

	case token.CHAR:
		return asm.GoRune

	case token.STRING:
		return asm.GoString
	}
	panic("bad token")
}

func IdentToUnity(kind *ast.Ident) asm.UdonTypeName {
	switch kind.Name {
	case "int":
		return asm.GoInt

	case "float":
		return asm.GoFloat32

	case "char":
		return asm.GoRune

	case "string":
		return asm.GoString
	}
	panic("bad token")
}
