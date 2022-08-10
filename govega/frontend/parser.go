package frontend

import (
	"errors"
	"fmt"
	"govega/govega/language/tokens"
)

// TODO: include typing system
// Parser interface which allows better testing capacities
type Parser interface {
	getToken() error
	lookAHead(tag int) bool
	matchToken(tag int) error
	parse(p Parser) error
	parseBlock(p Parser) error
	parseFunctionParamDeclaration(p Parser) error
	parseFunctionParamDefinition(p Parser) error
	parseFunctionReturnType(p Parser) error
	parseArrayDeclaration() error
	parseArrayAccess(p Parser) error
	parseTerminalVariableType() error
	parseScope(p Parser) error
	parseStatement(p Parser) error
	parseConditionalScope(p Parser) error
	parseExpression(p Parser) error
	parseTerm(p Parser) error
	parseFactor(p Parser) error
	parseUnary(p Parser) error
}

// parserObject stores needed objects to keep track during the parsing
type parserObject struct {
	errorState   ErrorState    // update information on the error state from the scanned tokens
	tokenStream  *TokenStream  // token stream generated from the lexer
	currentToken tokens.IToken // current token which is being analyzed
	table        *symbolTable  // symbolTable to store information about recognized identifiers
}

// NewParser generates a new Parser interface
func NewParser(code []byte, fileName string) (Parser, error) {
	lexer := NewLexer(code, fileName)
	tokenStream, err := lexer.Scan()
	if err != nil {
		return nil, err
	}
	var parser Parser = &parserObject{
		ErrorState{0, 0, fileName, ""},
		tokenStream,
		nil,
		NewSymbolTable(),
	}
	return parser, nil
}

// getToken updates the currentToken with the token on top of the tokenStream
func (pa *parserObject) getToken() error {
	token, err := pa.tokenStream.Remove()
	if err != nil {
		return NewSyntaxError(unexpectedEOF, pa.errorState)
	}
	pa.currentToken = token.GetToken()
	pa.errorState.lineNumber = token.GetTokenLine()
	return nil
}

// lookAHead compares a given tag with the current top of the tokenStream
func (pa *parserObject) lookAHead(tag int) bool {
	token := pa.tokenStream.GetHead()
	return token.GetTokenTag() == tag
}

// matchToken compares a given token with the currentToken
func (pa *parserObject) matchToken(tag int) error {
	if err := pa.getToken(); err != nil {
		return err
	}
	if pa.currentToken.GetTag() != tag {
		return NewSyntaxError(invalidSyntax, pa.errorState)
	}
	return nil
}

// parse starts parsing process. All functiones which are validating the grammar are using the Parser interface to make
// testing easier
func (pa *parserObject) parse(p Parser) error {
	return p.parseBlock(p)
}

// parseBlock parses block statements
//
// block:
//	(FUNC ID LBRACKET functionParamDeclaration? RBRACKET functionReturnType scopeStatement)+ EOF
// ;
func (pa *parserObject) parseBlock(p Parser) error {
	if err := pa.matchToken(tokens.FUNC); err != nil {
		return fmt.Errorf("test: %v", err)
	}
	if err := pa.matchToken(tokens.ID); err != nil {
		return err
	}
	if err := p.matchToken('('); err != nil {
		return err

	}
	if pa.lookAHead(tokens.ID) {
		if err := p.parseFunctionParamDeclaration(p); err != nil {
			return err
		}
	}
	if err := pa.matchToken(')'); err != nil {
		return fmt.Errorf("%v: Expected closing bracket: %v", err, pa.currentToken)
	}
	if err := p.parseFunctionReturnType(p); err != nil {
		return err
	}
	if err := p.parseScope(p); err != nil {
		return err
	}
	if pa.lookAHead(tokens.FUNC) {
		return p.parseBlock(p)
	}
	if err := pa.matchToken(tokens.EOF); err != nil {
		return fmt.Errorf("%v: Expected EOF, but got: %v", err, pa.currentToken)
	}
	return nil
}

// parseFunctionParamDeclaration parses function parameter list
//
// functionParameterDeclaration:
//	functionParameterDefinition (COMMA functionParameterDefinition)*
// ;
func (pa *parserObject) parseFunctionParamDeclaration(p Parser) error {
	if err := p.parseFunctionParamDefinition(p); err != nil {
		return err
	}
	if pa.lookAHead(',') {
		if err := pa.matchToken(','); err != nil {
			return err
		}
		return p.parseFunctionParamDeclaration(p)
	}
	return nil
}

// parseFunctionParamDefinition parse function parameter definition
//
// functionParameterDefinition:
//   terminalVariableType (arrayDeclaration)* ID
// ;
func (pa *parserObject) parseFunctionParamDefinition(p Parser) error {
	if err := p.parseTerminalVariableType(); err != nil {
		return err
	}
	for pa.lookAHead('[') {
		if err := p.parseArrayDeclaration(); err != nil {
			return err
		}
	}
	if err := pa.matchToken(tokens.ID); err != nil {
		return err
	}

	return nil
}

func (pa *parserObject) parseArrayDeclaration() error {
	if err := pa.matchToken('['); err != nil {
		return err
	}
	if err := pa.matchToken(tokens.NUM); err != nil {
		return err
	}
	if err := pa.matchToken(']'); err != nil {
		return err
	}
	return nil
}

// parseFunctionReturnType parse return type of function and sets the function identifier symbol type
//
// functionReturnType:
//   terminalVariableType (LARRAY RARRAY)*
func (pa *parserObject) parseFunctionReturnType(p Parser) error {
	if err := p.parseTerminalVariableType(); err != nil {
		return err
	}
	for pa.lookAHead('[') {
		if err := pa.matchToken('['); err != nil {
			return err
		}
		if err := pa.matchToken(']'); err != nil {
			return err
		}
	}
	return nil
}

// parseTerminalVariableType parse basic variable type terminals
//
// terminalVariableType:
// | INT_TYPE
// | FLOAT_TYPE
// | CHAR_TYPE
// | BOOL_TYPE
// | STRING_TYPE
// ;
func (pa *parserObject) parseTerminalVariableType() error {
	switch {
	case pa.lookAHead(tokens.BASIC):
		if err := pa.matchToken(tokens.BASIC); err != nil {
			return err
		}
		switch pa.currentToken.(tokens.IWord).GetLexeme() {
		case "int":
			return nil
		case "bool":
			return nil
		case "float":
			return nil
		case "char":
			return nil
		default:
			return errors.New("type not supported")
		}
	case pa.lookAHead(tokens.TYPE):
		if err := pa.matchToken(tokens.TYPE); err != nil {
			return err
		}
		switch pa.currentToken.(tokens.IWord).GetLexeme() {
		case "str":
			return nil
		default:
			return errors.New("type not supported")
		}
	default:
		return errors.New("type not supported")
	}
}

// parseScope parses scopes
//
// scopeStatement:
//   LCURLY statement RCURLY
// ;
func (pa *parserObject) parseScope(p Parser) error {
	if err := pa.matchToken('{'); err != nil {
		return err
	}
	if err := p.parseStatement(p); err != nil {
		return err
	}
	if err := pa.matchToken('}'); err != nil {
		return err
	}
	return nil
}

// parseStatement parses normal statements
//
// statement:
//	 (
//	       CONST? terminalVariableType arrayDeclaration* ID (COMMA ID)* (ASSIGN expression)? DELIMITER
//      |  ID arrayAccess* (COMMA ID arrayAccess* )* ASSIGN expression DELIMITER
//      |  RETURN expression DELIMITER
//      |  CONTINUE DELIMITER
//      |  BREAK DELIMITER
//      |  WHILE conditionalScope
//      |  IF conditionalScope (ELIF conditionalScope)* (ELSE scopeStatement)?
//      |  PASS DELIMITER
//   )+
// ;
func (pa *parserObject) parseStatement(p Parser) error {
	switch {
	// exit if scope ends
	case pa.lookAHead('}'):
		return nil
	// statement: PASS DELIMITER
	case pa.lookAHead(tokens.PASS):
		if err := pa.matchToken(tokens.PASS); err != nil {
			return err
		}
		if err := pa.matchToken(';'); err != nil {
			return err
		}
	// statement: CONTINUE DELIMITER
	case pa.lookAHead(tokens.CONTINUE):
		if err := pa.matchToken(tokens.CONTINUE); err != nil {
			return err
		}
		if err := pa.matchToken(';'); err != nil {
			return err
		}
	// statement: BREAK DELIMITER
	case pa.lookAHead(tokens.BREAK):
		if err := pa.matchToken(tokens.BREAK); err != nil {
			return err
		}
		if err := pa.matchToken(';'); err != nil {
			return err
		}
	// statement: IF conditionalScope (ELIF conditionalScope)* (ELSE scopeStatement)?
	case pa.lookAHead(tokens.IF):
		if err := pa.matchToken(tokens.IF); err != nil {
			return err
		}
		if err := p.parseConditionalScope(p); err != nil {
			return err
		}
		for pa.lookAHead(tokens.ELIF) {
			if err := pa.matchToken(tokens.ELIF); err != nil {
				return err
			}
			if err := p.parseConditionalScope(p); err != nil {
				return err
			}
		}
		if pa.lookAHead(tokens.ELSE) {
			if err := p.parseScope(p); err != nil {
				return err
			}
		}
	// statement: WHILE conditionalScope
	case pa.lookAHead(tokens.WHILE):
		if err := pa.matchToken(tokens.WHILE); err != nil {
			return err
		}
		if err := p.parseConditionalScope(p); err != nil {
			return err
		}
	// statement: PASS DELIMITER
	case pa.lookAHead(tokens.PASS):
		if err := pa.matchToken(tokens.PASS); err != nil {
			return err
		}
		if err := pa.matchToken(';'); err != nil {
			return err
		}
	// statement: RETURN expresion DELIMITER
	case pa.lookAHead(tokens.RETURN):
		if err := pa.matchToken(tokens.RETURN); err != nil {
			return err
		}
		if err := p.parseExpression(p); err != nil {
			return err
		}
		if err := pa.matchToken(';'); err != nil {
			return err
		}
	// CONST? terminalVariableType arrayDeclaration* ID (COMMA ID)* (ASSIGN expression)? DELIMITER
	case pa.lookAHead(tokens.CONST), pa.lookAHead(tokens.BASIC), pa.lookAHead(tokens.TYPE):
		if pa.lookAHead(tokens.CONST) {
			if err := pa.matchToken(tokens.CONST); err != nil {
				return err
			}
		}
		if err := p.parseTerminalVariableType(); err != nil {
			return err
		}
		for pa.lookAHead('[') {
			if err := p.parseArrayDeclaration(); err != nil {
				return err
			}
		}
		if err := pa.matchToken(tokens.ID); err != nil {
			return err
		}
		for pa.lookAHead(',') {
			if err := pa.matchToken(tokens.ID); err != nil {
				return err
			}
		}
		if pa.lookAHead('=') {
			if err := pa.matchToken('='); err != nil {
				return err
			}
			if err := p.parseExpression(p); err != nil {
				return err
			}
		}
		if err := pa.matchToken(';'); err != nil {
			return err
		}
	// ID arrayAccess* (COMMA ID arrayAccess* )* ASSIGN expression DELIMITER
	case pa.lookAHead(tokens.ID):
		if err := pa.matchToken(tokens.ID); err != nil {
			return err
		}
		for pa.lookAHead('[') {
			if err := p.parseArrayAccess(p); err != nil {
				return err
			}
		}
		for pa.lookAHead(',') {
			if err := pa.matchToken(tokens.ID); err != nil {
				return err
			}
			if err := p.parseArrayAccess(p); err != nil {
				return err
			}
		}
		if err := pa.matchToken('='); err != nil {
			return err
		}
		if err := p.parseExpression(p); err != nil {
			return err
		}
		if err := pa.matchToken(';'); err != nil {
			return err
		}
	default:
		return errors.New("invalid statement")
	}
	// parse further statements
	if err := p.parseStatement(p); err != nil {
		return err
	}
	return nil
}

// conditionalScope:
//   LBRACKET expression RBRACKET scopeStatement
// ;
func (pa *parserObject) parseConditionalScope(p Parser) error {
	if err := pa.matchToken('('); err != nil {
		return err
	}
	if err := p.parseExpression(p); err != nil {
		return err
	}
	if err := pa.matchToken(')'); err != nil {
		return err
	}
	return nil
}

// arrayAccess:
//   LARRAY expression RARRAY
// ;
func (pa *parserObject) parseArrayAccess(p Parser) error {
	if err := pa.matchToken('['); err != nil {
		return err
	}
	if err := p.parseExpression(p); err != nil {
		return err
	}
	if err := pa.matchToken('['); err != nil {
		return err
	}
	return nil
}

// expression:
//   term (PLUS term | MINUS term | OR term)*
// ;
func (pa *parserObject) parseExpression(p Parser) error {
	if err := p.parseTerm(p); err != nil {
		return err
	}
	var loopControl bool = true
	for loopControl {
		switch {
		case pa.lookAHead('+'):
			if err := pa.matchToken('+'); err != nil {
				return nil
			}
		case pa.lookAHead('-'):
			if err := pa.matchToken('-'); err != nil {
				return nil
			}
		case pa.lookAHead(tokens.OR):
			if err := pa.matchToken(tokens.OR); err != nil {
				return nil
			}
		case pa.lookAHead(tokens.BOOLOR):
			if err := pa.matchToken(tokens.BOOLOR); err != nil {
				return nil
			}
		default:
			loopControl = false
		}
	}
	return nil
}

// term:
//   factor (MULT factor | DIV factor | AND factor)*
// ;
func (pa *parserObject) parseTerm(p Parser) error {
	if err := p.parseFactor(p); err != nil {
		return err
	}
	var loopControl bool = true
	for loopControl {
		switch {
		case pa.lookAHead('*'):
			if err := pa.matchToken('*'); err != nil {
				return nil
			}
		case pa.lookAHead('/'):
			if err := pa.matchToken('/'); err != nil {
				return nil
			}
		case pa.lookAHead(tokens.AND):
			if err := pa.matchToken(tokens.AND); err != nil {
				return nil
			}
		case pa.lookAHead(tokens.BOOLAND):
			if err := pa.matchToken(tokens.BOOLAND); err != nil {
				return nil
			}
		default:
			loopControl = false
		}
	}
	return nil
}

// factor:
//   (NOT|MINUS)? unary (comparisonOperator unary)*
// ;
func (pa *parserObject) parseFactor(p Parser) error {
	switch {
	case pa.lookAHead('!'):
		if err := pa.matchToken('!'); err != nil {
			return err
		}
	case pa.lookAHead(tokens.NOT):
		if err := pa.matchToken(tokens.NOT); err != nil {
			return err
		}
	case pa.lookAHead('-'):
		if err := pa.matchToken('-'); err != nil {
			return err
		}
	}
	if err := p.parseUnary(p); err != nil {
		return err
	}
	var loopControl bool = true
	for loopControl {
		switch {
		case pa.lookAHead(tokens.EQ):
			if err := pa.matchToken(tokens.EQ); err != nil {
				return nil
			}
		case pa.lookAHead(tokens.NE):
			if err := pa.matchToken(tokens.NE); err != nil {
				return nil
			}
		case pa.lookAHead(tokens.GE):
			if err := pa.matchToken(tokens.GE); err != nil {
				return nil
			}
		case pa.lookAHead('>'):
			if err := pa.matchToken('>'); err != nil {
				return nil
			}
		case pa.lookAHead(tokens.LE):
			if err := pa.matchToken(tokens.LE); err != nil {
				return nil
			}
		case pa.lookAHead('<'):
			if err := pa.matchToken('<'); err != nil {
				return nil
			}
		default:
			loopControl = false
		}
	}
	return nil
}

// unary:
//    (BASIC | TRUE | FALSE | LITERAL)
//  | ID arrayAccess?
//  | ID LBRACKET ( expression (COMMA expression)* )? RBRACKET   // func call
//  | LBRACKET expression RBRACKET
//  | LARRAY (expression (COMMA expression)* )? RARRAY           // set array value
// ;
func (pa *parserObject) parseUnary(p Parser) error {
	switch {
	// ID (arrayAccess?) | ID LBRACKET ( expression (COMMA expression)* )? RBRACKET
	case pa.lookAHead(tokens.ID):
		if err := p.matchToken(tokens.ID); err != nil {
			return err
		}
		// arrayAccess?
		if pa.lookAHead('[') {
			if err := p.parseArrayAccess(p); err != nil {
				return err
			}
			// LBRACKET ( expression (COMMA expression)* )? RBRACKET
		} else if pa.lookAHead('(') {
			if err := p.matchToken('('); err != nil {
				return nil
			}
			if err := p.parseExpression(p); err != nil {
				return nil
			}
			for pa.lookAHead(',') {
				if err := p.matchToken(','); err != nil {
					return nil
				}
				if err := p.parseExpression(p); err != nil {
					return nil
				}
			}
			if err := p.matchToken(')'); err != nil {
				return nil
			}
		}
	// LBRACKET expression RBRACKET
	case pa.lookAHead('('):
		if err := p.matchToken('('); err != nil {
			return nil
		}
		if err := p.parseExpression(p); err != nil {
			return nil
		}
		if err := p.matchToken(')'); err != nil {
			return nil
		}
	// LARRAY (expression (COMMA expression)* )? RARRAY
	case pa.lookAHead('['):
		if err := p.matchToken('['); err != nil {
			return nil
		}
		if err := p.parseExpression(p); err != nil {
			return nil
		}
		for pa.lookAHead(',') {
			if err := p.matchToken(','); err != nil {
				return nil
			}
			if err := p.parseExpression(p); err != nil {
				return nil
			}
		}
		if err := p.matchToken(']'); err != nil {
			return nil
		}
	// NUM | FLOAT | TRUE | FALSE | LITERAL
	default:
		switch {
		case pa.lookAHead(tokens.NUM):
			if err := pa.matchToken(tokens.NUM); err != nil {
				return err
			}
		case pa.lookAHead(tokens.REAL):
			if err := pa.matchToken(tokens.REAL); err != nil {
				return err
			}
		case pa.lookAHead(tokens.TRUE):
			if err := pa.matchToken(tokens.TRUE); err != nil {
				return err
			}
		case pa.lookAHead(tokens.FALSE):
			if err := pa.matchToken(tokens.FALSE); err != nil {
				return err
			}
		case pa.lookAHead(tokens.LITERAL):
			if err := pa.matchToken(tokens.LITERAL); err != nil {
				return err
			}
		}
	}
	return nil
}
