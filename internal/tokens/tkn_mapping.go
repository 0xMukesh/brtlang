package tokens

type TokenType int

const (
	EOF TokenType = iota
	ILLEGAL
	IGNORE // placeholder type for tokens which can be ignored

	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE

	COMMA
	SEMICOLON

	PLUS
	MINUS
	DOT
	STAR
	SLASH

	EQUAL
	EQUAL_EQUAL
	BANG
	BANG_EQUAL
	LESS
	LESS_EQUAL
	GREATER
	GREATER_EQUAL

	STRING
	NUMBER
	IDENTIFIER

	AND
	CLASS
	ELSE
	FALSE
	FOR
	FUN
	IF
	NIL
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE
)

var TknLiteralMapping = map[TokenType]string{
	LEFT_PAREN:    "(",
	RIGHT_PAREN:   ")",
	LEFT_BRACE:    "{",
	RIGHT_BRACE:   "}",
	COMMA:         ",",
	SEMICOLON:     ";",
	PLUS:          "+",
	MINUS:         "-",
	DOT:           ".",
	STAR:          "*",
	SLASH:         "/",
	EQUAL:         "=",
	EQUAL_EQUAL:   "==",
	BANG:          "!",
	BANG_EQUAL:    "!=",
	LESS:          "<",
	LESS_EQUAL:    "<=",
	GREATER:       ">",
	GREATER_EQUAL: ">=",
}

var ReservedKeywordsMapping = map[TokenType]string{
	AND:    AND.String(),
	CLASS:  CLASS.String(),
	ELSE:   ELSE.String(),
	FALSE:  FALSE.String(),
	FOR:    FOR.String(),
	FUN:    FUN.String(),
	IF:     IF.String(),
	NIL:    NIL.String(),
	PRINT:  PRINT.String(),
	RETURN: RETURN.String(),
	SUPER:  SUPER.String(),
	THIS:   THIS.String(),
	TRUE:   TRUE.String(),
	VAR:    VAR.String(),
	WHILE:  WHILE.String(),
}

func (t TokenType) IsReserved() bool {
	_, ok := ReservedKeywordsMapping[t]
	return ok
}

func (t TokenType) Literal() string {
	return TknLiteralMapping[t]
}

func (t TokenType) String() string {
	switch t {
	case EOF:
		return "EOF"
	case IGNORE:
		return "IGNORE"
	case LEFT_PAREN:
		return "LEFT_PAREN"
	case RIGHT_PAREN:
		return "RIGHT_PAREN"
	case LEFT_BRACE:
		return "LEFT_BRACE"
	case RIGHT_BRACE:
		return "RIGHT_BRACE"
	case COMMA:
		return "COMMA"
	case SEMICOLON:
		return "SEMICOLON"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case DOT:
		return "DOT"
	case STAR:
		return "STAR"
	case SLASH:
		return "SLASH"
	case EQUAL:
		return "EQUAL"
	case EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	case BANG:
		return "BANG"
	case BANG_EQUAL:
		return "BANG_EQUAL"
	case LESS:
		return "LESS"
	case LESS_EQUAL:
		return "LESS_EQUAL"
	case GREATER:
		return "GREATER"
	case GREATER_EQUAL:
		return "GREATER_EQUAL"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case IDENTIFIER:
		return "IDENTIFIER"
	case AND:
		return "AND"
	case CLASS:
		return "CLASS"
	case ELSE:
		return "ELSE"
	case FALSE:
		return "FALSE"
	case FOR:
		return "FOR"
	case FUN:
		return "FUN"
	case IF:
		return "IF"
	case NIL:
		return "NIL"
	case PRINT:
		return "PRINT"
	case RETURN:
		return "RETURN"
	case SUPER:
		return "SUPER"
	case THIS:
		return "THIS"
	case TRUE:
		return "TRUE"
	case VAR:
		return "VAR"
	case WHILE:
		return "WHILE"
	default:
		return "ILLEGAL"
	}
}
