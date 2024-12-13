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
	MODULO

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

	TRUE
	FALSE
	NIL

	IF
	ELSE_IF
	ELSE

	VAR
	PRINT

	AND
	OR

	WHILE
	FOR

	FUNC
	RETURN
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
	MODULO:        "%",
	EQUAL:         "=",
	EQUAL_EQUAL:   "==",
	BANG:          "!",
	BANG_EQUAL:    "!=",
	LESS:          "<",
	LESS_EQUAL:    "<=",
	GREATER:       ">",
	GREATER_EQUAL: ">=",
	AND:           "&&",
	OR:            "||",
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
	WHILE:   WHILE.String(),
	FOR:     FOR.String(),
	FUNC:    FUNC.String(),
	RETURN:  RETURN.String(),
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
	case MODULO:
		return "MODULO"
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
		return "RIZZ"
	case PRINT:
		return "YAP"
	case IF:
		return "EDGING"
	case ELSE_IF:
		return "MID"
	case ELSE:
		return "AMOGUS"
	case TRUE:
		return "BET"
	case FALSE:
		return "CAP"
	case NIL:
		return "NADA"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case WHILE:
		return "VIBIN"
	case FOR:
		return "CHILLIN"
	case FUNC:
		return "SKIBIDI"
	case RETURN:
		return "BUSSIN"
	default:
		return "ILLEGAL"
	}
}
