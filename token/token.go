package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	INT   = "INT"

	ASSIGN          = "="
	PLUS            = "+"
	PLUS_ASSIGN     = "+="
	MINUS           = "-"
	MINUS_ASSIGN    = "-="
	ASTERISK        = "*"
	ASTERISK_ASSIGN = "*="
	SLASH           = "/"
	SLASH_ASSIGN    = "/="
	MODULO          = "%"
	MODULO_ASSIGN   = "%="
	BANG            = "!"

	LT     = "<"
	LTE    = "<="
	GT     = ">"
	GTE    = ">="
	EQ     = "=="
	NOT_EQ = "!="

	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN = "("
	RPAREN = ")"

	LBRACE = "{"
	RBRACE = "}"

	LBRACKET = "["
	RBRACKET = "]"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	RETURN   = "RETURN"
	IF       = "IF"
	ELSE     = "ELSE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	STRING   = "STRING"
	NULL     = "NULL"
	WHILE    = "WHILE"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
)

var keywords = map[string]TokenType{
	"fn":       FUNCTION,
	"let":      LET,
	"return":   RETURN,
	"if":       IF,
	"else":     ELSE,
	"true":     TRUE,
	"false":    FALSE,
	"null":     NULL,
	"while":    WHILE,
	"break":    BREAK,
	"continue": CONTINUE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
