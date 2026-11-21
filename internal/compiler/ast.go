package compiler

import "strings"

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
func (p *Program) String() string {
	out := ""
	for _, s := range p.Statements {
		out += s.String() + "\n"
	}
	return out
}

// --- DECLARACIONES ---

type LetStatement struct {
	Token     Token
	Name      *Identifier
	ValueType Token
	Value     Expression
	IsMutable bool
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	mut := "let"
	if ls.IsMutable {
		mut = "mut"
	}
	val := ""
	if ls.Value != nil {
		val = ls.Value.String()
	}
	return "📦 [" + mut + "] " + ls.Name.Value + " : " + ls.ValueType.Literal + " = " + val
}

type DenyStatement struct {
	Token     Token
	Condition string
}

func (ds *DenyStatement) statementNode()       {}
func (ds *DenyStatement) TokenLiteral() string { return ds.Token.Literal }
func (ds *DenyStatement) String() string       { return "🛡️ LEY [DENY] -> " + ds.Condition }

type ExpressionStatement struct {
	Token      Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// --- EXPRESIONES ---

type Identifier struct {
	Token Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegerLiteral struct {
	Token Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type StringLiteral struct {
	Token Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return "\"" + sl.Token.Literal + "\"" }

type CallExpression struct {
	Token     Token      // El paréntesis '('
	Function  Expression // El identificador o método (ej. puts, sys.Imprimir)
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	return "⚙️  LLAMADA: " + ce.Function.String() + "(" + strings.Join(args, ", ") + ")"
}
