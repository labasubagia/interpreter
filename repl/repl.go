package repl

import (
	"bufio"
	"fmt"
	"io"
	"os/user"

	"github.com/labasubagia/interpreter/evaluator"
	"github.com/labasubagia/interpreter/lexer"
	"github.com/labasubagia/interpreter/object"
	"github.com/labasubagia/interpreter/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the NEW Programming Language!\n", user.Username)
	fmt.Println("Feel free to type in commands")

	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
