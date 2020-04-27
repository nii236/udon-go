package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"strconv"
	"udon-go/assembly"
)

func handleDecls(uasm *assembly.UdonAssembly, out io.Writer, d *ast.File) error {
	// fmt.Println("run: handleDecls")
	var err error
	for _, decl := range d.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			err = handleGenDecl(uasm, out, decl)
			if err != nil {
				return fmt.Errorf("handle generic declaration: %w", err)
			}
		case *ast.FuncDecl:
			err = handleFuncDecl(uasm, out, decl)
			if err != nil {
				return fmt.Errorf("handle func declaration: %w", err)
			}
		default:
			fmt.Printf("unsupported decl: %#v", d)
		}
	}
	return nil
}

func handleGenDecl(uasm *assembly.UdonAssembly, out io.Writer, decl *ast.GenDecl) error {
	// fmt.Println("run: handleGenDecl")
	for _, s := range decl.Specs {
		vs, ok := s.(*ast.ValueSpec)
		if !ok {
			return nil
		}
		if len(vs.Names) > 1 {
			return fmt.Errorf("unsupported # of value names: %v", vs.Names)
		}
		if len(vs.Values) > 1 {
			return fmt.Errorf("unsupported # of values: %v", vs.Names)
		}

		switch l := vs.Values[0].(type) {
		case *ast.BasicLit:
			varName := assembly.VarName(vs.Names[0].Name)
			switch l.Kind {
			case token.INT:
				val, err := strconv.Atoi(l.Value)
				if err != nil {
					return fmt.Errorf("parse int: %w", err)
				}
				uasm.SetUint32(varName, val)
				uasm.VarTable.AddVar(varName, assembly.UdonTypeInt32, l.Value)
			case token.FLOAT:
				return fmt.Errorf("unsupported token: %s", l.Kind.String())
			case token.IMAG:
				return fmt.Errorf("unsupported token: %s", l.Kind.String())
			case token.CHAR:
				return fmt.Errorf("unsupported token: %s", l.Kind.String())
			case token.STRING:
				return fmt.Errorf("unsupported token: %s", l.Kind.String())
			default:
				return fmt.Errorf("unsupported token: %s", l.Kind.String())
			}

			uasm.VarTable.AddVarGlobal(varName)
		}

	}
	return nil
}
func handleFuncDecl(uasm *assembly.UdonAssembly, out io.Writer, decl *ast.FuncDecl) error {
	// fmt.Println("run: handleFuncDecl")
	argTypes := []assembly.UdonTypeName{}
	retTypes := []assembly.UdonTypeName{}
	argNames := []assembly.VarName{}

	for _, arg := range decl.Type.Params.List {
		if len(arg.Names) > 1 {
			return errors.New("multiple args to type not supported")
		}
		argTypes = append(argTypes, assembly.UdonTypeName(fmt.Sprintf("%s", arg.Type)))
		argNames = append(argNames, assembly.VarName(fmt.Sprintf("%s", arg.Names[0])))
	}
	for _, ret := range decl.Type.Results.List {
		retTypes = append(retTypes, assembly.UdonTypeName(fmt.Sprintf("%s", ret.Type)))
	}

	if len(retTypes) > 1 {
		return errors.New("multiple returns not supported")
	}
	uasm.DefFuncTable.AddFunc(assembly.FuncName(decl.Name.Name), argTypes, retTypes[0], argNames)

	err := handleBlockStmt(uasm, out, decl.Body)
	if err != nil {
		return fmt.Errorf("handle block: %w", err)
	}
	return nil
}
func handleBlockStmt(uasm *assembly.UdonAssembly, out io.Writer, bs *ast.BlockStmt) error {

	for _, s := range bs.List {
		switch st := s.(type) {
		case *ast.ExprStmt:
			// fmt.Println("handle ast.ExprStmt")
			_, err := handleExpr(uasm, out, st.X)
			if err != nil {
				return fmt.Errorf("error handling expr: %v", err)
			}
		case *ast.AssignStmt:
			// fmt.Println("handle ast.AssignStmt")
			if len(st.Lhs) > 1 {
				return fmt.Errorf("assign: unsupported # of lhs exprs: %v", st.Lhs)
			}
			if len(st.Rhs) > 1 {
				return fmt.Errorf("assign: unsupported # of rhs exprs: %v", st.Rhs)
			}

			lhs := st.Lhs[0]
			rhs := st.Rhs[0]

			lhsVarName, err := handleExpr(uasm, out, lhs)
			if err != nil {
				return fmt.Errorf("assign: left expr %v: %v", lhs, err)
			}
			rhsVarName, err := handleExpr(uasm, out, rhs)
			if err != nil {
				return fmt.Errorf("assign: right expr %v: %v", rhs, err)
			}
			uasm.Assign(lhsVarName, rhsVarName)
		case *ast.ReturnStmt:
			err := uasm.PopVar(assembly.VarName("ret_addr"))
			if err != nil {
				return fmt.Errorf("return stmt: %w", err)
			}

			if len(st.Results) > 1 {
				// tuple return
				// TODO: Add checks for return type
				for _, result := range st.Results {
					retVarName, err := handleExpr(uasm, out, result)
					if err != nil {
						return fmt.Errorf("handle expr: %w", err)
					}
					uasm.PushVar(retVarName)
				}
			} else if len(st.Results) == 1 {
				// standard return
				// TODO: Add checks for return type
				retVarName, err := handleExpr(uasm, out, st.Results[0])
				if err != nil {
					return fmt.Errorf("handle expr: %w", err)
				}
				uasm.PushVar(retVarName)
			} else {
				// void
				// TODO: Add checks for return type
			}
			uasm.JumpRetAddr()
		case *ast.IfStmt:
			// fmt.Println("handle ast.IfStmt")
			elseLabel := assembly.LabelName(uasm.GetNextId("else_label"))
			ifEndLabel := assembly.LabelName(uasm.GetNextId("if_end_label"))

			condVarName, err := handleExpr(uasm, out, st.Cond)
			if err != nil {
				return fmt.Errorf("error handling if cond: %v", err)
			}

			uasm.PushVar(condVarName)
			// if (!test) goto else
			uasm.JumpIfFalseLabel(elseLabel)
			// {}
			err = handleBlockStmt(uasm, out, st.Body)
			if err != nil {
				return fmt.Errorf("error handling if body: %v", err)
			}
			// goto if_end
			uasm.JumpLabel(ifEndLabel)
			// else:
			uasm.AddLabelCurrentAddr(elseLabel)
			err = handleBlockStmt(uasm, out, st.Else.(*ast.BlockStmt))
			if err != nil {
				return fmt.Errorf("error handling else body: %v", err)
			}
			// if_end:
			uasm.AddLabelCurrentAddr(ifEndLabel)

		case *ast.ForStmt:
		default:
			return fmt.Errorf("unsupported statement: %v", s)
		}
	}
	return nil
}

func handleCallExpr(uasm *assembly.UdonAssembly, out io.Writer, c *ast.CallExpr) (assembly.VarName, error) {
	return "", fmt.Errorf("%s: %w", "CallExpr", ErrNotImplemented)
}

func handleBinaryExpr(uasm *assembly.UdonAssembly, out io.Writer, be *ast.BinaryExpr) (assembly.VarName, error) {
	return "", fmt.Errorf("%s: %w", "binaryExpr", ErrNotImplemented)
}

func handleUnaryExpr(uasm *assembly.UdonAssembly, out io.Writer, ue *ast.UnaryExpr) (assembly.VarName, error) {
	return "", fmt.Errorf("%s: %w", "UnaryExpr", ErrNotImplemented)
}
func handleFuncType(uasm *assembly.UdonAssembly, out io.Writer, lit *ast.FuncType) (assembly.VarName, error) {
	return "", fmt.Errorf("%s: %w", "FuncType", ErrNotImplemented)
}
func handleFuncLit(uasm *assembly.UdonAssembly, out io.Writer, lit *ast.FuncLit) (assembly.VarName, error) {
	return "", fmt.Errorf("%s: %w", "FuncLit", ErrNotImplemented)
}
func handleIdent(uasm *assembly.UdonAssembly, out io.Writer, ident *ast.Ident) (assembly.VarName, error) {
	constNextID := uasm.GetNextId("ident")
	return constNextID, nil
}

func handleBasicLit(uasm *assembly.UdonAssembly, out io.Writer, lit *ast.BasicLit) (assembly.VarName, error) {
	constNextID := uasm.GetNextId("const")
	if lit.Kind == token.INT {
		uasm.VarTable.AddVar(constNextID, assembly.UdonTypeInt32, lit.Value)
	}
	return constNextID, nil
}

func handleExpr(uasm *assembly.UdonAssembly, out io.Writer, e ast.Expr) (assembly.VarName, error) {
	// fmt.Println("run: handleExpr")
	switch expr := e.(type) {
	case *ast.Ident:
		return handleIdent(uasm, out, expr)
	case *ast.FuncType:
		return handleFuncType(uasm, out, expr)
	case *ast.FuncLit:
		return handleFuncLit(uasm, out, expr)
	case *ast.CallExpr:
		return handleExpr(uasm, out, expr.Fun)
	case *ast.BinaryExpr:
		handleExpr(uasm, out, expr.X)
		return handleExpr(uasm, out, expr.Y)
	case *ast.UnaryExpr:
		return "", fmt.Errorf("UnaryExpr: %w", ErrNotImplemented)
	case *ast.BasicLit:
		return handleBasicLit(uasm, out, expr)
	default:
		return "", fmt.Errorf("%s: %w", expr, ErrNotImplemented)
	}
}
