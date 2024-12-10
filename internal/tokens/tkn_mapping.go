package tokens

type TokenType int

const (
	// 0
	EOF TokenType = iota
	ILLEGAL
	IGNORE // placeholder type for tokens which can be ignored

	// 3
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE

	// 7
	COMMA
	SEMICOLON

	// 9
	PLUS
	MINUS
	DOT
	STAR
	SLASH

	// 14
	EQUAL
	EQUAL_EQUAL
	BANG
	BANG_EQUAL
	LESS
	LESS_EQUAL
	GREATER
	GREATER_EQUAL

	// 22
	STRING
	NUMBER
	IDENTIFIER

	// 24
	TRUE
	FALSE
	NIL

	// 27
	IF
	ELSE_IF
	ELSE

	// 30
	VAR
	PRINT

	// 32
	AND
	OR

	// 34
	WHILE
	FOR

	// 36
	FUNC
	RETURN

	// 38
	CLASS
	SUPER
	THIS
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
	TRUE:    TRUE.String(),
	FALSE:   FALSE.String(),
	NIL:     NIL.String(),
	IF:      IF.String(),
	ELSE_IF: ELSE_IF.String(),
	ELSE:    ELSE.String(),
	VAR:     VAR.String(),
	PRINT:   PRINT.String(),
	AND:     AND.String(),
	OR:      OR.String(),
	WHILE:   WHILE.String(),
	FOR:     FOR.String(),
	FUNC:    FUNC.String(),
	RETURN:  RETURN.String(),
	CLASS:   CLASS.String(),
	SUPER:   SUPER.String(),
	THIS:    THIS.String(),
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
	case VAR:
		return "YO"
	case PRINT:
		return "YAP"
	case IF:
		return "HMM"
	case ELSE_IF:
		return "MID"
	case ELSE:
		return "NAH"
	case TRUE:
		return "BET"
	case FALSE:
		return "CAP"
	case NIL:
		return "NADA"
	case AND:
		return "N"
	case OR:
		return "EHH"
	case WHILE:
		return "VIBIN"
	case FOR:
		return "CHILLIN"
	case FUNC:
		return "VIBE"
	case RETURN:
		return "YOINK"
	case CLASS:
		return "CLASS"
	case SUPER:
		return "SUPER"
	case THIS:
		return "THIS"
	default:
		return "ILLEGAL"
	}
}
