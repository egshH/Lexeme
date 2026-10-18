package main

import (
	"os"

	"github.com/egshH/Lexeme/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
