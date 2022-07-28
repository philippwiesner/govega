package frontend

import (
	"fmt"
	"govega/govega/language/tokens"
)

type Parser interface {
	getToken() error
	lookAHead(tag int) bool
	matchToken(tag int) error
	alreadyDeclared(identifier string) error
	parse(p Parser) error
	parseBlock(p Parser) error
	parseFunctionDeclaration(p Parser) error
	parseFunctionReturnType(p Parser, s *Symbol) error
	parseScope(p Parser, n string) error
}

type parserObject struct {
	errorState   ErrorState
	tokenStream  *TokenStream
	currentToken tokens.IToken
	table        *SymbolTable
}

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

func (pa *parserObject) getToken() error {
	token, err := pa.tokenStream.Remove()
	if err != nil {
		return NewSyntaxError(unexpectedEOF, pa.errorState)
	}
	pa.currentToken = token.GetToken()
	pa.errorState.lineNumber = token.GetTokenLine()
	return nil
}

func (pa *parserObject) lookAHead(tag int) bool {
	token := pa.tokenStream.GetHead()
	return token.GetTokenTag() == tag
}

func (pa *parserObject) matchToken(tag int) error {
	if err := pa.getToken(); err != nil {
		return err
	}
	if pa.currentToken.GetTag() != tag {
		return NewSyntaxError(invalidSyntax, pa.errorState)
	}
	return nil
}

func (pa *parserObject) alreadyDeclared(identifier string) error {
	_, ok := pa.table.Lookup(identifier)
	if ok {
		return NewDeclarationError(alreadyDefined, pa.errorState)
	}
	return nil
}

func (pa *parserObject) parse(p Parser) error {
	return p.parseBlock(p)
}

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
	if err := p.matchToken('('); err != nil {
		return err

	}
	if pa.lookAHead(tokens.ID) {
		if err := p.parseFunctionDeclaration(p); err != nil {
			return err
		}
	}
	if err := pa.matchToken(')'); err != nil {
		return fmt.Errorf("%v: Expected identifer declaration, but got: %v", err, pa.currentToken)
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
	if pa.lookAHead(tokens.FUNC) {
		return p.parseBlock(p)
	}
	return nil
}

func (pa *parserObject) parseFunctionDeclaration(p Parser) error {
	return nil
}

func (pa *parserObject) parseFunctionReturnType(p Parser, s *Symbol) error {
	pa.table.Add(s)
	return nil
}

func (pa *parserObject) parseScope(p Parser, name string) error {
	pa.table.NewScope(name)
	pa.table.LeaveScope()
	return nil
}
