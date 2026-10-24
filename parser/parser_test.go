package parser

import (
	"fmt"
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

	program := p.ParseProgram()
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

func TestReturnStatement(t *testing.T) {
	input := `
	return 4;
	return 5;
	return 12341;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	for _, stmt := range program.Statements {
		ret, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt is not ReturnStatement, got %T (%+v)", stmt, stmt)
			continue
		}

		if ret.TokenLiteral() != "return" {
			t.Errorf("not a return, got %q", ret.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("len is not 1, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not an ast.ExpressionStatement, got %T (%+v)", program.Statements[0], program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("not an ast.Identifier, got %T (%+v)", stmt.Expression, stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("not a foobar, got %q", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("not a foobar, got %q", ident.TokenLiteral())
	}
}

func TestIntegerLiteral(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("len is not 1, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not an ast.ExpressionStatement, got %T (%+v)", program.Statements[0], program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("not an ast.Identifier, got %T (%+v)", stmt.Expression, stmt.Expression)
	}

	if ident.Value != 5 {
		t.Errorf("not a 5, got %d", ident.Value)
	}
	if ident.TokenLiteral() != "5" {
		t.Errorf("not a 5, got %q", ident.TokenLiteral())
	}
}

func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		input string
		oper  string
		value int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("len of statements is not 1, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("got an ExpresstionStatement, got %T (%+v)", program.Statements[0], program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("got an PrefixExpresstion, got %T (%+v)", program.Statements[0], program.Statements[0])
		}

		if exp.Operator != tt.oper {
			t.Fatalf("exp.Operator is not %s", tt.oper)
		}
		if !testIntegerLiteral(t, exp.Right, tt.value) {
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, exp int64) bool {
	intg, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("not an *ast.IntegerLiteral, got %T (%+v)", il, il)
		return false
	}

	if intg.Value != exp {
		t.Errorf("intg.Value is not %d, got %d", exp, intg.Value)
		return false
	}

	if intg.TokenLiteral() != fmt.Sprintf("%d", exp) {
		t.Errorf("intg.TokenLiteral() is not %d, got %s", exp, intg.TokenLiteral())
		return false
	}

	return true
}

func TestInfixExpressions(t *testing.T) {
	tests := []struct {
		input string
		left  int64
		oper  string
		right int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 / 5", 5, "/", 5},
		{"5 * 5", 5, "*", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("len of statements is not 1, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("got an ExpresstionStatement, got %T (%+v)", program.Statements[0], program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("got an InfixExpression, got %T (%+v)", program.Statements[0], program.Statements[0])
		}

		if !testIntegerLiteral(t, exp.Left, tt.left) {
			return
		}

		if exp.Operator != tt.oper {
			t.Fatalf("exp.Operator is not %s, got %s", tt.oper, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.right) {
			return
		}
	}
}
