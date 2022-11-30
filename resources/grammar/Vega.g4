grammar Vega;

block
    :   (FUNC ID LBRACKET functionParameterDeclaration? RBRACKET functionReturnType scopeStatement)+ EOF
    ;

functionParameterDeclaration
	:   functionParameterDefinition (COMMA functionParameterDefinition)*
    ;

functionParameterDefinition
    :   (LARRAY RARRAY)* terminalVariableType ID
    ;

functionReturnType
    :   (LARRAY RARRAY)* terminalVariableType
    ;

scopeStatement
    :   LCURLY ((PASS DELIMITER) | (statement)+) RCURLY
    ;

statement
	:	CONST? (LARRAY INT RARRAY)* terminalVariableType ID (ASSIGN booleanExpression)? DELIMITER
	|   ID arrayAccess* ASSIGN booleanExpression DELIMITER
	|   ID funcCall
	|   RETURN booleanExpression DELIMITER
	|   CONTINUE DELIMITER
	|   BREAK DELIMITER
	|	WHILE conditionalScope
	|	IF conditionalScope (ELIF conditionalScope)* (ELSE scopeStatement)?
	|   SWITCH expression LCURLY (CASE terminal COLON statement+)+ (DEFAULT COLON statement+)? RCURLY
	;

booleanExpression
    :   comparisonExpression ((OR | AND) comparisonExpression)*
    ;

comparisonExpression
    :   expression (comparisonOperator expression)*
    ;

expression
    :   term ((PLUS | MINUS) term)*
    ;

term
    :	factor ((MULT | DIV) factor)*
	;

factor
    :   (MINUS  |NOT)? unary
    ;

unary
    :   terminal
    |   ID arrayAccess? // potential array access
    |   ID funcCall
    |   LBRACKET booleanExpression RBRACKET
    |   LARRAY (expression (COMMA expression)*)? RARRAY
    ;

funcCall
    :   LBRACKET ( booleanExpression (COMMA booleanExpression)*)? RBRACKET
    ;

arrayAccess
    :   LARRAY expression RARRAY
    ;

// Terminals
terminal
    :   INT
    |   FLOAT
    |   BOOL
    |   LITERAL
    |   CHAR
    ;
terminalVariableType
    :   INT_TYPE
    |   FLOAT_TYPE
    |   STRING_TYPE
    |   CHAR_TYPE
    |   BOOL_TYPE
    ;
comparisonOperator
    :   EQUAL
    |   NOTEQUAL
    |   GREATER
    |   GREATEREQ
    |   LESS
    |   LESSEQ
    ;


// Tokens
PLUS
    :   '+'
    ;
MINUS
    :   '-'
    ;
MULT
    :   '*'
    ;
DIV
    :   '/'
    ;
LBRACKET
    :   '('
    ;
RBRACKET
    :   ')'
    ;
LCURLY
    :   '{'
    ;
RCURLY
    :   '}'
    ;
LARRAY
    :   '['
    ;
RARRAY
    :   ']'
    ;
COLON
    :   ':'
    ;
COMMA
    :   ','
    ;
ASSIGN
    :   '='
    ;
LESS
    :   '<'
    ;
GREATER
    :   '>'
    ;
EQUAL
    :   '=='
    ;
LESSEQ
    :   '<='
    ;
GREATEREQ
    :   '>='
    ;
NOTEQUAL
    :   '!='
    ;
DELIMITER
    :   ';'
    |   '\n'
    ;

// Words
CONST
    :   'const'
    ;
FUNC
    :   'func'
    ;
WHILE
    :   'while'
    ;
IF
    :   'if'
    ;
ELIF
    :   'elif'
    ;
ELSE
    :   'else'
    ;
SWITCH
    :   'switch'
    ;
CASE
    :   'case'
    ;
DEFAULT
    :   'default'
    ;
RETURN
    :   'return'
    ;
PASS
    :   'pass'
    ;
CONTINUE
    :   'continue'
    ;
BREAK
    :   'break'
    ;
INT_TYPE
    :   'int'
    ;
FLOAT_TYPE
    :   'float'
    ;
STRING_TYPE
    :   'str'
    ;
CHAR_TYPE
    :   'char'
    ;
BOOL_TYPE
    :   'bool'
    ;
NOT
    :   'not'
    |   '!'
    ;
AND
    :   'and'
    |   '&&'
    ;
OR
    :   'or'
    |   '||'
    ;


ID  :	('a'..'z'|'A'..'Z'|'_') ('a'..'z'|'A'..'Z'|'0'..'9'|'_')*
    ;

BOOL
    :   'true'
    |   'false'
    ;

INT :	'0'..'9'+
    ;

FLOAT
    :   ('0'..'9')+ '.' ('0'..'9')* EXPONENT?
    |   '.' ('0'..'9')+ EXPONENT?
    |   ('0'..'9')+ EXPONENT
    ;

LITERAL
    :  '"' ( ESC_SEQ | ~('\\'|'"') )* '"'
    ;

CHAR:  '\'' ( ESC_SEQ | ~('\''|'\\') ) '\''
    ;

fragment
EXPONENT : ('e'|'E') ('+'|'-')? ('0'..'9')+ ;

fragment
HEX_DIGIT : ('0'..'9'|'a'..'f'|'A'..'F') ;

fragment
ESC_SEQ
    :   '\\' ('b'|'t'|'n'|'f'|'r'|'\\"'|'\''|'\\')
    |   UNICODE_ESC
    |   OCTAL_ESC
    ;

fragment
OCTAL_ESC
    :   '\\' ('0'..'3') ('0'..'7') ('0'..'7')
    |   '\\' ('0'..'7') ('0'..'7')
    |   '\\' ('0'..'7')
    ;

fragment
UNICODE_ESC
    :   '\\' 'u' HEX_DIGIT HEX_DIGIT HEX_DIGIT HEX_DIGIT
    ;

WS : [ \t\n]+ -> skip;
