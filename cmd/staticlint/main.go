// Package main provides a custom static checker.
//
// This package registers and runs a set of analyzers for static code analysis.
// It includes both custom analyzers as well as analyzers from external packages.
//
// Usage:
//
//	To use this package, import it in your Go program and call the main function.
package main

import (
	"github.com/andreevym/metric-collector/pkg/noexitchecker"
	"github.com/gostaticanalysis/funcstat"
	"github.com/gostaticanalysis/zapvet/passes/fieldtype"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/staticcheck"
)

// main function registers and runs a set of analyzers.
func main() {
	var analyzers []*analysis.Analyzer

	// Add all analyzers from staticcheck package.
	for _, v := range staticcheck.Analyzers {
		analyzers = append(analyzers, v.Analyzer)
	}
	// Add analyzer from other package staticcheck.
	for _, qf := range quickfix.Analyzers {
		if qf.Analyzer.Name == "QF1006" {
			analyzers = append(analyzers, qf.Analyzer)
		}
	}

	// Add custom analyzer for checking os.Exit in main function.
	analyzers = append(analyzers, noexitchecker.Analyzer)

	// Add external analyzers.
	analyzers = append(analyzers, fieldtype.Analyzer)
	analyzers = append(analyzers, funcstat.Analyzer)

	// Add passes analyzers.
	analyzers = append(analyzers, assign.Analyzer)
	analyzers = append(analyzers, findcall.Analyzer)
	analyzers = append(analyzers, inspect.Analyzer)
	analyzers = append(analyzers, printf.Analyzer)
	analyzers = append(analyzers, shadow.Analyzer)
	analyzers = append(analyzers, shift.Analyzer)
	analyzers = append(analyzers, structtag.Analyzer)

	// Run all analyzers.
	multichecker.Main(analyzers...)
}
