package lexer

import (
	"testing"

	"github.com/labasubagia/interpreter/token"
)

func TestNextToken(t *testing.T) {
	input := `
		# comment at the start 1
		# comment at the start 2

		let five = 5;
		let ten = 10;

		let add =  fn(x, y) {
			x + y;
		};

		# comment in the middle 1

		let result = add(five, ten);

		!-/*5;
		5 < 10 > 5;

		if ( 5 < 10 ) {
			return true
		} else {
			return false
		}

		# comment in the middle 2

		10 == 10;
		9 != 10;
		"foobar"
		"foo bar"
		"hello \"world\""
		"hello\n world"
		"hello\t\t\tworld"
		[1, 2];
		{"foo": "bar"};

		10 <= 12;
		5 >= 3;

		a += 1;
		b -= 2;
		c *= 3;
		d /= 4;

		12 % 4;
		f %= 4;

		# comment at the end 1
		# comment at the end 2
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
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
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "result"},
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
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.RBRACE, "}"},

		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.INT, "9"},
		{token.NOT_EQ, "!="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.STRING, `hello \"world\"`},
		{token.STRING, `hello\n world`},
		{token.STRING, `hello\t\t\tworld`},

		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},

		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.INT, "10"},
		{token.LTE, "<="},
		{token.INT, "12"},
		{token.SEMICOLON, ";"},

		{token.INT, "5"},
		{token.GTE, ">="},
		{token.INT, "3"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "a"},
		{token.PLUS_ASSIGN, "+="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "b"},
		{token.MINUS_ASSIGN, "-="},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "c"},
		{token.ASTERISK_ASSIGN, "*="},
		{token.INT, "3"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "d"},
		{token.SLASH_ASSIGN, "/="},
		{token.INT, "4"},
		{token.SEMICOLON, ";"},

		{token.INT, "12"},
		{token.MODULO, "%"},
		{token.INT, "4"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "f"},
		{token.MODULO_ASSIGN, "%="},
		{token.INT, "4"},
		{token.SEMICOLON, ";"},

		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - token literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
