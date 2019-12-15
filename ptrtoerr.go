package ptrtoerr

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "ptrtoerr",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

const Doc = "ptrtoerr detects assignment pointer value to error interface."

var errType = types.Universe.Lookup("error").Type()

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.AssignStmt)(nil),
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.AssignStmt:
			checkAssign(pass, n)
		case *ast.FuncDecl:
			checkFuncReturn(pass, n.Type, n.Body)
		case *ast.FuncLit:
			checkFuncReturn(pass, n.Type, n.Body)
		}
	})

	return nil, nil
}

// checkAssign reports the assignment of a pointer to an error.
func checkAssign(pass *analysis.Pass, n *ast.AssignStmt) {
	for i := range n.Lhs {
		lt := pass.TypesInfo.TypeOf(n.Lhs[i])
		rt := pass.TypesInfo.TypeOf(n.Rhs[i])
		_, rtIsPtr := rt.(*types.Pointer)
		if lt == errType && rtIsPtr {
			pass.Reportf(n.Pos(), "Assign pointer to error")
		}
	}
}

// checkFuncReturn reports returning a pointer as error.
func checkFuncReturn(pass *analysis.Pass, t *ast.FuncType, b *ast.BlockStmt) {
	if t.Results == nil {
		return
	}
	var idxs []int
	for i, r := range t.Results.List {
		if pass.TypesInfo.TypeOf(r.Type) == errType {
			idxs = append(idxs, i)
		}
	}
	if len(idxs) == 0 {
		return
	}

	ast.Inspect(b, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.ReturnStmt:
			for _, i := range idxs {
				_, isPtr := pass.TypesInfo.TypeOf(n.Results[i]).(*types.Pointer)
				if isPtr {
					pass.Reportf(n.Pos(), "Return pointer as error")
				}
			}
		}
		return true
	})
}
