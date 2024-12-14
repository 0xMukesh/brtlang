package runtime

import (
	"strconv"
	"time"

	"github.com/0xmukesh/interpreter/internal/ast"
	"github.com/0xmukesh/interpreter/internal/tokens"
)

var (
	VibeCheck      = "vibeCheck"
	PutTheFriesBro = "putTheFriesBro"
)

var NativeFnsReturnExprMapping = map[string]ast.AstNode{
	VibeCheck: NativeClockFnReturnExprNode(),
}

func NativeClockFnReturnExprNode() ast.AstNode {
	currentTimestamp := time.Now().Unix()
	return *ast.NewAstNode(ast.EXPR, ast.NewLiteralExpr(tokens.NUMBER, strconv.FormatInt(currentTimestamp, 10), -1))
}

func NativeClockFnMapping() FuncMapping {
	return FuncMapping{
		Node: NativeClockFnReturnExprNode(),
		Args: nil,
	}
}
