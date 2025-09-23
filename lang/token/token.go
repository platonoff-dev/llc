package token

type TypeTocken = string

type Token struct {
	Type    TypeTocken
	Literal string
}

const (
	Illegal = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals.
	Ident  = "IDENT"  // add, foobar, x, y, ...
	Int    = "INT"    // 1343456
	String = "STRING" // "<unicode symbols>"

	// Operators.
	Assign   = "="
	Plus     = "+"
	Minus    = "-"
	Bang     = "!"
	Asterisk = "*"
	Slash    = "/"
	LT       = "<"
	GT       = ">"

	// Delimiters.
	Comma     = ","
	Semicolon = ";"
	Colon     = ":"
	LParen    = "("
	RParen    = ")"
	LBrace    = "{"
	RBrace    = "}"
	LBracket  = "["
	RBracket  = "]"

	// Keywords.
	Function = "FUNCTION"
	Let      = "LET"
	True     = "TRUE"
	False    = "FALSE"
	If       = "IF"
	Else     = "ELSE"
	Return   = "RETURN"

	Eq    = "=="
	NotEq = "!="
)

var keywords = map[string]TypeTocken{
	"fn":     Function,
	"let":    Let,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
}

func LookupIndent(indent string) TypeTocken {
	if tok, ok := keywords[indent]; ok {
		return tok
	}
	return Ident
}
