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
		if okFile && pkg.Name.Name == "main" {
			ast.Inspect(pkg, func(n2 ast.Node) bool {
				a, okFuncDecl := n2.(*ast.FuncDecl)
				if okFuncDecl && a.Name.Name == "main" {
					ast.Inspect(a.Body, func(n3 ast.Node) bool {
						b, okCallExpr := n3.(*ast.CallExpr)
						if okCallExpr {
							fun, ok := b.Fun.(*ast.SelectorExpr)
							if !ok {
								return false
							}
							xFunIdent, ok := fun.X.(*ast.Ident)
							if !ok {
								return false
							}
							selFunIdent := fun.Sel
							if xFunIdent.Name == "os" && selFunIdent.Name == "Exit" {
								t := xFunIdent.Pos()
								pos = &t
								return false
							}
						}
						return true
					})
					return false
				}
				return true
			})
			return false
		}
		return true
	})

	return pos
}
