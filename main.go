package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/labasubagia/interpreter/evaluator"
	"github.com/labasubagia/interpreter/lexer"
	"github.com/labasubagia/interpreter/object"
	"github.com/labasubagia/interpreter/parser"
	"github.com/labasubagia/interpreter/repl"
)

func eval(input string) {
	env := object.NewEnvironment()

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for i, e := range p.Errors() {
			fmt.Printf("error parse-%d: %s\n", i, e)
		}
	}
	obj := evaluator.Eval(program, env)
	if obj != nil && obj.Type() == object.ERROR_OBJ {
		fmt.Println(obj.Inspect())
	}
}

func main() {

	switch {
	case len(os.Args) >= 3:
		switch os.Args[1] {
		case "string":
			eval(os.Args[2])
		case "file":
			file := os.Args[2]

			// using file
			if ext := filepath.Ext(file); ext != ".newpl" {
				panic("file extension must be .newpl")
			}
			b, err := os.ReadFile(file)
			if err != nil {
				panic(err)
			}
			eval(string(b))
		}
	default:
		repl.Start(os.Stdin, os.Stdout)
	}

}
