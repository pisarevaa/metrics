package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Напишите и добавьте в multichecker собственный анализатор, запрещающий использовать
// прямой вызов os.Exit в функции main пакета main. При необходимости перепишите код
// своего проекта так, чтобы он удовлетворял данному анализатору.
var ErrCheckAnalyzer = &analysis.Analyzer{
	Name: "osexitmain",
	Doc:  "check for os.Exit in main.go",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// функцией ast.Inspect проходим по всем узлам AST
		ast.Inspect(file, func(node ast.Node) bool {
			if file.Name.Name == "main" {
				// Check if it's a call expression
				callExpr, ok := node.(*ast.CallExpr)
				if !ok {
					return true
				}
				// Check if it's an identifier and if the package is "os"
				if fun, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					if ident, ok := fun.X.(*ast.Ident); ok && ident.Name == "os" && fun.Sel.Name == "Exit" {
						pass.Reportf(ident.NamePos, "used os.Exit in main.go")
						return false // Stop traversing
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
