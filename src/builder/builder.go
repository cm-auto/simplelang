package builder

import (
	"fmt"
	"simplelang/src/ast"
	"simplelang/src/token"
	"strings"
)

func (recv *Builder) add_import_module(name string, path string) {
	if recv.Imports == nil {
		recv.Imports = []Import{}
	}
	for _, import_ := range recv.Imports {
		if import_.Path == path && import_.Name == name {
			// already in here -> return
			return
		}
		if name != "" && import_.Name == name {
			panic("import redeclared")
		}
	}
	recv.Imports = append(recv.Imports, Import{Name: name, Path: path})
}

func (recv *Builder) handle_interpolated_string_literal(literal ast.InterpolatedStringLiteral) string {
	recv.add_import_module("", "fmt")
	str := "fmt.Sprintf("
	// TODO: currently this uses %v, as this is the easiest
	// should we use %d, %s, if we can infer the type?
	string_arg := `"` + strings.Join(literal.StringParts, "%v") + `"`
	str += string_arg
	for _, expression := range literal.Expressions {
		str += ", " + recv.handleExpression(expression)
	}
	str += ")"
	return str
}

func (recv *Builder) handleLiteral(literal ast.Literal) string {
	switch literal := literal.(type) {
	case ast.FloatLiteral:
		return fmt.Sprintf("%f", literal.Value)
	case ast.IntLiteral:
		return fmt.Sprintf("%d", literal.Value)
	case ast.InterpolatedStringLiteral:
		return recv.handle_interpolated_string_literal(literal)
	case ast.StringLiteral:
		return fmt.Sprintf(`"%s"`, literal.Value)
	default:
		panic(fmt.Sprintf("unexpected ast.Literal: %#v", literal))
	}
}

func (recv *Builder) handleExpressionCall(call ast.ExpressionCall) string {
	identifer := call.Identifier
	switch identifer {
	case "print":
		identifer = "fmt.Println"
		recv.add_import_module("", "fmt")
	case "printf":
		identifer = "fmt.Printf"
		recv.add_import_module("", "fmt")
	// this was a simple trick I used before there were binary expressions
	// I left it in, to show how simple one problem could be solved by using
	// a different way
	case "add":
		return recv.handleExpression(call.Arguments[0]) + " + " + recv.handleExpression(call.Arguments[1])
	}
	str := identifer + "("
	for i, expression := range call.Arguments {
		if i > 0 {
			str += ", "
		}
		str += recv.handleExpression(expression)
	}
	str += ")"
	return str
}
func (recv *Builder) handleExpressionLiteral(literal ast.ExpressionLiteral) string {
	return recv.handleLiteral(literal.Literal)
}

func (recv *Builder) handleExpressionUnary(expression ast.ExpressionUnary) string {
	str := ""
	// ** has to be handled differently
	str += expression.Operator.String()
	// str += "("
	str += recv.handleExpression(expression.Expression)
	// str += ")"
	return str
}
func (recv *Builder) handleExpressionBinary(expression ast.ExpressionBinary) string {
	str := ""
	// ** has to be handled differently
	if expression.Operator == token.OperatorVariant_PowerOf {
		recv.add_import_module("", "math")
		str += "math.Pow("
		str += recv.handleExpression(expression.Left)
		str += ", "
		str += recv.handleExpression(expression.Right)
		str += ")"
		return str
	}
	// str += "("
	str += recv.handleExpression(expression.Left)
	str += " " + expression.Operator.String() + " "
	str += recv.handleExpression(expression.Right)
	// str += ")"

	return str
}
func (recv *Builder) handleExpression(expression ast.Expression) string {
	str := ""
	switch expression := expression.(type) {
	case ast.ExpressionIdentifier:
		str = expression.Identifier
	case ast.ExpressionLiteral:
		str = recv.handleLiteral(expression.Literal)
	case ast.ExpressionCall:
		str = recv.handleExpressionCall(expression)
	case ast.ExpressionUnary:
		str = recv.handleExpressionUnary(expression)
	case ast.ExpressionBinary:
		str = recv.handleExpressionBinary(expression)
	case ast.ExpressionParenthesized:
		str = "("
		str += recv.handleExpression(expression.Expression)
		str += ")"
	case ast.BlockExpression:
		str = "{\n"
		str += recv.handleStatements(expression.Statements)
		if expression.Expression != nil {
			potentialIdentifier, hasIdentifier := recv.identifierStack.peek()
			if hasIdentifier {
				str += potentialIdentifier + " = "
			}
			str += recv.handleExpression(*expression.Expression)
		}
		str += "\n}"
	case ast.IfExpression:
		str += recv.handleIfExpression(expression)
	default:
		panic(fmt.Sprintf("unexpected ast.Expression: %#v", expression))
	}
	return str
}

func (recv *Builder) handleFunctionDeclarationStatement(declaration ast.FunctionDeclarationStatement) string {
	str := "func " + declaration.Identifier + "("
	for _, param := range declaration.Parameters {
		str += param.Name + " " + param.Type + ","
	}
	returnTypesStr := ""
	if len(declaration.ReturnTypes) > 1 {
		returnTypesStr += "("
	}
	for _, returnType := range declaration.ReturnTypes {
		returnTypesStr += returnType
	}
	if len(declaration.ReturnTypes) > 1 {
		returnTypesStr += ")"
	}
	str += ")" + returnTypesStr + "{\n"
	str += recv.handleStatements(declaration.Statements)
	str += "}"
	return str
}

func (recv *Builder) handleImportStatement(importStatment ast.ImportStatement) {
	for _, importStatement := range importStatment.Imports {
		recv.add_import_module(importStatement.Name, importStatement.Path)
	}
}

func (recv *Builder) handleReturnStatement(returnStatement ast.ReturnStatement) string {
	return "return " + recv.handleExpression(returnStatement.Expression)
}

func (recv *Builder) handleBlockAssignment(expression ast.Expression) string {
	str := "\n"
	str += recv.handleExpression(expression)
	return str
}
func isBlockExpression(expression ast.Expression) bool {
	switch expression.(type) {
	case ast.BlockExpression:
		return true
	case ast.IfExpression:
		return true
	default:
		return false
	}
}

func (recv *Builder) handleDeclaration(declaration ast.ValueDeclaration) string {
	str := ""
	switch declaration.Variant {
	case ast.ValueDeclarationVariant_const:
		str += "const"
	case ast.ValueDeclarationVariant_let:
		str += "var"
	}
	str += " " + declaration.Identifier
	if declaration.ExplicitType != nil {
		str += " " + *declaration.ExplicitType
	}
	if declaration.Expression != nil {
		if isBlockExpression(*declaration.Expression) {
			str += recv.handleBlockAssignment(*declaration.Expression)
		} else {
			str += " = " + recv.handleExpression(*declaration.Expression)
		}
	}
	return str
}
func (recv *Builder) handleAssignment(assignment ast.Assignment) string {
	str := assignment.Identifier
	if isBlockExpression(assignment.Expression) {
		str += recv.handleBlockAssignment(assignment.Expression)
	} else {
		str += " = " + recv.handleExpression(assignment.Expression)
	}
	return str
}

func (recv *Builder) handleIfExpression(ifExpression ast.IfExpression) string {
	str := "if "
	str += recv.handleExpression(ifExpression.Condition)
	str += "{\n"
	for _, statement := range ifExpression.Consequent.Statements {
		str += recv.handleStatement(statement)
	}
	if ifExpression.Consequent.Expression != nil {
		potentialIdentifier, hasIdentifier := recv.identifierStack.peek()
		if hasIdentifier {
			str += potentialIdentifier + " = "
		}
		str += recv.handleExpression(*ifExpression.Consequent.Expression)
	}
	str += "\n}"
	if ifExpression.Alternate != nil {
		str += " else "
		str += recv.handleExpression(*ifExpression.Alternate)
	}
	return str
}

func (recv *Builder) handleLoop(loopStatement ast.LoopStatement) string {
	str := "for {"
	str += recv.handleStatements(loopStatement.Statements)
	str += "}"
	return str
}

func (recv *Builder) handleBreak(ast.BreakStatement) string {
	return "break"
}

type Import struct {
	Name string
	Path string
}

type identifierStack struct {
	stack []string
}

func (recv *identifierStack) push(identifier string) {
	recv.stack = append(recv.stack, identifier)
}

func (recv *identifierStack) pop() (string, bool) {
	if len(recv.stack) == 0 {
		return "", false
	}
	last := recv.stack[len(recv.stack)-1]
	recv.stack = recv.stack[:len(recv.stack)-1]
	return last, true
}

func (recv *identifierStack) peek() (string, bool) {
	if len(recv.stack) == 0 {
		return "", false
	}
	last := recv.stack[len(recv.stack)-1]
	return last, true
}

type Builder struct {
	Package         string
	Imports         []Import
	identifierStack identifierStack
}

func (recv *Builder) handleStatement(statement ast.Statement) string {
	switch statement := statement.(type) {
	// TODO: should this be moved to BuildProgram?
	case ast.ImportStatement:
		recv.handleImportStatement(statement)
		return ""
	case ast.FunctionDeclarationStatement:
		panic("FunctionDeclarationStatement is only allowed in file body")
	case ast.ReturnStatement:
		return recv.handleReturnStatement(statement)
	case ast.ValueDeclaration:
		recv.identifierStack.push(statement.Identifier)
		defer recv.identifierStack.pop()
		return recv.handleDeclaration(statement)
	case ast.ExpressionCall:
		return recv.handleExpressionCall(statement)
	case ast.ExpressionIdentifier:
		// TODO: implement
		panic("not implemented")
	case ast.ExpressionLiteral:
		return recv.handleExpressionLiteral(statement)
	case ast.Assignment:
		recv.identifierStack.push(statement.Identifier)
		defer recv.identifierStack.pop()
		return recv.handleAssignment(statement)
	case ast.IfExpression:
		return recv.handleIfExpression(statement)
	case ast.LoopStatement:
		return recv.handleLoop(statement)
	case ast.BreakStatement:
		return recv.handleBreak(statement)
	default:
		panic(fmt.Sprintf("unexpected ast.Statement: %#v", statement))
	}
}

func (recv *Builder) handleStatements(statements []ast.Statement) string {
	body := ""
	for _, statement := range statements {
		body += recv.handleStatement(statement) + "\n"
	}
	return body
}

func BuildProgram(ast_ ast.Ast) string {

	builder := Builder{}

	mainBody := ""
	for _, statement := range ast_.Statements {
		switch statement := statement.(type) {
		case ast.PackageStatement:
			builder.Package = statement.Name
		case ast.FunctionDeclarationStatement:
			mainBody += builder.handleFunctionDeclarationStatement(statement)
		default:
			mainBody += builder.handleStatement(statement)
		}
		mainBody += "\n"
	}

	importStr := ""
	if len(builder.Imports) > 0 {
		importStr = "import ("
		for _, import_ := range builder.Imports {
			importStr += "\n\t"
			if import_.Name != "" {
				importStr += import_.Name + " "
			}
			importStr += `"` + import_.Path + `"`
		}
		importStr += "\n)"
	}

	if builder.Package == "" {
		panic("package name has not been supplied")
	}
	packageStr := "package " + builder.Package
	return fmt.Sprintf("%s\n\n%s\n\n%s\n", packageStr, importStr, mainBody)
}
