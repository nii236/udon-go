package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"strconv"
	"strings"
	"udon-go/assembly"
)

func handleDecls(uasm *assembly.UdonAssembly, out io.Writer, d []ast.Decl) error {
	fmt.Println("run: handleDecls")
	for _, decl := range d {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			return handleGenDecl(uasm, out, decl)
		case *ast.FuncDecl:
			return handleFuncDecl(uasm, out, decl)
		default:
			return fmt.Errorf("unsupported decl: %#v", d)
		}
	}
	return ErrNotImplemented
}

func handleGenDecl(uasm *assembly.UdonAssembly, out io.Writer, decl *ast.GenDecl) error {
	fmt.Println("run: handleGenDecl")
	for _, s := range decl.Specs {
		vs, ok := s.(*ast.ValueSpec)
		if !ok {
			return fmt.Errorf("unsupported spec: %#v", s)
		}
		if len(vs.Names) > 1 {
			return fmt.Errorf("unsupported # of value names: %v", vs.Names)
		}
		decl := []string{}
		if vs.Names[0].Obj.Kind == ast.Con {
			decl = append(decl, "const")
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
		case *ast.FuncLit:
			// funcName := assembly.VarName(vs.Names[0].Name)
			// argTypes := []assembly.UdonTypeName{}
			// for _, field := range l.Type.Params.List {
			// 	switch f := field.Tag.Kind {
			// 	case token.INT:
			// 		argTypes = append(argTypes, f.String())
			// 	case token.FLOAT:
			// 		argTypes = append(argTypes, f.String())
			// 	case token.IMAG:
			// 		argTypes = append(argTypes, f.String())
			// 	case token.CHAR:
			// 		argTypes = append(argTypes, f.String())
			// 	case token.STRING:
			// 		argTypes = append(argTypes, f.String())
			// 	}

			// }
			// uasm.DefFuncTable.AddFunc(funcName)
		}

	}
	return nil
}
func handleFuncDecl(uasm *assembly.UdonAssembly, out io.Writer, decl *ast.FuncDecl) error {
	fmt.Println("run: handleFuncDecl")
	// for _, param := range decl.Type.Params.List {

	// }
	// for _, result := range decl.Type.Results.List {
	// 	result.Type
	// }
	// uasm.DefFuncTable.AddFunc(
	// 	assembly.FuncName(decl.Name.Name),
	// 	nil,
	// )
	return ErrNotImplemented
}
func handleBlockStmt(uasm *assembly.UdonAssembly, out io.Writer, bs *ast.BlockStmt) error {
	fmt.Println("run: handleBlockStmt")
	for _, s := range bs.List {
		switch st := s.(type) {
		case *ast.ExprStmt:
			if err := handleExpr(uasm, out, st.X); err != nil {
				return fmt.Errorf("error handling expr stmt %v: %v", st.X, err)
			}
			fmt.Fprint(out, ";\n")
		case *ast.AssignStmt:
			if len(st.Lhs) > 1 {
				return fmt.Errorf("unsupported # of lhs exprs: %v", st.Lhs)

			}
			if err := handleExpr(uasm, out, st.Lhs[0]); err != nil {
				return fmt.Errorf("error handling left expr %v: %v", st.Lhs[0], err)
			}
			fmt.Fprintf(out, "=")
			if len(st.Rhs) > 1 {
				return fmt.Errorf("unsupported # of rhs exprs: %v", st.Rhs)

			}
			if err := handleExpr(uasm, out, st.Rhs[0]); err != nil {
				return fmt.Errorf("error handling right expr %v: %v", st.Rhs[0], err)
			}
			fmt.Fprint(out, ";\n")

		case *ast.IfStmt:
			fmt.Fprintf(out, "if (")
			if err := handleExpr(uasm, out, st.Cond); err != nil {
				return fmt.Errorf("error handling if block conditionx: %v", err)
			}
			fmt.Fprint(out, ") {\n")
			if err := handleBlockStmt(uasm, out, st.Body); err != nil {
				return fmt.Errorf("error handling if block statements: %v", err)
			}
			fmt.Fprintf(out, "}")
			if st.Else != nil {
				bs, ok := st.Else.(*ast.BlockStmt)
				if !ok {
					return fmt.Errorf("unsupported statement: %v", st.Else)
				}
				fmt.Fprintf(out, " else {\n")
				if err := handleBlockStmt(uasm, out, bs); err != nil {
					return fmt.Errorf("error handling else block statements: %v", err)
				}
				fmt.Fprintf(out, "}")
			}
			fmt.Fprintln(out)
		default:
			return fmt.Errorf("unsupported statement: %v", s)

		}
	}
	return ErrNotImplemented
}

func handleCallExpr(uasm *assembly.UdonAssembly, out io.Writer, c *ast.CallExpr) error {
	fmt.Println("run: handleCallExpr")
	funcName, ok := c.Fun.(*ast.Ident)
	if !ok {
		return fmt.Errorf("unsupported func expr:uasm *assembly.UdonAssembly,  %#v", c.Fun)
	}
	args := []string{}
	for _, a := range c.Args {
		var buf bytes.Buffer
		if err := handleExpr(uasm, &buf, a); err != nil {
			return fmt.Errorf("error handling func arg uasm *assembly.UdonAssembly, expr %#v: %v", a, err)
		}
		args = append(args, buf.String())
	}
	fmt.Fprintf(out, "%s(%s)", funcName, strings.Join(args, ", "))
	return ErrNotImplemented
}

func handleBinaryExpr(uasm *assembly.UdonAssembly, out io.Writer, be *ast.BinaryExpr) error {
	fmt.Println("run: handleBinaryExpr")
	if err := handleExpr(uasm, out, be.X); err != nil {
		return fmt.Errorf("error handling left part %v of binary expr: %v", be.X, err)
	}
	fmt.Fprint(out, be.Op)
	if err := handleExpr(uasm, out, be.Y); err != nil {
		return fmt.Errorf("error handling right part %v of binary expr: %v", be.Y, err)
	}
	return ErrNotImplemented
}

func handleUnaryExpr(uasm *assembly.UdonAssembly, out io.Writer, ue *ast.UnaryExpr) error {
	fmt.Println("run: handleUnaryExpr")
	fmt.Fprint(out, ue.Op)
	if err := handleExpr(uasm, out, ue.X); err != nil {
		return err
	}
	return ErrNotImplemented
}

func handleIdent(uasm *assembly.UdonAssembly, out io.Writer, ident *ast.Ident) error {
	fmt.Println("run: handleIdent")
	fmt.Fprintf(out, ident.Name)
	return ErrNotImplemented
}

func handleBasicLit(uasm *assembly.UdonAssembly, out io.Writer, lit *ast.BasicLit) error {
	fmt.Println("run: handleBasicLit")
	fmt.Fprintf(out, lit.Value)
	return ErrNotImplemented
}

func handleExpr(uasm *assembly.UdonAssembly, out io.Writer, e ast.Expr) error {
	fmt.Println("run: handleExpr")
	switch expr := e.(type) {
	case *ast.CallExpr:
		return handleCallExpr(uasm, out, expr)
	case *ast.BinaryExpr:
		return handleBinaryExpr(uasm, out, expr)
	case *ast.UnaryExpr:
		return handleUnaryExpr(uasm, out, expr)
	case *ast.Ident:
		return handleIdent(uasm, out, expr)
	case *ast.BasicLit:
		return handleBasicLit(uasm, out, expr)
	default:
		return fmt.Errorf("unsupported expr: %#v", e)
	}
}
