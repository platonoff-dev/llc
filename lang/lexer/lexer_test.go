package lexer

import (
	"testing"

	"llc/lang/token"
)

func TestNextToken(t *testing.T) {
	input := `
let five = 5;
let ten=10;

let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);
!-/*5
5 < 10 > 5

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10
10 != 9
"foobar"
"foo bar"
"❦"
[1, 2, 3]
{"foo": "bar"}
macro(x, y) { x + y };
`
	tests := []struct {
		expectedType    token.TypeTocken
		expectedLiteral string
	}{
		{token.Let, "let"},
		{token.Ident, "five"},
		{token.Assign, "="},
		{token.Int, "5"},
		{token.Semicolon, ";"},

		{token.Let, "let"},
		{token.Ident, "ten"},
		{token.Assign, "="},
		{token.Int, "10"},
		{token.Semicolon, ";"},

		{token.Let, "let"},
		{token.Ident, "add"},
		{token.Assign, "="},
		{token.Function, "fn"},
		{token.LParen, "("},
		{token.Ident, "x"},
		{token.Comma, ","},
		{token.Ident, "y"},
		{token.RParen, ")"},
		{token.LBrace, "{"},
		{token.Ident, "x"},
		{token.Plus, "+"},
		{token.Ident, "y"},
		{token.Semicolon, ";"},
		{token.RBrace, "}"},
		{token.Semicolon, ";"},

		{token.Let, "let"},
		{token.Ident, "result"},
		{token.Assign, "="},
		{token.Ident, "add"},
		{token.LParen, "("},
		{token.Ident, "five"},
		{token.Comma, ","},
		{token.Ident, "ten"},
		{token.RParen, ")"},
		{token.Semicolon, ";"},

		{token.Bang, "!"},
		{token.Minus, "-"},
		{token.Slash, "/"},
		{token.Asterisk, "*"},
		{token.Int, "5"},

		{token.Int, "5"},
		{token.LT, "<"},
		{token.Int, "10"},
		{token.GT, ">"},
		{token.Int, "5"},

		{token.If, "if"},
		{token.LParen, "("},
		{token.Int, "5"},
		{token.LT, "<"},
		{token.Int, "10"},
		{token.RParen, ")"},
		{token.LBrace, "{"},
		{token.Return, "return"},
		{token.True, "true"},
		{token.Semicolon, ";"},
		{token.RBrace, "}"},
		{token.Else, "else"},
		{token.LBrace, "{"},
		{token.Return, "return"},
		{token.False, "false"},
		{token.Semicolon, ";"},
		{token.RBrace, "}"},

		{token.Int, "10"},
		{token.Eq, "=="},
		{token.Int, "10"},

		{token.Int, "10"},
		{token.NotEq, "!="},
		{token.Int, "9"},
		{token.String, "foobar"},

		{token.String, "foo bar"},

		{token.String, "❦"},

		{token.LBracket, "["},
		{token.Int, "1"},
		{token.Comma, ","},
		{token.Int, "2"},
		{token.Comma, ","},
		{token.Int, "3"},
		{token.RBracket, "]"},

		{token.LBrace, "{"},
		{token.String, "foo"},
		{token.Colon, ":"},
		{token.String, "bar"},
		{token.RBrace, "}"},

		{token.Macro, "macro"},
		{token.LParen, "("},
		{token.Ident, "x"},
		{token.Comma, ","},
		{token.Ident, "y"},
		{token.RParen, ")"},
		{token.LBrace, "{"},
		{token.Ident, "x"},
		{token.Plus, "+"},
		{token.Ident, "y"},
		{token.RBrace, "}"},
		{token.Semicolon, ";"},

		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
