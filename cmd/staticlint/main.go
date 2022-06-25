package main

import (
	"go/ast"

	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"honnef.co/go/tools/staticcheck"
)

// OsExitChecker os.Exit checker in main func
var OsExitChecker = &analysis.Analyzer{
	Name: "OsExitChecker",
	Doc:  "os.Exit checker in main func",
	Run:  runCheck,
}

func runCheck(source *analysis.Pass) (interface{}, error) {
	for _, file := range source.Files {
		if file.Name.Name == "main" {
			ast.Inspect(file, func(node ast.Node) bool {
				if x, ok := node.(*ast.CallExpr); ok {
					if f, ok := x.Fun.(*ast.SelectorExpr); ok {
						if p, ok := f.X.(*ast.Ident); ok {
							if p.Name == "os" && f.Sel.Name == "Exit" {
								source.Reportf(f.Pos(), "os.Exit() detected!")
							}
						}
					}
				}
				return true
			})
		}
	}
	return nil, nil
}

func main() {
	checkers := []*analysis.Analyzer{
		OsExitChecker,
		httpresponse.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
		errcheck.Analyzer,
		ineffassign.Analyzer,
		assign.Analyzer,
		bools.Analyzer,
		copylock.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
	}

	checkMap := map[string]bool{
		"QF1001": true,
		"SA":     true,
		"S1028":  true,
		"ST1006": true,
	}
	for _, v := range staticcheck.Analyzers {
		if checkMap[v.Name] {
			checkers = append(checkers, v)
		}
	}
	multichecker.Main(checkers...)
}
