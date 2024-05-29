package ttmp

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "ttmp analyzes the code that is using os.TempDir instead of t.TempDir."

var (
	A     = "all"
	aflag bool
)

func init() {
	Analyzer.Flags.BoolVar(&aflag, A, false, "the all option will run against all method in test file")
}

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "ttmp",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
		(*ast.ReturnStmt)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			checkFuncDecl(pass, n, pass.Fset.File(n.Pos()).Name())
		case *ast.FuncLit:
			checkFuncLit(pass, n, pass.Fset.File(n.Pos()).Name())
		case *ast.ReturnStmt:
			checkReturnStmt(pass, n)
		}
	})

	return nil, nil
}

func checkFuncDecl(pass *analysis.Pass, f *ast.FuncDecl, fileName string) {
	argName, ok := targetRunner(f.Type.Params.List, fileName)
	if !ok {
		return
	}
	checkStmts(pass, f.Body.List, f.Name.Name, argName)
}

func checkFuncLit(pass *analysis.Pass, f *ast.FuncLit, fileName string) {
	argName, ok := targetRunner(f.Type.Params.List, fileName)
	if !ok {
		return
	}
	checkStmts(pass, f.Body.List, "anonymous function", argName)
}

func checkReturnStmt(pass *analysis.Pass, stmt *ast.ReturnStmt) {
	for _, result := range stmt.Results {
		switch result := result.(type) {
		case *ast.CallExpr:
			fun, ok := result.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			x, ok := fun.X.(*ast.Ident)
			if !ok {
				continue
			}
			targetName := x.Name + "." + fun.Sel.Name
			if targetName == "os.MkdirTemp" {
				pass.Reportf(stmt.Pos(), "os.MkdirTemp() can be replaced by `testing.TempDir()`")
			}
		}
	}
}

func checkStmts(pass *analysis.Pass, stmts []ast.Stmt, funcName, argName string) {
	for _, stmt := range stmts {
		switch stmt := stmt.(type) {
		case *ast.ExprStmt:
			if !checkExprStmt(pass, stmt, funcName, argName) {
				continue
			}
		case *ast.IfStmt:
			if !checkIfStmt(pass, stmt, funcName, argName) {
				continue
			}
		case *ast.AssignStmt:
			if !checkAssignStmt(pass, stmt, funcName, argName) {
				continue
			}
		}
	}
}

func checkExprStmt(pass *analysis.Pass, stmt *ast.ExprStmt, funcName, argName string) bool {
	callExpr, ok := stmt.X.(*ast.CallExpr)
	if !ok {
		return false
	}
	fun, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	x, ok := fun.X.(*ast.Ident)
	if !ok {
		return false
	}
	targetName := x.Name + "." + fun.Sel.Name
	if targetName == "os.MkdirTemp" {
		if argName == "" {
			argName = "testing"
		}
		pass.Reportf(stmt.Pos(), "os.MkdirTemp() can be replaced by `%s.TempDir()` in %s", argName, funcName)
	}
	return true
}

func checkIfStmt(pass *analysis.Pass, stmt *ast.IfStmt, funcName, argName string) bool {
	assignStmt, ok := stmt.Init.(*ast.AssignStmt)
	if !ok {
		return false
	}
	rhs, ok := assignStmt.Rhs[0].(*ast.CallExpr)
	if !ok {
		return false
	}
	fun, ok := rhs.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	x, ok := fun.X.(*ast.Ident)
	if !ok {
		return false
	}
	targetName := x.Name + "." + fun.Sel.Name
	if targetName == "os.MkdirTemp" {
		if argName == "" {
			argName = "testing"
		}
		pass.Reportf(stmt.Pos(), "os.MkdirTemp() can be replaced by `%s.TempDir()` in %s", argName, funcName)
	}
	return true
}

func checkAssignStmt(pass *analysis.Pass, stmt *ast.AssignStmt, funcName, argName string) bool {
	rhs, ok := stmt.Rhs[0].(*ast.CallExpr)
	if !ok {
		return false
	}
	fun, ok := rhs.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	x, ok := fun.X.(*ast.Ident)
	if !ok {
		return false
	}
	targetName := x.Name + "." + fun.Sel.Name
	if targetName == "os.MkdirTemp" {
		if argName == "" {
			argName = "testing"
		}
		pass.Reportf(stmt.Pos(), "os.MkdirTemp() can be replaced by `%s.TempDir()` in %s", argName, funcName)
	}
	return true
}

func targetRunner(params []*ast.Field, fileName string) (string, bool) {
	for _, p := range params {
		switch typ := p.Type.(type) {
		case *ast.StarExpr:
			if checkStarExprTarget(typ) {
				if len(p.Names) == 0 {
					return "", false
				}
				argName := p.Names[0].Name
				return argName, true
			}
		case *ast.SelectorExpr:
			if checkSelectorExprTarget(typ) {
				if len(p.Names) == 0 {
					return "", false
				}
				argName := p.Names[0].Name
				return argName, true
			}
		}
	}
	if aflag && strings.HasSuffix(fileName, "_test.go") {
		return "", true
	}
	return "", false
}

func checkStarExprTarget(typ *ast.StarExpr) bool {
	selector, ok := typ.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	x, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}
	targetName := x.Name + "." + selector.Sel.Name
	switch targetName {
	case "testing.T", "testing.B":
		return true
	default:
		return false
	}
}

func checkSelectorExprTarget(typ *ast.SelectorExpr) bool {
	x, ok := typ.X.(*ast.Ident)
	if !ok {
		return false
	}
	targetName := x.Name + "." + typ.Sel.Name
	return targetName == "testing.TB"
}
