package parser

import (
	"github.com/0xmukesh/interpreter/internal/ast"
	"github.com/0xmukesh/interpreter/internal/tokens"
	"github.com/0xmukesh/interpreter/internal/utils"
)

func (p *Parser) parsePrintStmt() (*ast.AstNode, *ParserError) {
	nodePtr, err := p.Parse()
	if err != nil {
		return nil, err
	}

	if nodePtr != nil {
		node := *nodePtr
		err = NewParserError(MISSING_SEMICOLON, p.curr().Lexeme, p.curr().Line)
		p.consume(tokens.SEMICOLON, *err)
		return ast.NewAstNode(ast.STMT, ast.NewPrintStmt(node, p.curr().Line)), nil
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

	var varValueNode *ast.AstNode

	if p.peek().Type == tokens.EQUAL {
		p.advance()

		varValueNode, err = p.Parse()
		if err != nil {
			return nil, err
		}

		if varValueNode == nil {
			return nil, NewParserError(EXPRESSION_AFTER_ASSIGNMENT_EXPECTED, p.curr().Lexeme, p.curr().Line)
		}
	} else {
		varValueNode = ast.NewAstNode(ast.EXPR, ast.NewLiteralExpr(tokens.NIL, "", p.curr().Line))
	}

	err = NewParserError(MISSING_SEMICOLON, p.curr().Lexeme, p.curr().Line)
	p.consume(tokens.SEMICOLON, *err)
	return ast.NewAstNode(ast.STMT, ast.NewVarAssignStmt(varName, *varValueNode, p.curr().Line)), nil
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

		elseIfStmts = append(elseIfStmts, ast.NewElseIfStmt(*elseIfConditionNode, *elseIfBranch, p.curr().Line))
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

	return ast.NewAstNode(ast.STMT, ast.NewIfStmt(*ifConditionNode, *ifBranch, &elseIfStmts, &elseStmt, p.curr().Line)), nil
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

	return ast.NewAstNode(ast.STMT, ast.NewVarReassignStmt(varName, *varValueNode, p.curr().Line)), nil
}

func (p *Parser) parseWhileStmt() (*ast.AstNode, *ParserError) {
	node, err := p.Parse()
	if err != nil || node == nil {
		return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
	}

	branch, err := p.Parse()
	if err != nil || branch == nil {
		return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
	}

	return ast.NewAstNode(ast.STMT, ast.NewWhileStmt(*node, *branch, p.curr().Line)), nil
}

func (p *Parser) parseFuncDeclarationStmt() (*ast.AstNode, *ParserError) {
	funcName, err := p.Parse()
	if err != nil || funcName == nil || funcName.ExtractExpr() == nil {
		return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
	}

	literalExpr, isLiteralExpr := funcName.ExtractExpr().(ast.LiteralExpr)

	if !isLiteralExpr {
		return nil, NewParserError(INVALID_EXPRESSION, p.curr().Lexeme, p.curr().Line)
	}

	if literalExpr.TokenType != tokens.IDENTIFIER {
		return nil, NewParserError(INVALID_EXPRESSION, p.curr().Lexeme, p.curr().Line)
	}

	p.consume(tokens.LEFT_PAREN, *NewParserError(MISSING_LPAREN, p.curr().Lexeme, p.curr().Line))

	var args []ast.AstNode

	if p.peek().Type != tokens.RIGHT_PAREN {
		// equivalent to do-while loop in java
		for ok := true; ok; ok = p.matchAndAdvance(tokens.COMMA) {
			node, err := p.Parse()
			if err != nil || node == nil {
				return nil, NewParserError(INVALID_EXPRESSION, p.curr().Lexeme, p.curr().Line)
			}

			args = append(args, *node)
		}
	}

	if len(args) >= 255 {
		return nil, NewParserError("can't have more than 255 arguments", p.curr().Lexeme, p.curr().Line)
	}

	p.consume(tokens.RIGHT_PAREN, *NewParserError(MISSING_RPAREN, p.curr().Lexeme, p.curr().Line))

	nodeTbe, err := p.Parse()
	if err != nil || nodeTbe == nil {
		return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
	}

	funcDeclarationStmt := ast.NewFuncDeclarationStmt(literalExpr.Value, args, *nodeTbe, p.curr().Line)

	currEnv := p.Runtime.CurrEnv()
	funcNode, _ := currEnv.GetFunc(funcDeclarationStmt.Name)
	if funcNode != nil {
		return nil, NewParserError(IDENTIFIER_ALREADY_EXISTS, p.curr().Lexeme, p.curr().Line)
	}

	currEnv.SetFunc(funcDeclarationStmt.Name, funcDeclarationStmt.Node, funcDeclarationStmt.Args)

	return ast.NewAstNode(ast.STMT, funcDeclarationStmt), nil
}

func (p *Parser) parseFuncCallStmt() (*ast.AstNode, *ParserError) {
	funcName := p.curr().Lexeme
	// checking whether next token is "(" or not is handled within the switch-case statement
	p.advance()

	var args []ast.AstNode

	if p.peek().Type != tokens.RIGHT_PAREN {
		// equivalent to do-while loop in java
		for ok := true; ok; ok = p.matchAndAdvance(tokens.COMMA) {
			node, err := p.Parse()

			if err != nil || node == nil {
				return nil, NewParserError(INVALID_EXPRESSION, p.curr().Lexeme, p.curr().Line)
			}

			args = append(args, *node)
		}
	}

	var returnNode ast.AstNode

	currEnvPtr := p.Runtime.CurrEnv()
	if currEnvPtr != nil {
		currEnv := *currEnvPtr

		funcNode, ok := currEnv.Funcs[funcName]
		if ok {
			switch v := funcNode.Node.Value.(type) {
			case ast.CreateBlockStmt:
				for _, node := range v.Nodes {
					_, ok := node.Value.(ast.ReturnStmt)
					if ok {
						returnNode = node
					}
				}
			case ast.ReturnStmt:
				returnNode = v.Node
			}
		}
	}

	if len(args) >= 255 {
		return nil, NewParserError("can't have more than 255 arguments", p.curr().Lexeme, p.curr().Line)
	}

	p.consume(tokens.RIGHT_PAREN, *NewParserError(MISSING_RPAREN, p.curr().Lexeme, p.curr().Line))

	return ast.NewAstNode(ast.STMT, ast.NewFuncCallStmt(funcName, args, returnNode, p.curr().Line)), nil
}

func (p *Parser) parseReturnStmt() (*ast.AstNode, *ParserError) {
	node, err := p.Parse()
	if err != nil {
		return nil, err
	}

	if node == nil {
		return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
	}

	return ast.NewAstNode(ast.STMT, ast.NewReturnStmt(*node, p.curr().Line)), nil
}

func (p *Parser) parseIncrementStmt() (*ast.AstNode, *ParserError) {
	varName := p.curr().Lexeme
	p.advance()

	if utils.IsReservedKeyword(varName) {
		return nil, NewParserError(INVALID_EXPRESSION, p.curr().Lexeme, p.curr().Line)
	}

	if varName == "" {
		return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
	}

	if !utils.IsAlphaOnly(varName) {
		return nil, NewParserError(INVALID_EXPRESSION, p.curr().Lexeme, p.curr().Line)
	}

	return ast.NewAstNode(ast.STMT, ast.NewIncrementStmt(varName, p.curr().Line)), nil
}

func (p *Parser) parseDecrementStmt() (*ast.AstNode, *ParserError) {
	varName := p.curr().Lexeme
	p.advance()

	if utils.IsReservedKeyword(varName) {
		return nil, NewParserError(INVALID_EXPRESSION, p.curr().Lexeme, p.curr().Line)
	}

	if varName == "" {
		return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
	}

	if !utils.IsAlphaOnly(varName) {
		return nil, NewParserError(INVALID_EXPRESSION, p.curr().Lexeme, p.curr().Line)
	}

	return ast.NewAstNode(ast.STMT, ast.NewDecrementStmt(varName, p.curr().Line)), nil
}
