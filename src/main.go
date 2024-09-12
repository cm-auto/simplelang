package main

import (
	"os"
	"simplelang/src/ast"
	"simplelang/src/builder"
	"simplelang/src/token"
)

func main() {
	inputFileBytes, err := os.ReadFile("in/main.sl")
	if err != nil {
		panic(err)
	}
	inputFile := string(inputFileBytes)
	tokens := token.Tokenize(inputFile)

	ast_ := ast.NewAst(tokens)

	goSourceCode := builder.BuildProgram(ast_)
	err = os.WriteFile("out/main.go", []byte(goSourceCode), 0644)
	if err != nil {
		panic(err)
	}
}
