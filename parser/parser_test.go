package parser

import (
	"testing"

	"github.com/egshH/Lexeme/ast"
	"github.com/egshH/Lexeme/lexer"
)

func TestLetStatement(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foo = 123;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParserProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) == 0 {
		t.Fatalf("program.Statements are empty")
	}

	tests := []struct {
		expect string
	}{
		{"x"}, {"y"}, {"foo"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expect) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, exp string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral is not let, got %s", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("not a let Statement, got %T (%+v)", s, s)
		return false
	}

	if letStmt.Name.Value != exp {
		t.Errorf("letStmt.Name.Value is not '%s', got '%s'", exp, letStmt.Name.Value)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(p.errors))
	for _, err := range errors {
		t.Errorf("parser error: %s", err)
	}
	t.FailNow()
}
