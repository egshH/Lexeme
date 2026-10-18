package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/egshH/Lexeme/lexer"
	"github.com/egshH/Lexeme/token"
)

const prompt = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	for {
		fmt.Print(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
