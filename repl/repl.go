package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/egshH/Lexeme/lexer"
	"github.com/egshH/Lexeme/parser"
)

const prompt = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	for {
		fmt.Fprint(out, prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		fmt.Fprintln(out, program.String())
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		fmt.Fprintln(out, "\t"+msg)
	}
}
