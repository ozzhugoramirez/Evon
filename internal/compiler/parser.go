package compiler

import (
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	CALL // Prioridad para paréntesis ()
)

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

type Parser struct {
	l      *Lexer
	errors []string

	curToken  Token
	peekToken Token

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.registerPrefix(IDENT, p.parseIdentifier)
	p.registerPrefix(NUMBER, p.parseIntegerLiteral)
	p.registerPrefix(STRING, p.parseStringLiteral)

	p.infixParseFns = make(map[TokenType]infixParseFn)
	p.registerInfix(LPAREN, p.parseCallExpression) // Registramos el paréntesis para llamadas

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string { return p.errors }

func (p *Parser) peekPrecedence() int {
	if p.peekToken.Type == LPAREN {
		return CALL
	}
	return LOWEST
}
func (p *Parser) curPrecedence() int {
	if p.curToken.Type == LPAREN {
		return CALL
	}
	return LOWEST
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("Error sintáctico (Línea %d): esperaba '%s', encontré '%s'", p.peekToken.Line, t, p.peekToken.Literal)
	p.errors = append(p.errors, msg)
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}
	for p.curToken.Type != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case LET:
		return p.parseVariableStatement(false)
	case MUT:
		return p.parseVariableStatement(true)
	case DENY:
		return p.parseDenyStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.errors = append(p.errors, fmt.Sprintf("Error (Línea %d): no sé cómo evaluar el símbolo '%s'", p.curToken.Line, p.curToken.Literal))
		return nil
	}
	leftExp := prefix()

	for p.peekToken.Type != EOF && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() Expression {
	lit := &IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("no se pudo leer %q como número", p.curToken.Literal))
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseVariableStatement(isMutable bool) *LetStatement {
	stmt := &LetStatement{Token: p.curToken, IsMutable: isMutable}
	if !p.expectPeek(IDENT) {
		return nil
	}
	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(COLON) {
		return nil
	}
	p.nextToken()

	if p.curToken.Type != T_INT && p.curToken.Type != T_STRING && p.curToken.Type != T_BOOL && p.curToken.Type != T_LIST && p.curToken.Type != T_FILE {
		p.errors = append(p.errors, fmt.Sprintf("Error de tipado (Línea %d): '%s' no es un tipo válido.", p.curToken.Line, p.curToken.Literal))
		return nil
	}
	stmt.ValueType = p.curToken

	if !p.expectPeek(ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseDenyStatement() *DenyStatement {
	stmt := &DenyStatement{Token: p.curToken}
	p.nextToken()
	condicion := ""
	for p.curToken.Type != EOF {
		condicion += p.curToken.Literal + " "
		p.nextToken()
	}
	stmt.Condition = condicion
	return stmt
}

func (p *Parser) parseCallExpression(function Expression) Expression {
	exp := &CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []Expression {
	args := []Expression{}
	if p.peekToken.Type == RPAREN {
		p.nextToken()
		return args
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekToken.Type == COMMA {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(RPAREN) {
		return nil
	}
	return args
}
