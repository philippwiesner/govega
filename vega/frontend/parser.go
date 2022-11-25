package frontend

import (
	"errors"
	"fmt"

	"govega/vega/frontend/utils"
	"govega/vega/language"
	"govega/vega/language/tokens"
)

// TODO: include typing system

// parser stores needed objects to keep track during the parsing
type parser struct {
	*vega
	lexer              Lexer
	lexicalError       error
	lineBreakDelimiter bool               // flag to disrupt line break skipping for delimiter character
	nextToken          *lexicalToken      // next token read by looking a head
	currentToken       *lexicalToken      // current token which is being analyzed
	table              *utils.SymbolTable // symbolTable to store information about recognized identifiers
}

// NewParser generates a new Parser interface
func (v *vega) NewParser(lexer Lexer) Parser {
	var parser Parser = &parser{
		vega:               v,
		lexer:              lexer,
		lineBreakDelimiter: false,
		currentToken:       nil,
		lexicalError:       nil,
		nextToken:          nil,
		table:              utils.NewSymbolTable(),
	}
	return parser
}

// getToken gets token from lexer
func (parser *parser) getToken() (*lexicalToken, error) {
	var (
		err   error
		token *lexicalToken
	)
	for {
		if token, err = parser.lexer.scan(); err != nil {
			return nil, err
		}
		if !parser.lineBreakDelimiter {
			if token.GetTag() != tokens.LINEBREAK {
				break
			}
		} else {
			parser.lineBreakDelimiter = false
			break
		}
	}
	return token, nil
}

// readToken retrieves new tokens from lexer. First the current token is being updated with the previous next token and
// then the nextToken is updated.
func (parser *parser) readToken() error {
	var (
		err error
	)
	if parser.lexicalError != nil {
		return parser.lexicalError
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

// lookAHead compares a given tag with the next token, only update nextToken when previously match had cleared nextToken
func (parser *parser) lookAHead(tag int) bool {
	return parser.nextToken.GetTag() == tag
}

// matchToken compares a given token with the currentToken
func (parser *parser) matchToken(tag int) bool {
	if err := parser.readToken(); err != nil {
		parser.lexicalError = err
		return false
	}
	if parser.currentToken.GetTag() != tag {
		return false
	}
	return true
}

// syntaxError returns a vega error during on invalid syntax
func (parser *parser) syntaxError(errorMessage string) error {
	if parser.lexicalError != nil {
		return parser.lexicalError
	}
	if parser.currentToken.GetTag() == tokens.EOF {
		return parser.newParserSyntaxError(unexpectedEOF, parser.currentToken, "Unexpected End Of File", parser.lexer.getLineFeed())
	}
	errMsg := fmt.Sprintf(errorMessage, parser.currentToken.GetToken().String())
	return parser.newParserSyntaxError(invalidSyntax, parser.currentToken, errMsg, parser.lexer.getLineFeed())
}

// Parse starts parsing process. All functiones which are validating the grammar are using the Parser interface to make
// testing easier
func (parser *parser) Parse(parserInterface Parser) error {
	return parser.parseBlock(parserInterface)
}

// parseBlock parses block statements
//
// block
//   : (FUNC ID LBRACKET functionParamDeclaration? RBRACKET functionReturnType scopeStatement)+ EOF
//   ;
func (parser *parser) parseBlock(parserInterface Parser) error {
	if !parser.matchToken(tokens.FUNC) {
		return parser.syntaxError("Missing 'func' at '%v'")
	}
	if !parser.matchToken(tokens.ID) {
		return parser.syntaxError("Mismatched input '%v', expected <identifier>")
	}
	if !parser.matchToken(tokens.LBRACKET) {
		return parser.syntaxError("Mismatched input '%v', expected '('")
	}
	if parser.lookAHead(tokens.BASIC) || parser.lookAHead(tokens.TYPE) {
		if err := parserInterface.parseFunctionParamDeclaration(parserInterface); err != nil {
			return err
		}
	}
	if !parser.matchToken(tokens.RBRACKET) {
		return parser.syntaxError("Mismatched input '%v', expected <terminal_variable_type> or ')'")
	}
	if err := parserInterface.parseFunctionReturnType(parserInterface); err != nil {
		return err
	}
	if err := parserInterface.parseScope(parserInterface); err != nil {
		return err
	}
	if parser.lookAHead(tokens.FUNC) {
		return parserInterface.parseBlock(parserInterface) // !!! Declaration Stack !!!
	}
	if !parser.matchToken(tokens.EOF) {
		return parser.syntaxError("Extraneous input '%v', expected EOF or 'func'")
	}
	return nil
}

// parseFunctionParamDeclaration parses function parameter list
//
// functionParameterDeclaration
//   : functionParameterDefinition (COMMA functionParameterDeclaration)*
//   ;
func (parser *parser) parseFunctionParamDeclaration(parserInterface Parser) error {
	if err := parserInterface.parseFunctionParamDefinition(parserInterface); err != nil {
		return err
	}
	for parser.lookAHead(tokens.COMMA) {
		if !parser.matchToken(tokens.COMMA) {
			return parser.syntaxError("lexicalError")
		}
		if parser.lookAHead(tokens.BASIC) || parser.lookAHead(tokens.TYPE) {
			if err := parserInterface.parseFunctionParamDeclaration(parserInterface); err != nil {
				return err
			}
		} else {
			_ = parser.matchToken(-1)
			return parser.syntaxError("Mismatched input '%v', expected <terminal_variable_type>")
		}
	}
	// exit here with error when next token is not ')'
	if !parser.lookAHead(tokens.RBRACKET) {
		_ = parser.matchToken(-1)
		return parser.syntaxError("Mismatched input '%v', expected ',' or ')'")
	}
	return nil
}

// parseFunctionParamDefinition parse function parameter definition
//
// functionParameterDefinition
//   : terminalVariableType (LARRAY RARRAY)* ID
//   ;
func (parser *parser) parseFunctionParamDefinition(parserInterface Parser) error {
	if err := parserInterface.parseTerminalVariableType(); err != nil {
		return err
	}
	for parser.lookAHead(tokens.LSBRACKET) {
		if !parser.matchToken(tokens.LSBRACKET) {
			return parser.syntaxError("lexicalError")
		}
		if !parser.matchToken(tokens.RSBRACKET) {
			return parser.syntaxError("Extraneous input '%v', expected ']'")
		}
	}
	if !parser.matchToken(tokens.ID) {
		return parser.syntaxError("Mismatched input '%v', expected '[' or <identifier>")
	}
	return nil
}

// parseFunctionReturnType parse return type of function and sets the function identifier symbol type
//
// functionReturnType
//   : terminalVariableType (LARRAY RARRAY)*
//   ;
func (parser *parser) parseFunctionReturnType(parserInterface Parser) error {
	if err := parserInterface.parseTerminalVariableType(); err != nil {
		return err
	}
	for parser.lookAHead(tokens.LSBRACKET) {
		if !parser.matchToken(tokens.LSBRACKET) {
			return parser.syntaxError("lexicalError")
		}
		if !parser.matchToken(tokens.RSBRACKET) {
			return parser.syntaxError("Mismatched input '%v', expected ']'")
		}
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
func (parser *parser) parseTerminalVariableType() error {
	switch {
	case parser.lookAHead(tokens.BASIC):
		if !parser.matchToken(tokens.BASIC) {
			return parser.syntaxError("lexicalError")
		}
		switch parser.currentToken.GetToken().(tokens.IWord).GetLexeme() {
		case language.IntType.GetLexeme():
			return nil
		case language.BoolType.GetLexeme():
			return nil
		case language.FloatType.GetLexeme():
			return nil
		case language.CharType.GetLexeme():
			return nil
		}
		return nil
	case parser.lookAHead(tokens.TYPE):
		if !parser.matchToken(tokens.TYPE) {
			return parser.syntaxError("lexicalError")
		}
		switch parser.currentToken.GetToken().(tokens.IWord).GetLexeme() {
		case "str":
			return nil
		}
		return nil
	default:
		_ = parser.matchToken(-1)
		return parser.syntaxError("Mismatched input '%v', expected <variable_type>")
	}
}

// parseDelimiter parses delimiter characters
//
// delimiter
//   : delimiterCharacters
//   ;
func (parser *parser) parseDelimiter() error {
	switch {
	case parser.lookAHead(tokens.DELIMITER):
		if !parser.matchToken(tokens.DELIMITER) {
			return parser.syntaxError("lexicalError")
		}
	case parser.lookAHead(tokens.LINEBREAK):
		if !parser.matchToken(tokens.LINEBREAK) {
			return parser.syntaxError("lexicalError")
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
func (parser *parser) parseScope(parserInterface Parser) error {
	if !parser.matchToken(tokens.LCBRACKET) {
		return parser.syntaxError("Mismatched input '%v', expected '{'")
	}
	// statement: PASS delimiter
	if parser.lookAHead(tokens.PASS) {
		parser.lineBreakDelimiter = true
		if !parser.matchToken(tokens.PASS) {
			return parser.syntaxError("lexicalError")
		}
		if err := parser.parseDelimiter(); err != nil {
			return err
		}
	} else {
		// error on empty body
		if parser.lookAHead(tokens.RCBRACKET) {
			_ = parser.matchToken(-1)
			return parser.syntaxError("Mismatched input '%v', expected 'pass;' or <statement>")
		}
		// at least one statement has to be defined
		if err := parserInterface.parseStatement(parserInterface); err != nil {
			if err.Error() == "StatementNotDefined" {
				_ = parser.matchToken(-1)
				return parser.syntaxError("Mismatched input '%v', expected 'pass;' or <statement>")
			}
			return err
		}
		for !parser.lookAHead(tokens.RCBRACKET) {
			if err := parserInterface.parseStatement(parserInterface); err != nil {
				if err.Error() == "StatementNotDefined" {
					_ = parser.matchToken(-1)
					return parser.syntaxError("Mismatched input '%v', expected <statement> or '}'")
				}
				return err
			}
		}
	}
	if !parser.matchToken(tokens.RCBRACKET) {
		return parser.syntaxError("Mismatched input '%v', expected '}'")
	}
	return nil
}

// parseStatement parses normal statements
//
// statement
//   : CONST? terminalVariableType (LARRAY INT RARRAY)* ID (ASSIGN booleanExpression)? DELIMITER
//   |  ID (
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
func (parser *parser) parseStatement(parserInterface Parser) error {
	switch {
	// statement: CONTINUE delimiter
	case parser.lookAHead(tokens.CONTINUE):
		parser.lineBreakDelimiter = true
		if !parser.matchToken(tokens.CONTINUE) {
			return parser.syntaxError("lexicalError")
		}
		return parser.parseDelimiter()
	// statement: BREAK delimiter
	case parser.lookAHead(tokens.BREAK):
		parser.lineBreakDelimiter = true
		if !parser.matchToken(tokens.BREAK) {
			return parser.syntaxError("lexicalError")
		}
		return parser.parseDelimiter()
	// statement: IF conditionalScope (ELIF conditionalScope)* (ELSE scopeStatement)?
	case parser.lookAHead(tokens.IF):
		if !parser.matchToken(tokens.IF) {
			return parser.syntaxError("lexicalError")
		}
		if err := parserInterface.parseConditionalScope(parserInterface); err != nil {
			return err
		}
		for parser.lookAHead(tokens.ELIF) {
			if !parser.matchToken(tokens.ELIF) {
				return parser.syntaxError("lexicalError")
			}
			if err := parserInterface.parseConditionalScope(parserInterface); err != nil {
				return err
			}
		}
		if parser.lookAHead(tokens.ELSE) {
			if !parser.matchToken(tokens.ELSE) {
				return parser.syntaxError("lexicalError")
			}
			if err := parserInterface.parseScope(parserInterface); err != nil {
				return err
			}
		}
		return nil
	// statement: SWITCH expression LCURLY (CASE terminal COLON statement+)+ (DEFAULT COLON statement+)? RCURLY
	case parser.lookAHead(tokens.SWITCH):
		if !parser.matchToken(tokens.SWITCH) {
			return parser.syntaxError("lexicalError")
		}
		if err := parserInterface.parseExpression(parserInterface); err != nil {
			return err
		}
		if !parser.matchToken(tokens.LCBRACKET) {
			return parser.syntaxError("Mismatched input '%v', expected '{'")
		}
		if !parser.matchToken(tokens.CASE) {
			return parser.syntaxError("Mismatched input '%v', expected 'case' or 'default'")
		}
		if err := parserInterface.parseTerminal(); err != nil {
			return err
		}
		if !parser.matchToken(tokens.COLON) {
			return parser.syntaxError("Mismatched input '%v', expected ':'")
		}
		if err := parserInterface.parseStatement(parserInterface); err != nil {
			if err.Error() == "StatementNotDefined" {
				_ = parser.matchToken(-1)
				return parser.syntaxError("Mismatched input '%v', expected <statement>")
			}
			return err
		}
		for !parser.lookAHead(tokens.RCBRACKET) && !parser.lookAHead(tokens.CASE) && !parser.lookAHead(tokens.DEFAULT) {
			if err := parserInterface.parseStatement(parserInterface); err != nil {
				if err.Error() == "StatementNotDefined" {
					_ = parser.matchToken(-1)
					return parser.syntaxError("Mismatched input '%v', expected <statement>, another 'case' or 'default' keyword or '}'")
				}
				return err
			}
		}
		for parser.lookAHead(tokens.CASE) {
			if !parser.matchToken(tokens.CASE) {
				return parser.syntaxError("lexicalError")
			}
			if err := parserInterface.parseTerminal(); err != nil {
				return err
			}
			if !parser.matchToken(tokens.COLON) {
				return parser.syntaxError("Mismatched input '%v', expected ':'")
			}
			if err := parserInterface.parseStatement(parserInterface); err != nil {
				if err.Error() == "StatementNotDefined" {
					_ = parser.matchToken(-1)
					return parser.syntaxError("Mismatched input '%v', expected <statement>")
				}
				return err
			}
			for !parser.lookAHead(tokens.RCBRACKET) && !parser.lookAHead(tokens.CASE) && !parser.lookAHead(tokens.DEFAULT) {
				if err := parserInterface.parseStatement(parserInterface); err != nil {
					if err.Error() == "StatementNotDefined" {
						_ = parser.matchToken(-1)
						return parser.syntaxError("Mismatched input '%v', expected <statement>, another 'case' or 'default' keyword or '}'")
					}
					return err
				}
			}
		}
		if parser.lookAHead(tokens.DEFAULT) {
			if !parser.matchToken(tokens.DEFAULT) {
				return parser.syntaxError("lexicalError")
			}
			if !parser.matchToken(tokens.COLON) {
				return parser.syntaxError("Mismatched input '%v', expected ':'")
			}
			if err := parserInterface.parseStatement(parserInterface); err != nil {
				if err.Error() == "StatementNotDefined" {
					_ = parser.matchToken(-1)
					return parser.syntaxError("Mismatched input '%v', expected <statement>")
				}
				return err
			}
			for !parser.lookAHead(tokens.RCBRACKET) {
				if err := parserInterface.parseStatement(parserInterface); err != nil {
					if err.Error() == "StatementNotDefined" {
						_ = parser.matchToken(-1)
						return parser.syntaxError("Mismatched input '%v', expected <statement> or '}'")
					}
					return err
				}
			}
		}
		if !parser.matchToken(tokens.RCBRACKET) {
			return parser.syntaxError("lexicalError")
		}
		return nil
	// statement: WHILE conditionalScope
	case parser.lookAHead(tokens.WHILE):
		if !parser.matchToken(tokens.WHILE) {
			return parser.syntaxError("lexicalError")
		}
		return parserInterface.parseConditionalScope(parserInterface)
	// statement: RETURN booleanExpression delimiter
	case parser.lookAHead(tokens.RETURN):
		if !parser.matchToken(tokens.RETURN) {
			return parser.syntaxError("lexicalError")
		}
		if err := parserInterface.parseBooleanExpression(parserInterface); err != nil {
			return err
		}
		return parser.parseDelimiter()
	// CONST? terminalVariableType (LARRAY INT RARRAY)* ID (ASSIGN booleanExpression)? delimiter
	case parser.lookAHead(tokens.CONST), parser.lookAHead(tokens.BASIC), parser.lookAHead(tokens.TYPE):
		if parser.lookAHead(tokens.CONST) {
			if !parser.matchToken(tokens.CONST) {
				return parser.syntaxError("lexicalError")
			}
		}
		if err := parserInterface.parseTerminalVariableType(); err != nil {
			return err
		}
		for parser.lookAHead(tokens.LSBRACKET) {
			if !parser.matchToken(tokens.LSBRACKET) {
				return parser.syntaxError("lexicalError")
			}
			if !parser.matchToken(tokens.NUM) {
				return parser.syntaxError("Mismatched input '%v', expected <INT>")
			}
			if !parser.matchToken(tokens.RSBRACKET) {
				return parser.syntaxError("Mismatched input '%v', expected ']'")
			}
		}
		parser.lineBreakDelimiter = true
		if !parser.matchToken(tokens.ID) {
			return parser.syntaxError("Mismatched input '%v', expected <identifier> or '['")
		}
		if parser.lookAHead(tokens.ASSIGN) {
			if !parser.matchToken(tokens.ASSIGN) {
				return parser.syntaxError("lexicalError")
			}
			if err := parserInterface.parseBooleanExpression(parserInterface); err != nil {
				return err
			}
		}
		return parser.parseDelimiter()
	// ID ((LBRACKET ( booleanExpression (COMMA booleanExpression)* )? RBRACKET) | (arrayAccess* ASSIGN booleanExpression)) delimiter
	case parser.lookAHead(tokens.ID):
		if !parser.matchToken(tokens.ID) {
			return parser.syntaxError("lexicalError")
		}
		switch {
		// function call
		case parser.lookAHead(tokens.LBRACKET):
			if !parser.matchToken(tokens.LBRACKET) {
				return parser.syntaxError("lexicalError")
			}
			if !parser.lookAHead(tokens.RBRACKET) {
				if err := parserInterface.parseBooleanExpression(parserInterface); err != nil {
					return err
				}
				for parser.lookAHead(tokens.COMMA) {
					if !parser.matchToken(tokens.COMMA) {
						return parser.syntaxError("Mismatched input '%v', expected ',' or ')'")
					}
					if err := parserInterface.parseBooleanExpression(parserInterface); err != nil {
						return err
					}
				}
			}
			parser.lineBreakDelimiter = true
			if !parser.matchToken(tokens.RBRACKET) {
				return parser.syntaxError("Mismatched input '%v', expected ',' or ')'")
			}
		// array definiton
		case parser.lookAHead(tokens.LSBRACKET):
			for parser.lookAHead(tokens.LSBRACKET) {
				if err := parserInterface.parseArrayAccess(parserInterface); err != nil {
					return err
				}
			}
			fallthrough
		case parser.lookAHead(tokens.ASSIGN):
			if !parser.matchToken(tokens.ASSIGN) {
				return parser.syntaxError("Mismatched input '%v', expected '[', ',' or '='")
			}
			if err := parserInterface.parseBooleanExpression(parserInterface); err != nil {
				return err
			}
		default:
			_ = parser.matchToken(-1)
			return parser.syntaxError("Mismatched input '%v', expected '(', '[', ',' or '='")
		}
		return parser.parseDelimiter()
	default:
		return errors.New("StatementNotDefined")
	}
}

// conditionalScope
//   : booleanExpression scopeStatement
//   ;
func (parser *parser) parseConditionalScope(parserInterface Parser) error {
	if err := parserInterface.parseBooleanExpression(parserInterface); err != nil {
		return err
	}
	return parserInterface.parseScope(parserInterface)
}

// booleanExpression
//   :   comparisonExpression ((OR | AND) comparisonExpression)*
//   ;
func (parser *parser) parseBooleanExpression(parserInterface Parser) error {
	if err := parserInterface.parseComparisonExpression(parserInterface); err != nil {
		return err
	}
	var loopControl = false
	for {
		switch {
		case parser.lookAHead(tokens.OR):
			if !parser.matchToken(tokens.OR) {
				return parser.syntaxError("lexicalError")
			}
		case parser.lookAHead(tokens.BOOLOR):
			if !parser.matchToken(tokens.BOOLOR) {
				return parser.syntaxError("lexicalError")
			}
		case parser.lookAHead(tokens.AND):
			if !parser.matchToken(tokens.AND) {
				return parser.syntaxError("lexicalError")
			}
		case parser.lookAHead(tokens.BOOLAND):
			if !parser.matchToken(tokens.BOOLAND) {
				return parser.syntaxError("lexicalError")
			}
		default:
			loopControl = true
		}
		if loopControl {
			break
		}
		if err := parserInterface.parseComparisonExpression(parserInterface); err != nil {
			return err
		}
	}
	return nil
}

// comparisonExpression
// :   expression (comparisonOperator expression)*
// ;
func (parser *parser) parseComparisonExpression(parserInterface Parser) error {
	if err := parserInterface.parseExpression(parserInterface); err != nil {
		return err
	}
	var loopControl = false
	for {
		switch {
		case parser.lookAHead(tokens.EQ):
			if !parser.matchToken(tokens.EQ) {
				return parser.syntaxError("lexicalError")
			}
		case parser.lookAHead(tokens.NE):
			if !parser.matchToken(tokens.NE) {
				return parser.syntaxError("lexicalError")
			}
		case parser.lookAHead(tokens.GE):
			if !parser.matchToken(tokens.GE) {
				return parser.syntaxError("lexicalError")
			}
		case parser.lookAHead(tokens.GREATER):
			if !parser.matchToken(tokens.GREATER) {
				return parser.syntaxError("lexicalError")
			}
		case parser.lookAHead(tokens.LE):
			if !parser.matchToken(tokens.LE) {
				return parser.syntaxError("lexicalError")
			}
		case parser.lookAHead(tokens.LESS):
			if !parser.matchToken(tokens.LESS) {
				return parser.syntaxError("lexicalError")
			}
		default:
			loopControl = true
		}
		if loopControl {
			break
		}
		if err := parserInterface.parseExpression(parserInterface); err != nil {
			return err
		}
	}
	return nil
}

// arrayAccess
// : LARRAY expression RARRAY
// ;
func (parser *parser) parseArrayAccess(parserInterface Parser) error {
	if !parser.matchToken(tokens.LSBRACKET) {
		return parser.syntaxError("lexicalError")
	}
	if err := parserInterface.parseExpression(parserInterface); err != nil {
		return err
	}
	parser.lineBreakDelimiter = true
	if !parser.matchToken(tokens.RSBRACKET) {
		return parser.syntaxError("Mismatched input '%v', expected ']'")
	}
	return nil
}

// expression
// : term ((PLUS | MINUS) term)*
// ;
func (parser *parser) parseExpression(parserInterface Parser) error {
	if err := parserInterface.parseTerm(parserInterface); err != nil {
		return err
	}
	var loopControl = false
	for {
		switch {
		case parser.lookAHead(tokens.ADD):
			if !parser.matchToken(tokens.ADD) {
				return parser.syntaxError("lexicalError")
			}
		case parser.lookAHead(tokens.SUB):
			if !parser.matchToken(tokens.SUB) {
				return parser.syntaxError("lexicalError")
			}
		default:
			loopControl = true
		}
		if loopControl {
			break
		}
		if err := parserInterface.parseTerm(parserInterface); err != nil {
			return err
		}
	}
	return nil
}

// term
// : factor ((MULT | DIV) factor)*
// ;
func (parser *parser) parseTerm(parserInterface Parser) error {
	if err := parserInterface.parseFactor(parserInterface); err != nil {
		return err
	}
	var loopControl = false
	for loopControl {
		switch {
		case parser.lookAHead(tokens.MULT):
			if !parser.matchToken(tokens.MULT) {
				return parser.syntaxError("lexicalError")
			}
		case parser.lookAHead(tokens.DIV):
			if !parser.matchToken(tokens.DIV) {
				return parser.syntaxError("lexicalError")
			}
		default:
			loopControl = true
		}
		if loopControl {
			break
		}
		if err := parserInterface.parseFactor(parserInterface); err != nil {
			return err
		}
	}
	return nil
}

// factor
// : (MINUS | NOT)? unary
// ;
func (parser *parser) parseFactor(parserInterface Parser) error {
	switch {
	case parser.lookAHead(tokens.EXCLAMATION):
		if !parser.matchToken(tokens.EXCLAMATION) {
			return parser.syntaxError("lexicalError")
		}
	case parser.lookAHead(tokens.NOT):
		if !parser.matchToken(tokens.NOT) {
			return parser.syntaxError("lexicalError")
		}
	}
	return parserInterface.parseUnary(parserInterface)
}

// unary
// : (BASIC | TRUE | FALSE | LITERAL)
// | ID arrayAccess?
// | ID LBRACKET ( booleanExpression (COMMA booleanExpression)* )? RBRACKET   // func call
// | LBRACKET booleanExpression RBRACKET
// | LARRAY (expression (COMMA expression)* )? RARRAY           // set array value
// ;
func (parser *parser) parseUnary(parserInterface Parser) error {
	switch {
	// ID (arrayAccess?) | ID LBRACKET ( expression (COMMA expression)* )? RBRACKET
	case parser.lookAHead(tokens.ID):
		if !parser.matchToken(tokens.ID) {
			return parser.syntaxError("lexicalError")
		}
		// arrayAccess?
		if parser.lookAHead(tokens.LSBRACKET) {
			if err := parserInterface.parseArrayAccess(parserInterface); err != nil {
				return err
			}
			// LBRACKET ( booleanExpression (COMMA booleanExpression)* )? RBRACKET
		} else if parser.lookAHead(tokens.LBRACKET) {
			if !parser.matchToken(tokens.LBRACKET) {
				return parser.syntaxError("lexicalError")
			}
			if !parser.lookAHead(tokens.RBRACKET) {
				if err := parserInterface.parseBooleanExpression(parserInterface); err != nil {
					return err
				}
				for parser.lookAHead(tokens.COMMA) {
					if !parser.matchToken(tokens.COMMA) {
						return parser.syntaxError("lexicalError")
					}
					if err := parserInterface.parseBooleanExpression(parserInterface); err != nil {
						return err
					}
				}
			}
			parser.lineBreakDelimiter = true
			if !parser.matchToken(tokens.RBRACKET) {
				return parser.syntaxError("Mismatched input '%v', expected ',' or ')'")
			}
		}
	// LBRACKET booleanExpression RBRACKET
	case parser.lookAHead(tokens.LBRACKET):
		if !parser.matchToken(tokens.LBRACKET) {
			return parser.syntaxError("lexicalError")
		}
		if err := parserInterface.parseBooleanExpression(parserInterface); err != nil {
			return err
		}
		parser.lineBreakDelimiter = true
		if !parser.matchToken(tokens.RBRACKET) {
			return parser.syntaxError("Mismatched input '%v', expected ')'")
		}
	// LARRAY (expression (COMMA expression)* )? RARRAY
	case parser.lookAHead(tokens.LSBRACKET):
		if !parser.matchToken(tokens.LSBRACKET) {
			return parser.syntaxError("lexicalError")
		}
		if err := parserInterface.parseExpression(parserInterface); err != nil {
			return err
		}
		for parser.lookAHead(tokens.COMMA) {
			if !parser.matchToken(tokens.COMMA) {
				return parser.syntaxError("lexicalError")
			}
			if err := parserInterface.parseExpression(parserInterface); err != nil {
				return err
			}
		}
		parser.lineBreakDelimiter = true
		if !parser.matchToken(tokens.RSBRACKET) {
			return parser.syntaxError("Mismatched input '%v', expected ',' or ']'")
		}
	//
	default:
		if err := parserInterface.parseTerminal(); err != nil {
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
func (parser *parser) parseTerminal() error {
	parser.lineBreakDelimiter = true
	switch {
	case parser.lookAHead(tokens.NUM):
		if !parser.matchToken(tokens.NUM) {
			return parser.syntaxError("lexicalError")
		}
	case parser.lookAHead(tokens.REAL):
		if !parser.matchToken(tokens.REAL) {
			return parser.syntaxError("lexicalError")
		}
	case parser.lookAHead(tokens.TRUE):
		if !parser.matchToken(tokens.TRUE) {
			return parser.syntaxError("lexicalError")
		}
	case parser.lookAHead(tokens.FALSE):
		if !parser.matchToken(tokens.FALSE) {
			return parser.syntaxError("lexicalError")
		}
	case parser.lookAHead(tokens.LITERAL):
		if !parser.matchToken(tokens.LITERAL) {
			return parser.syntaxError("lexicalError")
		}
	default:
		_ = parser.matchToken(-1)
		return parser.syntaxError("Mismatched input '%v', expected <terminal>")
	}
	return nil
}
