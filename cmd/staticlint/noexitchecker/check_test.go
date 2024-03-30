package noexitchecker

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
)

func IsOsExitFound(src string) bool {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", src, 0)
	if err != nil {
		panic(err)
	}

	return FindOsExitInMain(file) != nil
}

func TestMainExit_True(t *testing.T) {
	src := `package main
import "fmt"

func main() {
    fmt.Println("Hello, world!")
	os.Exit(0)
}`

	ok := IsOsExitFound(src)
	require.True(t, ok)
}

func TestMainWithOtherFunInFileExit_True(t *testing.T) {
	src := `package main
import "fmt"

func main2() {
    fmt.Println("Hello, world!")
}

func main() {
    fmt.Println("Hello, world!")
	os.Exit(0)
}`

	ok := IsOsExitFound(src)
	require.True(t, ok)
}

func TestMainWithOtherFunExitInFile_True(t *testing.T) {
	src := `package main
import "fmt"

func d() {
    fmt.Println("Hello, world!")
    os.Exit(0)
}

func main() {
    fmt.Println("Hello, world!")
	d()
}`

	ok := IsOsExitFound(src)
	require.False(t, ok)
}

func TestMainWithOtherFunInFileNoExit_True(t *testing.T) {
	src := `package main
import "fmt"

func main2() {
    fmt.Println("Hello, world!")
	os.Exit(0)
}

func main() {
    fmt.Println("Hello, world!")
}`

	ok := IsOsExitFound(src)
	require.False(t, ok)
}

func TestMainIfExit_True(t *testing.T) {
	src := `package main
import "fmt"

func main() {
    fmt.Println("Hello, world!")
	if true {
	  os.Exit(0)
	}
}`

	ok := IsOsExitFound(src)
	require.True(t, ok)
}

func TestOtherPackageMainExit_False(t *testing.T) {
	src := `package main2
import "fmt"

func main() {
    fmt.Println("Hello, world!")
	os.Exit(0)
}`

	ok := IsOsExitFound(src)
	require.False(t, ok)
}

func TestOtherPackageMain2Exit_False(t *testing.T) {
	src := `package main2
import "fmt"

func main2() {
    fmt.Println("Hello, world!")
	os.Exit(0)
}`

	ok := IsOsExitFound(src)
	require.False(t, ok)
}

func TestMainPackageMainNoExit_False(t *testing.T) {
	src := `package main
import "fmt"

func main() {
    fmt.Println("Hello, world!")
}`

	ok := IsOsExitFound(src)
	require.False(t, ok)
}
