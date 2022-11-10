package frontend

import "govega/vega/language/tokens"

type Vega interface {
	NewLexer(code []byte) Lexer
	NewParser(lexer Lexer) Parser
}

// Parser interface which allows better testing capacities
type Parser interface {
	Parse(p Parser) error
	parseBlock(p Parser) error
	parseFunctionParamDeclaration(p Parser) error
	parseFunctionParamDefinition(p Parser) error
	parseFunctionReturnType(p Parser) error
	parseArrayAccess(p Parser) error
	parseTerminalVariableType() error
	parseScope(p Parser) error
	parseStatement(p Parser) error
	parseConditionalScope(p Parser) error
	parseBooleanExpression(p Parser) error
	parseComparisonExpression(p Parser) error
	parseExpression(p Parser) error
	parseTerm(p Parser) error
	parseFactor(p Parser) error
	parseUnary(p Parser) error
}

type Lexer interface {
	getLineFeed() string
	scan() (*lexicalToken, error)
	newLexicalToken(token tokens.IToken) *lexicalToken
}
