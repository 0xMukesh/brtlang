package parser

import (
	"github.com/0xmukesh/interpreter/internal/ast"
	"github.com/0xmukesh/interpreter/internal/tokens"
)

func (p *Parser) parsePrintStmt() (*ast.AstNode, *ParserError) {
	node, err := p.Parse()
	if err != nil {
		return nil, err
	}

	if node != nil {
		expr, err := p.extractExpr(*node)
		if err != nil {
			return nil, err
		}

		err = NewParserError(MISSING_SEMICOLON, p.curr().Lexeme, p.curr().Line)
		p.consume(tokens.SEMICOLON, *err)
		return ast.NewAstNode(ast.STMT, ast.NewPrintStmt(expr, p.curr().Line)), nil
	} else {
		return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
	}
}

func (p *Parser) parseVarAssignStmt() (*ast.AstNode, *ParserError) {
	varNameNode, err := p.Parse()
	if err != nil {
		return nil, err
	}
	if varNameNode == nil {
		return nil, NewParserError(VARIABLE_NAME_EXPECTED, p.curr().Lexeme, p.curr().Line)
	}

	varNameExpr := varNameNode.ExtractExpr()
	if varNameExpr == nil {
		return nil, NewParserError(INVALID_VARIABLE_NAME, p.curr().Lexeme, p.curr().Line)
	}

	varNameLiteralExpr, isLiteralExpr := varNameExpr.(ast.LiteralExpr)
	if !isLiteralExpr {
		return nil, NewParserError(INVALID_VARIABLE_NAME, p.curr().Lexeme, p.curr().Line)
	}

	varName := varNameLiteralExpr.Value

	var varValueExpr ast.Expr

	if p.peek().Type == tokens.EQUAL {
		p.advance()

		varValueNode, err := p.Parse()
		if err != nil {
			return nil, err
		}

		if varValueNode == nil {
			return nil, NewParserError(EXPRESSION_AFTER_ASSIGNMENT_EXPECTED, p.curr().Lexeme, p.curr().Line)
		}

		varValueExpr = varValueNode.ExtractExpr()

		if varValueExpr == nil {
			return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
		}
	} else {
		varValueExpr = ast.NewLiteralExpr(tokens.NIL, "", p.curr().Line)
	}

	err = NewParserError(MISSING_SEMICOLON, p.curr().Lexeme, p.curr().Line)
	p.consume(tokens.SEMICOLON, *err)
	return ast.NewAstNode(ast.STMT, ast.NewVarAssignStmt(varName, varValueExpr, p.curr().Line)), nil
}

func (p *Parser) parseCreateBlockStmt() (*ast.AstNode, *ParserError) {
	rBraceFound := false
	var nodes []ast.AstNode

	for !p.isAtEnd() {
		nextTkn := p.peek()

		if nextTkn.Type == tokens.RIGHT_BRACE {
			p.advance()
			nodes = append(nodes, *ast.NewAstNode(ast.STMT, ast.NewCloseBlockStmt(p.curr().Line)))
			rBraceFound = true
			break
		}

		node, err := p.Parse()
		if err != nil {
			return nil, err
		}

		if node != nil {
			nodes = append(nodes, *node)
		}
	}

	if !rBraceFound {
		return nil, NewParserError(MISSING_RBRACE, p.curr().Lexeme, p.curr().Line)
	}

	return ast.NewAstNode(ast.STMT, ast.NewCreateBlockStmt(nodes, p.curr().Line)), nil
}

func (p *Parser) parseIfStmt() (*ast.AstNode, *ParserError) {
	ifConditionNode, err := p.Parse()
	if err != nil {
		return nil, err
	}

	if ifConditionNode == nil {
		return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
	}

	if !ifConditionNode.Value.IsExpr() {
		return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
	}

	ifBranch, err := p.Parse()
	if err != nil {
		return nil, err
	}

	if ifBranch == nil {
		return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
	}

	var elseIfStmts []ast.ElseIfStmt
	var elseStmt ast.ElseStmt

	for p.peek().Type == tokens.ELSE_IF {
		p.advance()

		elseIfConditionNode, err := p.Parse()
		if err != nil {
			return nil, err
		}

		if elseIfConditionNode == nil {
			return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
		}

		if !elseIfConditionNode.Value.IsExpr() {
			return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
		}

		elseIfBranch, err := p.Parse()
		if err != nil {
			return nil, err
		}

		if elseIfBranch == nil {
			return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
		}

		elseIfStmts = append(elseIfStmts, ast.NewElseIfStmt(elseIfConditionNode.ExtractExpr(), *elseIfBranch, p.curr().Line))
	}

	if p.peek().Type == tokens.ELSE {
		p.advance()

		elseBranch, err := p.Parse()
		if err != nil {
			return nil, err
		}

		if elseBranch == nil {
			return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
		}

		elseStmt = ast.NewElseStmt(*elseBranch, p.curr().Line)
	}

	return ast.NewAstNode(ast.STMT, ast.NewIfStmt(ifConditionNode.ExtractExpr(), *ifBranch, &elseIfStmts, &elseStmt, p.curr().Line)), nil
}

func (p *Parser) parseVarReassignStmt() (*ast.AstNode, *ParserError) {
	varName := p.curr().Lexeme
	// checking whether next token is "=" or not is handled within the switch-case statement
	p.advance()

	varValueNode, err := p.Parse()
	if err != nil {
		return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
	}

	if varValueNode == nil {
		return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
	}

	varValueExpr := varValueNode.ExtractExpr()
	if varValueExpr == nil {
		return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
	}

	return ast.NewAstNode(ast.STMT, ast.NewVarReassignStmt(varName, varValueExpr, p.curr().Line)), nil
}

func (p *Parser) parseWhileStmt() (*ast.AstNode, *ParserError) {
	expr, err := p.Parse()
	if err != nil || expr == nil {
		return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
	}

	node, err := p.Parse()
	if err != nil || node == nil {
		return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
	}

	return ast.NewAstNode(ast.STMT, ast.NewWhileStmt(expr.ExtractExpr(), *node, p.curr().Line)), nil
}
