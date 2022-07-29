package frontend

import (
	"fmt"
	"govega/govega/language"
	"govega/govega/language/tokens"
)

// Parser interface which allows better testing capacities
type Parser interface {
	getToken() error
	lookAHead(tag int) bool
	matchToken(tag int) error
	alreadyDeclared(identifier string) error
	parse(p Parser) error
	parseBlock(p Parser) error
	parseFunctionParamDeclaration(p Parser) error
	parseFunctionParamDefinition(p Parser) error
	parseVariableType(p Parser, s *Symbol) error
	parseFunctionReturnType(p Parser, s *Symbol) error
	parseTerminalVariableTypes(s *Symbol) (language.IBasicType, error)
	parseScope(p Parser, n string) error
	parseStatement(p Parser) error
	parseExpression(p Parser) error
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

// alreadyDeclared checks if an identifier for redeclaration
func (pa *parserObject) alreadyDeclared(identifier string) error {
	_, ok := pa.table.Lookup(identifier)
	if ok {
		return NewDeclarationError(alreadyDefined, pa.errorState)
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
//	(FUNC ID LBRACKET functionParamDeclaration? RBRACKET RETURN_TYPE functionReturnType scopeStatement)+ EOF
// ;
func (pa *parserObject) parseBlock(p Parser) error {
	if err := pa.matchToken(tokens.FUNC); err != nil {
		return fmt.Errorf("test: %v", err)
	}
	if err := pa.matchToken(tokens.ID); err != nil {
		return err
	}
	identifier := pa.currentToken.(tokens.IWord)
	if err := pa.alreadyDeclared(identifier.GetLexeme()); err != nil {
		return err
	}
	symbol := NewSymbol(identifier.GetLexeme(), nil, true, false)
	pa.table.Add(symbol)
	// add new scope for function param declaration
	pa.table.NewScope(identifier.GetLexeme())
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
	if err := pa.matchToken(tokens.RETURNTYPE); err != nil {
		return err
	}
	if err := p.parseFunctionReturnType(p, symbol); err != nil {
		return err
	}
	if err := p.parseScope(p, identifier.GetLexeme()); err != nil {
		return err
	}
	// leave function param scope
	pa.table.LeaveScope()
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
//   ID COLON variableTypes (ASSIGN expression)?
// ;
func (pa *parserObject) parseFunctionParamDefinition(p Parser) error {
	if err := p.matchToken(tokens.ID); err != nil {
		return err
	}
	identifier := pa.currentToken.(tokens.IWord)
	if err := pa.alreadyDeclared(identifier.GetLexeme()); err != nil {
		return err
	}
	symbol := NewSymbol(identifier.GetLexeme(), nil, true, false)
	if err := pa.matchToken(':'); err != nil {
		return err
	}
	if err := pa.parseVariableType(p, symbol); err != nil {
		return err
	}
	if pa.lookAHead('=') {
		if err := pa.matchToken('='); err != nil {
			return err
		}
		return p.parseExpression(p)
	}
	return nil
}

// parseVariableType parses identifier types
//
// variableTypes:
//   terminalVariableType (LARRAY INT RARRAY)*
// ;
func (pa *parserObject) parseVariableType(p Parser, s *Symbol) error {
	symbolType, err := p.parseTerminalVariableTypes(s)
	if err != nil {
		return err
	}
	for pa.lookAHead('[') {
		if err := pa.matchToken('['); err != nil {
			return err
		}
		if err := pa.matchToken(tokens.NUM); err != nil {
			return err
		}
		arraySize := pa.currentToken.(tokens.INum).GetValue()
		if err := pa.matchToken(']'); err != nil {
			return err
		}
		symbolType = language.NewArray(symbolType, arraySize)
	}
	s.SymbolType = symbolType
	pa.table.Add(s)
	return nil
}

// parseFunctionReturnType parse return type of function and sets the function identifier symbol type
//
// functionReturnType:
//   terminalVariableType (LARRAY RARRAY)*
// ;
func (pa *parserObject) parseFunctionReturnType(p Parser, s *Symbol) error {
	symbolType, err := p.parseTerminalVariableTypes(s)
	if err != nil {
		return err
	}
	for pa.lookAHead('[') {
		if err := pa.matchToken('['); err != nil {
			return err
		}
		if err := pa.matchToken(']'); err != nil {
			return err
		}
		symbolType = language.NewArray(symbolType, 0)
	}
	s.SymbolType = symbolType
	pa.table.Add(s)
	return nil
}

// parseTerminalVariableTypes parse basic variable type terminals
//
// terminalVariableType:
// | INT_TYPE
// | FLOAT_TYPE
// | CHAR_TYPE
// | BOOL_TYPE
// | STRING_TYPE
// ;
func (pa *parserObject) parseTerminalVariableTypes(s *Symbol) (language.IBasicType, error) {
	switch {
	case pa.lookAHead(tokens.BASIC):
		if err := pa.matchToken(tokens.BASIC); err != nil {
			return nil, err
		}
		return pa.currentToken.(language.IBasicType), nil
	case pa.lookAHead(tokens.TYPE):
		if err := pa.matchToken(tokens.TYPE); err != nil {
			return nil, err
		}
		switch pa.currentToken.(tokens.IWord).GetLexeme() {
		case "str":
			return pa.currentToken.(language.IStringType), nil
		}
	}
	return nil, NewSyntaxError("invalid type", pa.errorState)
}

// parseScope parses scopes
//
// scopeStatement:
//   LCURLY statement RCURLY
// ;
func (pa *parserObject) parseScope(p Parser, name string) error {
	pa.table.NewScope(name)
	if err := pa.matchToken('{'); err != nil {
		return err
	}
	if err := p.parseStatement(p); err != nil {
		return err
	}
	if err := pa.matchToken('}'); err != nil {
		return err
	}
	pa.table.LeaveScope()
	return nil
}

func (pa *parserObject) parseStatement(p Parser) error {
	return nil
}

func (pa *parserObject) parseExpression(p Parser) error {
	return nil
}
