package evaluator

import (
	"testing"

	"github.com/egshH/Lexeme/lexer"
	"github.com/egshH/Lexeme/object"
	"github.com/egshH/Lexeme/parser"
)

func evaluate(t *testing.T, input string, env *object.Environment) object.Object {
	t.Helper()
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	if errors := p.Errors(); len(errors) != 0 {
		t.Fatalf("parser errors: %v", errors)
	}
	return Eval(program, env)
}

func TestEvaluation(t *testing.T) {
	for _, tt := range []struct {
		input string
		want  any
	}{
		{"5", int64(5)},
		{"-10", int64(-10)},
		{"5 + 2 * 10 - 8 / 2", int64(21)},
		{"(5 + 2) * 3", int64(21)},
		{"-5 * -2 + 10", int64(20)},
		{"7 / 2", int64(3)},
		{"-7 / 2", int64(-3)},
		{"1; 2; 3", int64(3)},
		{"true", true},
		{"false", false},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!5", true},
		{"!0", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 <= 1", true},
		{"1 >= 2", false},
		{"5 == 5", true},
		{"5 != 5", false},
		{"5 == 6", false},
		{"true == (1 < 2)", true},
		{"false != true", true},
		{"1 == true", false},
		{"1 != true", true},
		{"if (true) { 10 } else { 20 }", int64(10)},
		{"if (false) { 10 } else { 20 }", int64(20)},
		{"if (0) { 10 }", int64(10)},
		{"if (false) { 10 }", NULL},
		{"if (true) {}", NULL},
		{"if (true) { let a = 1; }", NULL},
		{"!(if (false) { 10 })", true},
		{"if (if (false) { 1 }) { 10 } else { 20 }", int64(20)},
		{"(if (false) { 1 }) == (if (false) { 2 })", true},
		{"true && false", false},
		{"false || true", true},
		{"false && (1 / 0)", false},
		{"true || unknown", true},
		{"0 && 2", true},
		{"(if (false) { 1 }) || false", false},
		{"if (true) { 10 } else { 1 / 0 }", int64(10)},
		{"if (false) { unknown } else { 20 }", int64(20)},
		{"return 10; 9", int64(10)},
		{"9; return 2 * 5; 9", int64(10)},
		{"if (true) { if (true) { return 10; } return 1; } 99", int64(10)},
		{"return; 99", NULL},
		{"if (true) { return; } 99", NULL},
		{"let a = 5; let b = a * 2; a + b", int64(15)},
		{"let a = if (1 < 2) { 10 } else { 20 }; a", int64(10)},
		{"let a = if (false) { 10 }; !a", true},
		{"let a = if (true) { return 10; }; 99", int64(10)},
		{"1 + if (true) { return 10; }; 99", int64(10)},
		{"return if (true) { return 10; }; 99", int64(10)},
	} {
		t.Run(tt.input, func(t *testing.T) {
			got := evaluate(t, tt.input, object.NewEnvironment())
			switch want := tt.want.(type) {
			case int64:
				value, ok := got.(*object.Integer)
				if !ok || value.Value != want {
					t.Fatalf("got %T (%+v), want integer %d", got, got, want)
				}
			case bool:
				value, ok := got.(*object.Boolean)
				if !ok || value.Value != want {
					t.Fatalf("got %T (%+v), want boolean %t", got, got, want)
				}
			default:
				if got != NULL {
					t.Fatalf("got %T (%+v), want NULL", got, got)
				}
			}
		})
	}
}

func TestErrors(t *testing.T) {
	for _, tt := range []struct{ input, message string }{
		{"5 + true; 5", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false", "unknown operator: BOOLEAN + BOOLEAN"},
		{"1 / 0; 99", "division by zero"},
		{"missing", "identifier not found: missing"},
		{"if (true) { if (true) { return true + false; } return 1; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (5 + true) { 1 }", "type mismatch: INTEGER + BOOLEAN"},
		{"-(5 + true)", "type mismatch: INTEGER + BOOLEAN"},
		{"(5 + true) * 2", "type mismatch: INTEGER + BOOLEAN"},
		{"1 + missing", "identifier not found: missing"},
		{"true && missing", "identifier not found: missing"},
		{"false || missing", "identifier not found: missing"},
		{"missing || true", "identifier not found: missing"},
		{"let a = 5 + true; 99", "type mismatch: INTEGER + BOOLEAN"},
		{"return 1 / 0; 99", "division by zero"},
		{"fn() {}", "evaluation not implemented for *ast.FunctionLiteral"},
		{"f()", "evaluation not implemented for *ast.CallExpression"},
		{"[]", "evaluation not implemented for *ast.ArrayLiteral"},
		{"{}", "evaluation not implemented for *ast.HashLiteral"},
		{`"hello"`, "evaluation not implemented for *ast.StringLiteral"},
		{"for {}", "evaluation not implemented for *ast.ForStatement"},
		{"try {}", "evaluation not implemented for *ast.TryStatement"},
		{"const a = 1", "evaluation not implemented for *ast.ConstStatement"},
		{"a = 1", "evaluation not implemented for *ast.AssignExpression"},
		{"a++", "evaluation not implemented for *ast.PostfixExpression"},
		{"break", "evaluation not implemented for *ast.BreakStatement"},
		{"continue", "evaluation not implemented for *ast.ContinueStatement"},
	} {
		t.Run(tt.input, func(t *testing.T) {
			got := evaluate(t, tt.input, object.NewEnvironment())
			err, ok := got.(*object.Error)
			if !ok || err.Message != tt.message {
				t.Fatalf("got %T (%+v), want error %q", got, got, tt.message)
			}
		})
	}
}

func TestEnvironmentPersistsAndErrorsStopExecution(t *testing.T) {
	env := object.NewEnvironment()
	if got := evaluate(t, "let a = 10", env); got != nil {
		t.Fatalf("declaration produced output: %v", got)
	}
	if got := evaluate(t, "a + 5", env).(*object.Integer).Value; got != 15 {
		t.Fatalf("got %d, want 15", got)
	}
	evaluate(t, "let a = 1 / 0; let b = 20", env)
	if a, _ := env.Get("a"); a.(*object.Integer).Value != 10 {
		t.Fatal("failed declaration changed previous value")
	}
	if _, ok := env.Get("b"); ok {
		t.Fatal("continued evaluating after error")
	}
	if got := evaluate(t, "", env); got != nil {
		t.Fatalf("empty program produced output: %v", got)
	}
	if got := evaluate(t, "a", object.NewEnvironment()); got.Type() != object.ERROR_OBJ {
		t.Fatal("binding leaked into an independent environment")
	}
}
