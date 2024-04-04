package noexitchecker

import (
	"go/ast"
	"go/token"
)

// FindOsExitInMain returns an error if a call to os.Exit is found in the main function of the main package.
// The error contains the position where the os.Exit call was found.
func FindOsExitInMain(file ast.Node) *token.Pos {
	var pos *token.Pos
	ast.Inspect(file, func(n1 ast.Node) bool {
		pkg, okFile := n1.(*ast.File)
		if !okFile || pkg.Name.Name != "main" {
			return true
		}

		ast.Inspect(pkg, func(n2 ast.Node) bool {
			a, okFuncDecl := n2.(*ast.FuncDecl)
			if !okFuncDecl || a.Name.Name != "main" {
				return true
			}

			ast.Inspect(a.Body, func(bodyNode ast.Node) bool {
				callExpr, ok := bodyNode.(*ast.CallExpr)
				if !ok {
					return true
				}

				selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				ident, ok := selectorExpr.X.(*ast.Ident)
				if !ok || ident.Name != "os" || selectorExpr.Sel.Name != "Exit" {
					return true
				}

				pos = &ident.NamePos
				return false
			})
			return false
		})
		return false
	})

	return pos
}
