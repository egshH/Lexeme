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

func parseChapterInput(t *testing.T, input string) *ast.Program {
	t.Helper()
	p := New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)
	return program
}

func TestPrecedence(t *testing.T) {
	tests := []struct{ input, want string }{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"true == false", "(true == false)"},
		{"!true == false", "((!true) == false)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
		{"fn(x) { x; }(5)(6)", "fn(x)x(5)(6)"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := parseChapterInput(t, tt.input).String(); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBindingsAndReturns(t *testing.T) {
	for _, tt := range []struct {
		value string
		want  string
	}{
		{"5", "5"}, {"true", "true"}, {"foobar", "foobar"},
		{"1 + 2 * 3", "(1 + (2 * 3))"}, {"fn(x) { x }", "fn(x)x"},
	} {
		for _, ending := range []string{"", ";"} {
			program := parseChapterInput(t, "let my_value = "+tt.value+ending)
			stmt := program.Statements[0].(*ast.LetStatement)
			if stmt.Name.Value != "my_value" || stmt.Value == nil || stmt.Value.String() != tt.want {
				t.Fatalf("incorrect binding: %+v", stmt)
			}
			program = parseChapterInput(t, "return "+tt.value+ending)
			ret := program.Statements[0].(*ast.ReturnStatement)
			if ret.ReturnValue == nil || ret.ReturnValue.String() != tt.want {
				t.Fatalf("incorrect return: %+v", ret)
			}
		}
	}
}

func TestBooleans(t *testing.T) {
	for _, input := range []string{"true", "false"} {
		stmt := parseChapterInput(t, input).Statements[0].(*ast.ExpressionStatement)
		b, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("got %T", stmt.Expression)
		}
		if b.Value != (input == "true") || b.TokenLiteral() != input || stmt.TokenLiteral() != input {
			t.Fatalf("incorrect boolean: %+v", b)
		}
	}
}

func TestIfElse(t *testing.T) {
	for _, alternative := range []string{"", " else { y }"} {
		stmt := parseChapterInput(t, "if (x < y) { x }"+alternative).Statements[0].(*ast.ExpressionStatement)
		exp, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("got %T", stmt.Expression)
		}
		if exp.TokenLiteral() != "if" || exp.Condition.String() != "(x < y)" || len(exp.Consequence.Statements) != 1 || exp.Consequence.String() != "x" {
			t.Fatalf("incorrect conditional: %+v", exp)
		}
		if alternative == "" {
			if exp.Alternative != nil {
				t.Fatal("unexpected alternative")
			}
		} else if exp.Alternative == nil || exp.Alternative.String() != "y" {
			t.Fatal("missing alternative")
		}
	}
}

func TestFunctionAndCall(t *testing.T) {
	for _, tt := range []struct {
		input  string
		params []string
	}{
		{"fn() {}", nil}, {"fn(x) { x }", []string{"x"}}, {"fn(x, y, z) { x + y + z }", []string{"x", "y", "z"}},
	} {
		exp := parseChapterInput(t, tt.input).Statements[0].(*ast.ExpressionStatement).Expression.(*ast.FunctionLiteral)
		if exp.TokenLiteral() != "fn" || len(exp.Parameters) != len(tt.params) {
			t.Fatal("incorrect function parameters")
		}
		for i, name := range tt.params {
			if exp.Parameters[i].Value != name || exp.Parameters[i].TokenLiteral() != name {
				t.Fatal("incorrect parameter")
			}
		}
	}
	exp := parseChapterInput(t, "add(1, 2 * 3, 4 + 5)").Statements[0].(*ast.ExpressionStatement).Expression.(*ast.CallExpression)
	if exp.Function.(*ast.Identifier).Value != "add" || len(exp.Arguments) != 3 {
		t.Fatal("incorrect call")
	}
	for i, want := range []string{"1", "(2 * 3)", "(4 + 5)"} {
		if exp.Arguments[i].String() != want {
			t.Fatalf("argument %d", i)
		}
	}
	empty := parseChapterInput(t, "add()").Statements[0].(*ast.ExpressionStatement).Expression.(*ast.CallExpression)
	if len(empty.Arguments) != 0 {
		t.Fatal("unexpected arguments")
	}
	program := parseChapterInput(t, "let add = fn(x, y) { if (x < y) { return x } else { return y } }; add(1, 2)")
	if len(program.Statements) != 2 {
		t.Fatal("incorrect program length")
	}
}

func TestInvalidSyntax(t *testing.T) {
	for _, input := range []string{"let x =", "return", "(1 + 2", "if (true) { 1", "fn(x) { x", "fn(1) {}", "fn(x,) {}", "add(1,)", "add(1", "if true { 1 }", "if (true) { 1 } else 2", "1 && 2", "for (true) {}", "const x = 1"} {
		t.Run(input, func(t *testing.T) {
			p := New(lexer.New(input))
			p.ParseProgram()
			if len(p.Errors()) == 0 {
				t.Fatal("expected parser error")
			}
		})
	}
}
