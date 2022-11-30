package parser

import (
	"errors"
	"fmt"

	"govega/vega/scanner"
	"govega/vega/token"
)

// TODO: include typing system

// Parser stores needed objects to keep track during the parsing
type Parser struct {
	scanner            *scanner.Scanner
	scanError          error
	lineBreakDelimiter bool           // flag to disrupt line break skipping for delimiter character
	nextToken          *scanner.Token // next token read by looking a head
	currentToken       *scanner.Token // current token which is being analyzed
	//	table              *utils.SymbolTable // symbolTable to store information about recognized identifiers
}

// NewParser generates a new Parser interface
func NewParser(file token.File) *Parser {
	return &Parser{
		scanner:            scanner.NewScanner(file),
		lineBreakDelimiter: false,
		currentToken:       nil,
		scanError:          nil,
		nextToken:          nil,
		//		table:              utils.NewSymbolTable(),
	}
}

// getToken gets token from scanner
func (parser *Parser) getToken() (*scanner.Token, error) {
	var (
		err error
		tok *scanner.Token
	)
	for {
		if tok, err = parser.scanner.Scan(); err != nil {
			return nil, err
		}
		if !parser.lineBreakDelimiter {
			if tok.Token != token.LINEBREAK {
				break
			}
		} else {
			parser.lineBreakDelimiter = false
			break
		}
	}
	return tok, nil
}

// readToken retrieves new token from scanner. First the current token is being updated with the previous next token and
// then the nextToken is updated.
func (parser *Parser) readToken() error {
	var (
		err error
	)
	if parser.scanError != nil {
		return parser.scanError
	}
	if parser.nextToken == nil {
		if parser.currentToken, err = parser.getToken(); err != nil {
			return err
		}
	} else {
		parser.currentToken = parser.nextToken
	}
	if parser.nextToken, err = parser.getToken(); err != nil {
		return err
	}
	return nil
}

func (parser *Parser) test() {
	parser.lineBreakDelimiter = false
}

// lookAHead compares a given tag with the next token, only update nextToken when previously match had cleared nextToken
func (parser *Parser) lookAHead(token token.Token) bool {
	return parser.nextToken.Token == token
}

// matchToken compares a given token with the currentToken
func (parser *Parser) matchToken(token token.Token) bool {
	if err := parser.readToken(); err != nil {
		parser.scanError = err
		return false
	}
	if parser.currentToken.Token != token {
		return false
	}
	return true
}

// syntaxError returns a vega error during on invalid syntax
func (parser *Parser) syntaxError(errorMessage string) error {
	if parser.scanError != nil {
		return parser.scanError
	}
	if parser.currentToken.Token == token.EOF {
		return newParserSyntaxError(unexpectedEOF, parser.currentToken, "Unexpected End Of File", parser.scanner.LineFeed)
	}
	errMsg := fmt.Sprintf(errorMessage, parser.currentToken.Literal)
	return newParserSyntaxError(invalidSyntax, parser.currentToken, errMsg, parser.scanner.LineFeed)
}

// Parse starts parsing process. All functiones which are validating the grammar are using the Parser interface to make
// testing easier
func (parser *Parser) Parse() error {
	return parser.parseBlock()
}

// parseBlock parses block statements
//
// block
//   : (FUNC IDENT LBRACKET functionParamDeclaration? RBRACKET functionReturnType scopeStatement)+ EOF
//   ;
func (parser *Parser) parseBlock() error {
	if !parser.matchToken(token.FUNC) {
		return parser.syntaxError("Missing 'func' at '%v'")
	}
	if !parser.matchToken(token.IDENT) {
		return parser.syntaxError("Mismatched input '%v', expected <identifier>")
	}
	if !parser.matchToken(token.LBRACKET) {
		return parser.syntaxError("Mismatched input '%v', expected '('")
	}
	if parser.lookAHead(token.IDENT) || parser.lookAHead(token.LSBRACKET) {
		if err := parser.parseFunctionParamDeclaration(); err != nil {
			return err
		}
	}
	if !parser.matchToken(token.RBRACKET) {
		return parser.syntaxError("Mismatched input '%v', expected <terminal_variable_type> or ')'")
	}
	if err := parser.parseFunctionReturnType(); err != nil {
		return err
	}
	if err := parser.parseScope(); err != nil {
		return err
	}
	if parser.lookAHead(token.FUNC) {
		return parser.parseBlock() // !!! Declaration Stack !!!
	}
	if !parser.matchToken(token.EOF) {
		return parser.syntaxError("Extraneous input '%v', expected EOF or 'func'")
	}
	return nil
}

// parseFunctionParamDeclaration parses function parameter list
//
// functionParameterDeclaration
//   : functionParameterDefinition (COMMA functionParameterDeclaration)*
//   ;
func (parser *Parser) parseFunctionParamDeclaration() error {
	if err := parser.parseFunctionParamDefinition(); err != nil {
		return err
	}
	for parser.lookAHead(token.COMMA) {
		if !parser.matchToken(token.COMMA) {
			return parser.syntaxError("scanError")
		}
		if parser.lookAHead(token.IDENT) || parser.lookAHead(token.LSBRACKET) {
			if err := parser.parseFunctionParamDeclaration(); err != nil {
				return err
			}
		} else {
			_ = parser.matchToken(-1)
			return parser.syntaxError("Mismatched input '%v', expected <terminal_variable_type>")
		}
	}
	// exit here with error when next token is not ')'
	if !parser.lookAHead(token.RBRACKET) {
		_ = parser.matchToken(-1)
		return parser.syntaxError("Mismatched input '%v', expected ',' or ')'")
	}
	return nil
}

// parseFunctionParamDefinition parse function parameter definition
//
// functionParameterDefinition
//   : (LARRAY RARRAY)* terminalVariableType IDENT
//   ;
func (parser *Parser) parseFunctionParamDefinition() error {
	for parser.lookAHead(token.LSBRACKET) {
		if !parser.matchToken(token.LSBRACKET) {
			return parser.syntaxError("scanError")
		}
		if !parser.matchToken(token.RSBRACKET) {
			return parser.syntaxError("Extraneous input '%v', expected ']'")
		}
	}
	if err := parser.parseTerminalVariableType(); err != nil {
		return err
	}
	if !parser.matchToken(token.IDENT) {
		return parser.syntaxError("Mismatched input '%v', expected <identifier>")
	}
	return nil
}

// parseFunctionReturnType parse return type of function and sets the function identifier symbol type
//
// functionReturnType
//   : terminalVariableType (LARRAY RARRAY)*
//   ;
func (parser *Parser) parseFunctionReturnType() error {
	for parser.lookAHead(token.LSBRACKET) {
		if !parser.matchToken(token.LSBRACKET) {
			return parser.syntaxError("scanError")
		}
		if !parser.matchToken(token.RSBRACKET) {
			return parser.syntaxError("Mismatched input '%v', expected ']'")
		}
	}
	if err := parser.parseTerminalVariableType(); err != nil {
		return err
	}
	return nil
}

// parseTerminalVariableType parse basic variable type terminals
//
// terminalVariableType
//   : INT_TYPE
//   | FLOAT_TYPE
//   | CHAR_TYPE
//   | BOOL_TYPE
//   | STRING_TYPE
//   ;
func (parser *Parser) parseTerminalVariableType() error {
	if !parser.matchToken(token.IDENT) {
		return parser.syntaxError("Mismatched input '%v', expected <variable_type>")
	}
	// TODO: implement typing
	return nil
}

// parseDelimiter parses delimiter characters
//
// delimiter
//   : delimiterCharacters
//   ;
func (parser *Parser) parseDelimiter() error {
	switch {
	case parser.lookAHead(token.DELIMITER):
		if !parser.matchToken(token.DELIMITER) {
			return parser.syntaxError("scanError")
		}
	case parser.lookAHead(token.LINEBREAK):
		if !parser.matchToken(token.LINEBREAK) {
			return parser.syntaxError("scanError")
		}
	default:
		_ = parser.matchToken(-1)
		return parser.syntaxError("Mismatched input '%v', expected ';' or line break")
	}
	return nil
}

// parseScope parses scopes
//
// scopeStatement
//   : LCURLY PASS delimiter | (statement)+ RCURLY
//   ;
func (parser *Parser) parseScope() error {
	if !parser.matchToken(token.LCBRACKET) {
		return parser.syntaxError("Mismatched input '%v', expected '{'")
	}
	// statement: PASS delimiter
	if parser.lookAHead(token.PASS) {
		parser.lineBreakDelimiter = true
		if !parser.matchToken(token.PASS) {
			return parser.syntaxError("scanError")
		}
		if err := parser.parseDelimiter(); err != nil {
			return err
		}
	} else {
		// error on empty body
		if parser.lookAHead(token.RCBRACKET) {
			_ = parser.matchToken(-1)
			return parser.syntaxError("Mismatched input '%v', expected 'pass;' or <statement>")
		}
		// at least one statement has to be defined
		if err := parser.parseStatement(); err != nil {
			if err.Error() == "StatementNotDefined" {
				_ = parser.matchToken(-1)
				return parser.syntaxError("Mismatched input '%v', expected 'pass;' or <statement>")
			}
			return err
		}
		for !parser.lookAHead(token.RCBRACKET) {
			if err := parser.parseStatement(); err != nil {
				if err.Error() == "StatementNotDefined" {
					_ = parser.matchToken(-1)
					return parser.syntaxError("Mismatched input '%v', expected <statement> or '}'")
				}
				return err
			}
		}
	}
	if !parser.matchToken(token.RCBRACKET) {
		return parser.syntaxError("Mismatched input '%v', expected '}'")
	}
	return nil
}

func (parser *Parser) parseDeclaration() error {
	switch {
	case parser.lookAHead(token.CONST):
		if !parser.matchToken(token.CONST) {
			return parser.syntaxError("scanError")
		}
	case parser.lookAHead(token.VAR):
		if !parser.matchToken(token.VAR) {
			return parser.syntaxError("scanError")
		}
	}
	for parser.lookAHead(token.LSBRACKET) {
		if !parser.matchToken(token.LSBRACKET) {
			return parser.syntaxError("scanError")
		}
		if !parser.matchToken(token.INT) {
			return parser.syntaxError("Mismatched input '%v', expected <INT>")
		}
		if !parser.matchToken(token.RSBRACKET) {
			return parser.syntaxError("Mismatched input '%v', expected ']'")
		}
	}
	return parser.parseTerminalVariableType()
}

// parseStatement parses normal statements
//
// statement
//   : CONST? terminalVariableType (LARRAY INT RARRAY)* IDENT (ASSIGN booleanExpression)? DELIMITER
//   |  IDENT (
//  		(LBRACKET ( booleanExpression (COMMA booleanExpression)* )? RBRACKET)
//  	  | (arrayAccess* ASSIGN booleanExpression)
// 	    ) delimiter
//   |  RETURN booleanExpression delimiter
//   |  CONTINUE delimiter
//   |  BREAK delimiter
//   |  WHILE conditionalScope
//   |  IF conditionalScope (ELIF conditionalScope)* (ELSE scopeStatement)?
//   |  SWITCH booleanExpression LCURLY (CASE terminal COLON statement+)+ (DEFAULT COLON statement+)? RCURLY
// ;
func (parser *Parser) parseStatement() error {
	switch {
	// statement: CONTINUE delimiter
	case parser.lookAHead(token.CONTINUE):
		parser.lineBreakDelimiter = true
		if !parser.matchToken(token.CONTINUE) {
			return parser.syntaxError("scanError")
		}
		return parser.parseDelimiter()
	// statement: BREAK delimiter
	case parser.lookAHead(token.BREAK):
		parser.lineBreakDelimiter = true
		if !parser.matchToken(token.BREAK) {
			return parser.syntaxError("scanError")
		}
		return parser.parseDelimiter()
	// statement: IF conditionalScope (ELIF conditionalScope)* (ELSE scopeStatement)?
	case parser.lookAHead(token.IF):
		if !parser.matchToken(token.IF) {
			return parser.syntaxError("scanError")
		}
		if err := parser.parseConditionalScope(); err != nil {
			return err
		}
		for parser.lookAHead(token.ELIF) {
			if !parser.matchToken(token.ELIF) {
				return parser.syntaxError("scanError")
			}
			if err := parser.parseConditionalScope(); err != nil {
				return err
			}
		}
		if parser.lookAHead(token.ELSE) {
			if !parser.matchToken(token.ELSE) {
				return parser.syntaxError("scanError")
			}
			if err := parser.parseScope(); err != nil {
				return err
			}
		}
		return nil
	// statement: SWITCH expression LCURLY (CASE terminal COLON statement+)+ (DEFAULT COLON statement+)? RCURLY
	case parser.lookAHead(token.SWITCH):
		if !parser.matchToken(token.SWITCH) {
			return parser.syntaxError("scanError")
		}
		if err := parser.parseExpression(); err != nil {
			return err
		}
		if !parser.matchToken(token.LCBRACKET) {
			return parser.syntaxError("Mismatched input '%v', expected '{'")
		}
		if !parser.matchToken(token.CASE) {
			return parser.syntaxError("Mismatched input '%v', expected 'case' or 'default'")
		}
		if err := parser.parseTerminal(); err != nil {
			return err
		}
		if !parser.matchToken(token.COLON) {
			return parser.syntaxError("Mismatched input '%v', expected ':'")
		}
		if err := parser.parseStatement(); err != nil {
			if err.Error() == "StatementNotDefined" {
				_ = parser.matchToken(-1)
				return parser.syntaxError("Mismatched input '%v', expected <statement>")
			}
			return err
		}
		for !parser.lookAHead(token.RCBRACKET) && !parser.lookAHead(token.CASE) && !parser.lookAHead(token.DEFAULT) {
			if err := parser.parseStatement(); err != nil {
				if err.Error() == "StatementNotDefined" {
					_ = parser.matchToken(-1)
					return parser.syntaxError("Mismatched input '%v', expected <statement>, another 'case' or 'default' keyword or '}'")
				}
				return err
			}
		}
		for parser.lookAHead(token.CASE) {
			if !parser.matchToken(token.CASE) {
				return parser.syntaxError("scanError")
			}
			if err := parser.parseTerminal(); err != nil {
				return err
			}
			if !parser.matchToken(token.COLON) {
				return parser.syntaxError("Mismatched input '%v', expected ':'")
			}
			if err := parser.parseStatement(); err != nil {
				if err.Error() == "StatementNotDefined" {
					_ = parser.matchToken(-1)
					return parser.syntaxError("Mismatched input '%v', expected <statement>")
				}
				return err
			}
			for !parser.lookAHead(token.RCBRACKET) && !parser.lookAHead(token.CASE) && !parser.lookAHead(token.DEFAULT) {
				if err := parser.parseStatement(); err != nil {
					if err.Error() == "StatementNotDefined" {
						_ = parser.matchToken(-1)
						return parser.syntaxError("Mismatched input '%v', expected <statement>, another 'case' or 'default' keyword or '}'")
					}
					return err
				}
			}
		}
		if parser.lookAHead(token.DEFAULT) {
			if !parser.matchToken(token.DEFAULT) {
				return parser.syntaxError("scanError")
			}
			if !parser.matchToken(token.COLON) {
				return parser.syntaxError("Mismatched input '%v', expected ':'")
			}
			if err := parser.parseStatement(); err != nil {
				if err.Error() == "StatementNotDefined" {
					_ = parser.matchToken(-1)
					return parser.syntaxError("Mismatched input '%v', expected <statement>")
				}
				return err
			}
			for !parser.lookAHead(token.RCBRACKET) {
				if err := parser.parseStatement(); err != nil {
					if err.Error() == "StatementNotDefined" {
						_ = parser.matchToken(-1)
						return parser.syntaxError("Mismatched input '%v', expected <statement> or '}'")
					}
					return err
				}
			}
		}
		if !parser.matchToken(token.RCBRACKET) {
			return parser.syntaxError("scanError")
		}
		return nil
	// statement: WHILE conditionalScope
	case parser.lookAHead(token.WHILE):
		if !parser.matchToken(token.WHILE) {
			return parser.syntaxError("scanError")
		}
		return parser.parseConditionalScope()
	// statement: RETURN booleanExpression delimiter
	case parser.lookAHead(token.RETURN):
		if !parser.matchToken(token.RETURN) {
			return parser.syntaxError("scanError")
		}
		if err := parser.parseBooleanExpression(); err != nil {
			return err
		}
		return parser.parseDelimiter()
	// CONST? terminalVariableType (LARRAY INT RARRAY)* IDENT (ASSIGN booleanExpression)? delimiter
	// TODO fallthrough to enable CONST var ident | ident())
	case parser.lookAHead(token.CONST), parser.lookAHead(token.VAR):
		if err := parser.parseDeclaration(); err != nil {
			return err
		}
		parser.lineBreakDelimiter = true
		if !parser.matchToken(token.IDENT) {
			return parser.syntaxError("Mismatched input '%v', expected <identifier>")
		}
		if parser.lookAHead(token.ASSIGN) {
			if !parser.matchToken(token.ASSIGN) {
				return parser.syntaxError("scanError")
			}
			if err := parser.parseBooleanExpression(); err != nil {
				return err
			}
		}
		return parser.parseDelimiter()
	// IDENT ((LBRACKET ( booleanExpression (COMMA booleanExpression)* )? RBRACKET) | (arrayAccess* ASSIGN booleanExpression)) delimiter
	case parser.lookAHead(token.IDENT):
		if !parser.matchToken(token.IDENT) {
			return parser.syntaxError("scanError")
		}
		switch {
		// function call
		case parser.lookAHead(token.LBRACKET):
			if !parser.matchToken(token.LBRACKET) {
				return parser.syntaxError("scanError")
			}
			if !parser.lookAHead(token.RBRACKET) {
				if err := parser.parseBooleanExpression(); err != nil {
					return err
				}
				for parser.lookAHead(token.COMMA) {
					if !parser.matchToken(token.COMMA) {
						return parser.syntaxError("Mismatched input '%v', expected ',' or ')'")
					}
					if err := parser.parseBooleanExpression(); err != nil {
						return err
					}
				}
			}
			parser.lineBreakDelimiter = true
			if !parser.matchToken(token.RBRACKET) {
				return parser.syntaxError("Mismatched input '%v', expected ',' or ')'")
			}
		// array definiton
		case parser.lookAHead(token.LSBRACKET):
			for parser.lookAHead(token.LSBRACKET) {
				if err := parser.parseArrayAccess(); err != nil {
					return err
				}
			}
			fallthrough
		case parser.lookAHead(token.ASSIGN):
			if !parser.matchToken(token.ASSIGN) {
				return parser.syntaxError("Mismatched input '%v', expected '[', or '='")
			}
			if err := parser.parseBooleanExpression(); err != nil {
				return err
			}
		default:
			_ = parser.matchToken(-1)
			return parser.syntaxError("Mismatched input '%v', expected '(', '[', or '='")
		}
		return parser.parseDelimiter()
	default:
		return errors.New("StatementNotDefined")
	}
}

// conditionalScope
//   : booleanExpression scopeStatement
//   ;
func (parser *Parser) parseConditionalScope() error {
	if err := parser.parseBooleanExpression(); err != nil {
		return err
	}
	return parser.parseScope()
}

// booleanExpression
//   :   comparisonExpression ((OR | AND) comparisonExpression)*
//   ;
func (parser *Parser) parseBooleanExpression() error {
	if err := parser.parseComparisonExpression(); err != nil {
		return err
	}
	var loopControl = false
	for {
		switch {
		case parser.lookAHead(token.OR):
			if !parser.matchToken(token.OR) {
				return parser.syntaxError("scanError")
			}
		case parser.lookAHead(token.LOR):
			if !parser.matchToken(token.LOR) {
				return parser.syntaxError("scanError")
			}
		case parser.lookAHead(token.AND):
			if !parser.matchToken(token.AND) {
				return parser.syntaxError("scanError")
			}
		case parser.lookAHead(token.LAND):
			if !parser.matchToken(token.LAND) {
				return parser.syntaxError("scanError")
			}
		default:
			loopControl = true
		}
		if loopControl {
			break
		}
		if err := parser.parseComparisonExpression(); err != nil {
			return err
		}
	}
	return nil
}

// comparisonExpression
// :   expression (comparisonOperator expression)*
// ;
func (parser *Parser) parseComparisonExpression() error {
	if err := parser.parseExpression(); err != nil {
		return err
	}
	var loopControl = false
	for {
		switch {
		case parser.lookAHead(token.EQ):
			if !parser.matchToken(token.EQ) {
				return parser.syntaxError("scanError")
			}
		case parser.lookAHead(token.NE):
			if !parser.matchToken(token.NE) {
				return parser.syntaxError("scanError")
			}
		case parser.lookAHead(token.GE):
			if !parser.matchToken(token.GE) {
				return parser.syntaxError("scanError")
			}
		case parser.lookAHead(token.GREATER):
			if !parser.matchToken(token.GREATER) {
				return parser.syntaxError("scanError")
			}
		case parser.lookAHead(token.LE):
			if !parser.matchToken(token.LE) {
				return parser.syntaxError("scanError")
			}
		case parser.lookAHead(token.LESS):
			if !parser.matchToken(token.LESS) {
				return parser.syntaxError("scanError")
			}
		default:
			loopControl = true
		}
		if loopControl {
			break
		}
		if err := parser.parseExpression(); err != nil {
			return err
		}
	}
	return nil
}

// arrayAccess
// : LARRAY expression RARRAY
// ;
func (parser *Parser) parseArrayAccess() error {
	if !parser.matchToken(token.LSBRACKET) {
		return parser.syntaxError("scanError")
	}
	if err := parser.parseExpression(); err != nil {
		return err
	}
	parser.lineBreakDelimiter = true
	if !parser.matchToken(token.RSBRACKET) {
		return parser.syntaxError("Mismatched input '%v', expected ']'")
	}
	return nil
}

// expression
// : term ((PLUS | MINUS) term)*
// ;
func (parser *Parser) parseExpression() error {
	if err := parser.parseTerm(); err != nil {
		return err
	}
	var loopControl = false
	for {
		switch {
		case parser.lookAHead(token.ADD):
			if !parser.matchToken(token.ADD) {
				return parser.syntaxError("scanError")
			}
		case parser.lookAHead(token.SUB):
			if !parser.matchToken(token.SUB) {
				return parser.syntaxError("scanError")
			}
		default:
			loopControl = true
		}
		if loopControl {
			break
		}
		if err := parser.parseTerm(); err != nil {
			return err
		}
	}
	return nil
}

// term
// : factor ((MUL | DIV) factor)*
// ;
func (parser *Parser) parseTerm() error {
	if err := parser.parseFactor(); err != nil {
		return err
	}
	var loopControl = false
	for loopControl {
		switch {
		case parser.lookAHead(token.MUL):
			if !parser.matchToken(token.MUL) {
				return parser.syntaxError("scanError")
			}
		case parser.lookAHead(token.DIV):
			if !parser.matchToken(token.DIV) {
				return parser.syntaxError("scanError")
			}
		default:
			loopControl = true
		}
		if loopControl {
			break
		}
		if err := parser.parseFactor(); err != nil {
			return err
		}
	}
	return nil
}

// factor
// : (MINUS | NOT)? unary
// ;
func (parser *Parser) parseFactor() error {
	switch {
	case parser.lookAHead(token.LNOT):
		if !parser.matchToken(token.LNOT) {
			return parser.syntaxError("scanError")
		}
	case parser.lookAHead(token.NOT):
		if !parser.matchToken(token.NOT) {
			return parser.syntaxError("scanError")
		}
	}
	return parser.parseUnary()
}

// unary
// : (BASIC | TRUE | FALSE | LITERAL)
// | IDENT arrayAccess?
// | IDENT LBRACKET ( booleanExpression (COMMA booleanExpression)* )? RBRACKET   // func call
// | LBRACKET booleanExpression RBRACKET
// | LARRAY (expression (COMMA expression)* )? RARRAY           // set array value
// ;
func (parser *Parser) parseUnary() error {
	switch {
	// IDENT (arrayAccess?) | IDENT LBRACKET ( expression (COMMA expression)* )? RBRACKET
	case parser.lookAHead(token.IDENT):
		if !parser.matchToken(token.IDENT) {
			return parser.syntaxError("scanError")
		}
		// arrayAccess?
		if parser.lookAHead(token.LSBRACKET) {
			if err := parser.parseArrayAccess(); err != nil {
				return err
			}
			// LBRACKET ( booleanExpression (COMMA booleanExpression)* )? RBRACKET
		} else if parser.lookAHead(token.LBRACKET) {
			if !parser.matchToken(token.LBRACKET) {
				return parser.syntaxError("scanError")
			}
			if !parser.lookAHead(token.RBRACKET) {
				if err := parser.parseBooleanExpression(); err != nil {
					return err
				}
				for parser.lookAHead(token.COMMA) {
					if !parser.matchToken(token.COMMA) {
						return parser.syntaxError("scanError")
					}
					if err := parser.parseBooleanExpression(); err != nil {
						return err
					}
				}
			}
			parser.lineBreakDelimiter = true
			if !parser.matchToken(token.RBRACKET) {
				return parser.syntaxError("Mismatched input '%v', expected ',' or ')'")
			}
		}
	// LBRACKET booleanExpression RBRACKET
	case parser.lookAHead(token.LBRACKET):
		if !parser.matchToken(token.LBRACKET) {
			return parser.syntaxError("scanError")
		}
		if err := parser.parseBooleanExpression(); err != nil {
			return err
		}
		parser.lineBreakDelimiter = true
		if !parser.matchToken(token.RBRACKET) {
			return parser.syntaxError("Mismatched input '%v', expected ')'")
		}
	// LARRAY (expression (COMMA expression)* )? RARRAY
	case parser.lookAHead(token.LSBRACKET):
		if !parser.matchToken(token.LSBRACKET) {
			return parser.syntaxError("scanError")
		}
		if err := parser.parseExpression(); err != nil {
			return err
		}
		for parser.lookAHead(token.COMMA) {
			if !parser.matchToken(token.COMMA) {
				return parser.syntaxError("scanError")
			}
			if err := parser.parseExpression(); err != nil {
				return err
			}
		}
		parser.lineBreakDelimiter = true
		if !parser.matchToken(token.RSBRACKET) {
			return parser.syntaxError("Mismatched input '%v', expected ',' or ']'")
		}
	//
	default:
		if err := parser.parseTerminal(); err != nil {
			return parser.syntaxError("Mismatched input '%v', expected <unary>")
		}
	}
	return nil
}

// terminal
//   : NUM
//   | FLOAT
//   | TRUE
//   | FALSE
//   | LITERAL
//   ;
func (parser *Parser) parseTerminal() error {
	parser.lineBreakDelimiter = true
	switch {
	case parser.lookAHead(token.INT):
		if !parser.matchToken(token.INT) {
			return parser.syntaxError("scanError")
		}
	case parser.lookAHead(token.FLOAT):
		if !parser.matchToken(token.FLOAT) {
			return parser.syntaxError("scanError")
		}
	case parser.lookAHead(token.TRUE):
		if !parser.matchToken(token.TRUE) {
			return parser.syntaxError("scanError")
		}
	case parser.lookAHead(token.FALSE):
		if !parser.matchToken(token.FALSE) {
			return parser.syntaxError("scanError")
		}
	case parser.lookAHead(token.CHAR):
		if !parser.matchToken(token.CHAR) {
			return parser.syntaxError("scanError")
		}
	case parser.lookAHead(token.STRING):
		if !parser.matchToken(token.STRING) {
			return parser.syntaxError("scanError")
		}
	default:
		_ = parser.matchToken(-1)
		return parser.syntaxError("Mismatched input '%v', expected <terminal>")
	}
	return nil
}
