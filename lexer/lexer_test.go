package lexer

import (
	"testing"

	"github.com/egshH/Lexeme/token"
)

func TestNextToken(t *testing.T) {
	input := `
	let five = 5;
	let ten = 10;

	let add = fn(x, y) {
		x + y;
	}

	let res = add(five, ten);
	!-/*5;
	5 < 10 > 5;

	if (5 < 10) {
		return true;
	} else {
		return false;
	}

	10 == 10
	5 != 10

	for break continue
	try catch
	const
	&& ||

	`

	tests := []struct {
		expectType    token.TokenType
		expectLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},

		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.LET, "let"},
		{token.IDENT, "res"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},

		{token.INT, "5"},
		{token.NotEQ, "!="},
		{token.INT, "10"},

		{token.FOR, "for"},
		{token.BREAK, "break"},
		{token.CONTINUE, "continue"},
		{token.TRY, "try"},
		{token.CATCH, "catch"},
		{token.CONST, "const"},
		{token.AND, "&&"},
		{token.OR, "||"},

		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectType {
			t.Fatalf("[%d] wrong token type: want %s got %s", i, tt.expectType, tok.Type)
		}

		if tok.Literal != tt.expectLiteral {
			t.Fatalf("[%d] wrong token type: want %s got %s", i, tt.expectLiteral, tok.Literal)
		}
	}
}

func TestExtendedTokens(t *testing.T) {
	l := New("// comment\nitem2 += 09; item2--; a++ >= b <= c -= 1; [\"a\\\"b\\\\c\\n\"]: {} // eof")
	expected := []token.Token{
		{Type: token.IDENT, Literal: "item2"}, {Type: token.PLUS_ASSIGN, Literal: "+="}, {Type: token.INT, Literal: "09"}, {Type: token.SEMICOLON, Literal: ";"},
		{Type: token.IDENT, Literal: "item2"}, {Type: token.DECREMENT, Literal: "--"}, {Type: token.SEMICOLON, Literal: ";"},
		{Type: token.IDENT, Literal: "a"}, {Type: token.INCREMENT, Literal: "++"}, {Type: token.GT_EQ, Literal: ">="}, {Type: token.IDENT, Literal: "b"},
		{Type: token.LT_EQ, Literal: "<="}, {Type: token.IDENT, Literal: "c"}, {Type: token.MINUS_ASSIGN, Literal: "-="}, {Type: token.INT, Literal: "1"}, {Type: token.SEMICOLON, Literal: ";"},
		{Type: token.LBRACKET, Literal: "["}, {Type: token.STRING, Literal: "a\"b\\c\n"}, {Type: token.RBRACKET, Literal: "]"}, {Type: token.COLON, Literal: ":"},
		{Type: token.LBRACE, Literal: "{"}, {Type: token.RBRACE, Literal: "}"}, {Type: token.EOF, Literal: ""},
	}
	for i, want := range expected {
		if got := l.NextToken(); got != want {
			t.Fatalf("token %d: got %+v, want %+v", i, got, want)
		}
	}
}

func TestInvalidStrings(t *testing.T) {
	for _, input := range []string{"\"unfinished", "\"unfinished\\", `"bad\q"`} {
		l := New(input)
		if got := l.NextToken(); got.Type != token.ILLEGAL {
			t.Fatalf("%q: got %+v", input, got)
		}
		if got := l.NextToken(); got.Type != token.EOF {
			t.Fatalf("expected EOF, got %+v", got)
		}
	}
}
