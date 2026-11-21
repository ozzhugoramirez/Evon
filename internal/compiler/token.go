package compiler

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	NUMBER = "NUMBER"
	STRING = "STRING"

	// Operadores
	ASSIGN   = "="
	EQ       = "=="
	NOT_EQ   = "!="
	LT       = "<"
	GT       = ">"
	LTE      = "<="
	GTE      = ">="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	REDIRECT = ">>"
	ARROW    = "->"
	DOT      = "."
	COLON    = ":"
	LPAREN   = "("
	RPAREN   = ")"
	COMMA    = "," // Nuevo operador para separar argumentos en funciones

	// Palabras Clave Evon
	ENTITY = "ENTITY"
	LAW    = "LAW"
	FLOW   = "FLOW"
	LET    = "LET"
	MUT    = "MUT"

	T_STRING = "T_STRING"
	T_INT    = "T_INT"
	T_BOOL   = "T_BOOL"
	T_LIST   = "T_LIST"
	T_FILE   = "T_FILE"

	ENSURE = "ENSURE"
	DENY   = "DENY"
	MATCH  = "MATCH"

	SELF  = "SELF"
	EMPTY = "EMPTY"
	HAS   = "HAS"
	IN    = "IN"

	TRUE  = "TRUE"
	FALSE = "FALSE"
	NULL  = "NULL"
	ANY   = "ANY"
)

var keywords = map[string]TokenType{
	"entity": ENTITY,
	"law":    LAW,
	"flow":   FLOW,
	"let":    LET,
	"mut":    MUT,
	"String": T_STRING,
	"Int":    T_INT,
	"Bool":   T_BOOL,
	"List":   T_LIST,
	"File":   T_FILE,
	"ensure": ENSURE,
	"deny":   DENY,
	"match":  MATCH,
	"self":   SELF,
	"empty":  EMPTY,
	"has":    HAS,
	"in":     IN,
	"true":   TRUE,
	"false":  FALSE,
	"null":   NULL,
	"any":    ANY,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
