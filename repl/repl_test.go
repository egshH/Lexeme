package repl

import (
	"bytes"
	"strings"
	"testing"
)

func TestReadParsePrintLoop(t *testing.T) {
	var out bytes.Buffer
	Start(strings.NewReader("let x = 1 + 2 * 3\ntrue == false\nadd()\n"), &out)
	want := ">> let x = (1 + (2 * 3));\n>> (true == false)\n>> add()\n>> "
	if out.String() != want {
		t.Fatalf("got %q, want %q", out.String(), want)
	}
}

func TestParserErrorsAndContinue(t *testing.T) {
	var out bytes.Buffer
	Start(strings.NewReader("let x 5;\n5\n"), &out)
	if !strings.Contains(out.String(), "parser errors:") || !strings.Contains(out.String(), "expected next token") || !strings.HasSuffix(out.String(), ">> 5\n>> ") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}
