package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT        = "IDENT"
	INT          = "INT"
	STRING       = "STRING"
	COLON        = ":"
	LBRACKET     = "["
	RBRACKET     = "]"
	LT_EQ        = "<="
	GT_EQ        = ">="
	PLUS_ASSIGN  = "+="
	MINUS_ASSIGN = "-="
	INCREMENT    = "++"
	DECREMENT    = "--"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	DOTDOT   = ".."

	LT    = "<" // less then
	GT    = ">" // greater then
	EQ    = "=="
	NotEQ = "!="

	AND = "&&"
	OR  = "||"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	CONST    = "CONST"
	IF       = "IF"
	ELSE     = "ELSE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	RETURN   = "RETURN"
	FOR      = "FOR"
	RANGE    = "RANGE"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
	TRY      = "TRY"
	CATCH    = "CATCH"
)

var keywords = map[string]TokenType{
	"fn":       FUNCTION,
	"let":      LET,
	"const":    CONST,
	"if":       IF,
	"else":     ELSE,
	"true":     TRUE,
	"false":    FALSE,
	"return":   RETURN,
	"for":      FOR,
	"range":    RANGE,
	"break":    BREAK,
	"continue": CONTINUE,
	"try":      TRY,
	"catch":    CATCH,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
