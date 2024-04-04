package noexitchecker

import "golang.org/x/tools/go/analysis"

var Analyzer = &analysis.Analyzer{
	Name: "noexitchecker",
	Doc:  "check is exist in package 'main' in func 'main'",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		foundPos := FindOsExitInMain(file)
		if foundPos != nil {
			pos := *foundPos
			pass.Reportf(pos, "detected direct os.Exit() call in main function")
		}
	}
	return nil, nil
}
