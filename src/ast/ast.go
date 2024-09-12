package ast

import (
	"fmt"
	"simplelang/src/token"
	"strconv"
)

// This is what I originally intended for my languages to be capable of,
// however I kept getting more and more ideas :D

// declaration
// kind: let | const, name, expression, span

// expression
// kind: identifier | literal | call, span // for now there are no unary, binary and parenthesized expressions

// callexpression
// name, args: expression[], span

type Statement interface {
	isStatement()
}

type ValueDeclarationVariant int

const (
	ValueDeclarationVariant_const ValueDeclarationVariant = iota
	ValueDeclarationVariant_let
)

type ValueDeclaration struct {
	Variant    ValueDeclarationVariant
	Identifier string
	// TODO: should this just be a string, instead of a string pointer?
	ExplicitType *string
	Expression   *Expression
	token.Span
}

func (recv ValueDeclaration) isStatement() {}

type Parameter struct {
	Name string
	Type string
}

type FunctionDeclarationStatement struct {
	Identifier string
	Parameters []Parameter
	// for now without variadic params,
	// receiver
	// and generics
	ReturnTypes []string
	Statements  []Statement
}

func (recv FunctionDeclarationStatement) isStatement() {}

type ReturnStatement struct {
	Expression Expression
}

func (recv ReturnStatement) isStatement() {}

type LoopStatement struct {
	Statements []Statement
}

func (recv LoopStatement) isStatement() {}

type BreakStatement struct{}

func (recv BreakStatement) isStatement() {}

type Assignment struct {
	// for now only '='
	Identifier string
	Expression Expression
	token.Span
}

func (recv Assignment) isStatement() {}

type Expression interface {
	Statement
	isExpression()
}

type BlockExpression struct {
	Statements []Statement
	Expression *Expression
}

func (recv BlockExpression) isStatement()  {}
func (recv BlockExpression) isExpression() {}

type IfExpression struct {
	Condition  Expression
	Consequent BlockExpression
	Alternate  *Expression
}

func (recv IfExpression) isStatement()  {}
func (recv IfExpression) isExpression() {}

type ExpressionIdentifier struct {
	Identifier string
	token.Span
}

func (recv ExpressionIdentifier) isExpression() {}
func (recv ExpressionIdentifier) isStatement()  {}

type ExpressionLiteral struct {
	Literal Literal
	token.Span
}

func (recv ExpressionLiteral) isExpression() {}
func (recv ExpressionLiteral) isStatement()  {}

type ExpressionUnary struct {
	Operator   token.OperatorVariant
	Expression Expression
	token.Span
}

func (recv ExpressionUnary) isExpression() {}
func (recv ExpressionUnary) isStatement()  {}

type ExpressionBinary struct {
	Left     Expression
	Operator token.OperatorVariant
	Right    Expression
	token.Span
}

func (recv ExpressionBinary) isExpression() {}
func (recv ExpressionBinary) isStatement()  {}

type ExpressionParenthesized struct {
	Expression Expression
	token.Span
}

func (recv ExpressionParenthesized) isExpression() {}
func (recv ExpressionParenthesized) isStatement()  {}

type ExpressionCall struct {
	Identifier string
	Arguments  []Expression
	token.Span
}

func (recv ExpressionCall) isExpression() {}
func (recv ExpressionCall) isStatement()  {}

type Literal interface {
	isLiteral()
}

type StringLiteral struct {
	Value string
}

func (recv StringLiteral) isLiteral() {}

type IntLiteral struct {
	Value int64
}

func (recv IntLiteral) isLiteral() {}

type FloatLiteral struct {
	Value float64
}

func (recv FloatLiteral) isLiteral() {}

type InterpolatedStringLiteral struct {
	Value       string
	StringParts []string
	// TODO: this shouldn't only be an expression, but also a format string (:%s or even :%d02 for padding)
	Expressions []Expression
}

func (recv InterpolatedStringLiteral) isLiteral() {}

type Ast struct {
	tokens        []token.Token
	current_index int
	Statements    []Statement
}

func NewAst(tokens []token.Token) Ast {
	ast := Ast{
		tokens:     tokens,
		Statements: []Statement{},
	}
	ast.parse()
	return ast
}

func (recv *Ast) get_current_token() token.Token {
	return recv.tokens[recv.current_index]
}

func (recv *Ast) increment(by int) {
	recv.current_index += by
}

func (recv *Ast) parse() {
	recv.Statements = append(recv.Statements, recv.handle_body()...)
}
func (recv *Ast) handle_body() []Statement {
	statements := []Statement{}
for_label:
	for {
		if recv.current_index >= len(recv.tokens) {
			break
		}
		token_ := recv.tokens[recv.current_index]
		var statement Statement
		switch token_ := token_.(type) {
		case *token.Identifier:
			statement = recv.handle_identifier()
		case *token.Keyword:
			statement = recv.handle_keyword(token_)
		case *token.NewLine:
			// keep in mind that this also signals end of some expressions
			recv.increment(1)
			continue
		case *token.RightCurlyBrace:
			recv.increment(1)
			break for_label
		case *token.StringLiteral:
			statement = recv.handle_expression()
		default:
			panic(fmt.Sprintf("unexpected token.Token: %#v", token_))
		}
		statements = append(statements, statement)
	}
	return statements
}
func (recv *Ast) handle_block_expression() BlockExpression {
	// skip '{'
	recv.increment(1)
	statements := recv.handle_body()
	if len(statements) > 0 {
		if expr, last_is_expression := statements[len(statements)-1].(Expression); last_is_expression {
			statements = statements[:len(statements)-1]
			return BlockExpression{Statements: statements, Expression: &expr}
		}
	}
	return BlockExpression{Statements: statements, Expression: nil}
}

func (recv *Ast) handle_call_arguments() []Expression {
	arguments := []Expression{}
	// skip '('
	recv.increment(1)
	for {
		current_token := recv.get_current_token()
		if _, is_right_parenthesis := current_token.(*token.RightParenthesis); is_right_parenthesis {
			break
		}
		if _, is_comma := current_token.(*token.Comma); is_comma {
			continue
		}
		expression := recv.handle_expression()
		arguments = append(arguments, expression)
		current_token = recv.get_current_token()
		if _, is_right_parenthesis := current_token.(*token.RightParenthesis); is_right_parenthesis {
			break
		}
		recv.increment(1)
	}
	// skipping closing parenthesis
	recv.increment(1)
	return arguments
}

// handles identfiers like `a.b.c.d`
func (recv *Ast) handle_potentially_complex_identifier() string {
	must_be_identifier := recv.get_current_token()
	identifer, is_identifier := must_be_identifier.(*token.Identifier)
	if !is_identifier {
		panic("expected identifier")
	}
	recv.increment(1)
	might_be_dot := recv.get_current_token()
	_, is_dot := might_be_dot.(*token.Dot)
	if is_dot {
		recv.increment(1)
		return identifer.Name + "." + recv.handle_potentially_complex_identifier()
	}
	return identifer.Name
}

func (recv *Ast) handle_function_parameters() []Parameter {
	parameters := []Parameter{}

	must_be_left_parenthesis := recv.get_current_token()
	_, is_left_parenthesis := must_be_left_parenthesis.(*token.LeftParenthesis)
	if !is_left_parenthesis {
		fmt.Println(must_be_left_parenthesis)
		panic("expected left parenthesis")
	}

	// skip '('
	recv.increment(1)
	for {
		current_token := recv.get_current_token()
		if _, is_right_parenthesis := current_token.(*token.RightParenthesis); is_right_parenthesis {
			// skipping closing parenthesis
			recv.increment(1)
			break
		}
		param := Parameter{}
		// identifier for name
		must_be_identifier := recv.get_current_token()
		identifer, is_identifier := must_be_identifier.(*token.Identifier)
		if !is_identifier {
			panic("expected identifier")
		}
		param.Name = identifer.Name
		// identifier for type
		recv.increment(1)
		typeStr := ""
		might_be_operator := recv.get_current_token()
		operator, is_operator := might_be_operator.(*token.Operator)
		if is_operator {
			// technically not a multiply, but mistakes have been made :)
			if operator.OperatorVariant == token.OperatorVariant_Multiply {
				typeStr += "*"
			} else {
				panic(fmt.Sprintf("expected '*' or identifier but got operator: '%s'", operator.OperatorVariant))
			}
			recv.increment(1)
		}
		typeStr += recv.handle_potentially_complex_identifier()
		param.Type = typeStr
		// no need to increment, as "handle_potentially_complex_identifier()" has done it
		current_token = recv.get_current_token()
		if _, is_comma := current_token.(*token.Comma); is_comma {
			recv.increment(1)
		} else if _, is_parenthesis := current_token.(*token.RightParenthesis); is_parenthesis {
			parameters = append(parameters, param)
			recv.increment(1)
			break
		} else {
			panic("expected comma in function params")
		}
		parameters = append(parameters, param)
	}
	return parameters
}

func (recv *Ast) handle_function_return_types() []string {
	returnTypes := []string{}
	current_token := recv.get_current_token()
	if _, is_left_curly_brace := current_token.(*token.LeftCurlyBrace); is_left_curly_brace {
		return returnTypes
	}
	if _, is_left_parenthesis := current_token.(*token.LeftParenthesis); !is_left_parenthesis {
		identifier := recv.handle_potentially_complex_identifier()
		returnTypes = append(returnTypes, identifier)
		return returnTypes
	}
	// TODO: handle multiple return types
	return returnTypes
}

func (recv *Ast) handle_identifier() Statement {
	identifier := recv.handle_potentially_complex_identifier()

	current_token := recv.get_current_token()
	switch current_token := current_token.(type) {
	// dot is handled via handle_potentially_complex_identifer()
	case *token.EqualAssignment:
		recv.increment(1)
		expression := recv.handle_expression()
		return Assignment{Identifier: identifier, Expression: expression}
	case *token.LeftParenthesis:
		arguments := recv.handle_call_arguments()
		return ExpressionCall{Identifier: identifier, Arguments: arguments}
	case *token.NewLine:
		return ExpressionIdentifier{Identifier: identifier}
	default:
		panic(fmt.Sprintf("unexpected token.Token: %#v", current_token))
	}
}

func (recv *Ast) handle_identifier_expression() Expression {
	identifier := recv.handle_potentially_complex_identifier()
	current_token := recv.get_current_token()
	if _, ok := current_token.(*token.LeftParenthesis); ok {
		arguments := recv.handle_call_arguments()
		return ExpressionCall{Identifier: identifier, Arguments: arguments}
	}
	return ExpressionIdentifier{Identifier: identifier}
}

func (recv *Ast) handle_variable_declaration_explicit_type() string {
	// skipping colon
	recv.increment(1)
	must_be_identifier := recv.get_current_token()
	identifier, is_identifer := must_be_identifier.(*token.Identifier)
	if !is_identifer {
		panic("expected identifier")
	}

	recv.increment(1)

	return identifier.Name
}

func handle_interpolated_string_expression(runes []rune, i *int) Expression {
	// maybe the lexer could help here to parse more complex expressions?
	// for now only identifiers (although something like `i + 1`
	// is interpreted as an identifier and the builder just outputs
	// it, so free simple expressions :) )
	raw_expression := ""
	for {
		*i++
		current_rune := runes[*i]
		if current_rune == '}' {
			break
		}
		raw_expression += string(current_rune)
	}
	return ExpressionIdentifier{Identifier: raw_expression}
}

func string_to_interpolated_string(input string) InterpolatedStringLiteral {
	returnValue := InterpolatedStringLiteral{Value: input}
	runes := []rune(input)
	parts := []string{}
	expressions := []Expression{}
	current_part := ""
	for i := 0; i < len(runes); i++ {
		current_rune := runes[i]
		if current_rune == '{' {
			parts = append(parts, current_part)
			current_part = ""
			expressions = append(expressions, handle_interpolated_string_expression(runes, &i))
		} else {
			current_part += string(current_rune)
		}
	}
	parts = append(parts, current_part)
	returnValue.StringParts = parts
	returnValue.Expressions = expressions

	return returnValue
}

func (recv *Ast) handle_interpolated_string_expression() InterpolatedStringLiteral {
	// for now just simple strings "", not raw ``
	current_token := recv.get_current_token()
	switch current_token := current_token.(type) {
	case *token.StringLiteral:
		value := string_to_interpolated_string(current_token.Value)
		recv.increment(1)
		return value
	default:
		panic(fmt.Sprintf("unexpected token.Token: %#v", current_token))
	}
}

func (recv *Ast) handle_expression() Expression {
	current_token := recv.get_current_token()

	var left_expression Expression
	switch current_token := current_token.(type) {
	case *token.Dollar:
		recv.increment(1)
		left_expression = ExpressionLiteral{Literal: recv.handle_interpolated_string_expression()}
	case *token.Identifier:
		left_expression = recv.handle_identifier_expression()
	case *token.NumericLiteral:
		recv.increment(1)
		// this only allows base 10 ints
		if int_value, err := strconv.ParseInt(current_token.Value, 10, 64); err == nil {
			left_expression = ExpressionLiteral{Literal: IntLiteral{Value: int_value}}
			break
		}
		float_value, err := strconv.ParseFloat(current_token.Value, 64)
		if err != nil {
			panic(err)
		}
		left_expression = ExpressionLiteral{Literal: FloatLiteral{Value: float_value}}
	case *token.StringLiteral:
		recv.increment(1)
		left_expression = ExpressionLiteral{Literal: StringLiteral{Value: current_token.Value}}
	case *token.Operator:
		recv.increment(1)
		expression := recv.handle_expression()
		exprUnary := ExpressionUnary{Expression: expression}
		exprUnary.Operator = current_token.OperatorVariant
		left_expression = exprUnary
	case *token.LeftParenthesis:
		recv.increment(1)
		expression := recv.handle_expression()
		// assert that expression has been closed with right parenthesis
		if _, is_right_parenthesis := recv.get_current_token().(*token.RightParenthesis); !is_right_parenthesis {
			panic("missing right parenthesis")
		}
		recv.increment(1)
		left_expression = ExpressionParenthesized{Expression: expression}
	case *token.LeftCurlyBrace:
		left_expression = recv.handle_block_expression()
	case *token.Keyword:
		if current_token.KeywordVariant == token.KeywordVariant_If {
			left_expression = recv.handle_if_expression()
		} else {
			panic(fmt.Sprintf("unexpected keyword in expression: %s", current_token.KeywordVariant))
		}
	default:
		panic(fmt.Sprintf("unexpected token.Token: %#v", current_token))
	}

	if operator_token, is_operator := recv.get_current_token().(*token.Operator); is_operator {
		recv.increment(1)
		right := recv.handle_expression()
		right_binary_expression, is_binary_expression := right.(ExpressionBinary)
		if !is_binary_expression {
			return ExpressionBinary{Left: left_expression, Operator: operator_token.OperatorVariant, Right: right}
		}
		right_operator := right_binary_expression.Operator
		left_operator := operator_token.OperatorVariant
		left_first := left_operator.HasHigherPrecedenceThan(right_operator)
		if !left_first {
			return ExpressionBinary{Left: left_expression, Operator: operator_token.OperatorVariant, Right: right}
		}
		new_left_expression := ExpressionBinary{Left: left_expression, Operator: left_operator, Right: right_binary_expression.Left}
		new_right_expression := right_binary_expression.Right
		return ExpressionBinary{Left: new_left_expression, Operator: right_binary_expression.Operator, Right: new_right_expression}
	}

	return left_expression
}

type Import struct {
	Name string
	Path string
}
type ImportStatement struct {
	Imports []Import
}

func (recv ImportStatement) isStatement() {}

func (recv *Ast) handle_import_statement() ImportStatement {
	// this accepts only `import "math"` or `import m "math"`
	// but not
	// `import (
	// 		"fmt"
	//		m "math"
	// )`

	// skip import keyword
	recv.increment(1)
	import_ := Import{}
	current_token := recv.get_current_token()
	if identifier, is_identifier := current_token.(*token.Identifier); is_identifier {
		import_.Name = identifier.Name
		recv.increment(1)
	}
	importPath := recv.get_current_token().(*token.StringLiteral)
	import_.Path = importPath.Value
	recv.increment(1)
	return ImportStatement{Imports: []Import{import_}}
}

type PackageStatement struct {
	Name string
}

func (recv PackageStatement) isStatement() {}

func (recv *Ast) handle_package_statement() PackageStatement {
	statement := PackageStatement{}
	recv.increment(1)
	must_be_identifier := recv.get_current_token()
	identifier, is_identifier := must_be_identifier.(*token.Identifier)
	if !is_identifier {
		panic("expected identifier after package keyword")
	}
	statement.Name = identifier.Name
	recv.increment(1)
	if _, ok := recv.get_current_token().(*token.NewLine); !ok {
		panic("expected NewLine after package name")
	}
	recv.increment(1)
	return statement
}

func (recv *Ast) handle_value_variable_declaration(declaration_type ValueDeclarationVariant) ValueDeclaration {
	declaration := ValueDeclaration{Variant: declaration_type}

	// skipping declaration_type token
	recv.increment(1)
	must_be_identifier := recv.get_current_token()
	identifier, is_identifer := must_be_identifier.(*token.Identifier)
	if !is_identifer {
		panic("expected identifier")
	}
	declaration.Identifier = identifier.Name

	recv.increment(1)

	might_be_colon := recv.get_current_token()
	if _, ok := might_be_colon.(*token.Colon); ok {
		explicit_type := recv.handle_variable_declaration_explicit_type()
		// TODO: should type be just a string, instead of string pointer?
		declaration.ExplicitType = &explicit_type
	}

	// TODO: constants must be initialized
	might_be_equal_sign := recv.get_current_token()
	if _, ok := might_be_equal_sign.(*token.EqualAssignment); ok {
		recv.increment(1)
		init_expression := recv.handle_expression()
		declaration.Expression = &init_expression
	} else {
		if declaration.ExplicitType == nil {
			// go does not infer the type, if it is set later
			//
			// var test
			// test = "hi"
			//
			// it expects a type after "test"
			// and we would need a semantic analyzer to do that
			// so for now just panicking
			panic("if variable has not init, it must have an explicit type")
		}
	}

	return declaration
}

// should function declaration `fn test()` and function expression `let test = fn()`
// handled by different functions (that should share code of course)?
// I think yes, because `var test = func hi()` is not allowed
// `var test = func(recv *Abc)` is allowed, however Abc can't use it as method...
func (recv *Ast) handle_function_declaration() FunctionDeclarationStatement {
	declaration := FunctionDeclarationStatement{}

	// skipping func token
	recv.increment(1)

	// for now, no receiver

	must_be_identifier := recv.get_current_token()
	identifier, is_identifer := must_be_identifier.(*token.Identifier)
	if !is_identifer {
		panic("expected identifier")
	}
	declaration.Identifier = identifier.Name

	// skipping identifier
	recv.increment(1)
	parameters := recv.handle_function_parameters()
	declaration.Parameters = parameters

	// handle returns types
	returnTypes := recv.handle_function_return_types()
	declaration.ReturnTypes = returnTypes

	// TODO: expect '{’
	must_be_left_curly_brace := recv.get_current_token()
	_, is_left_curly_brace := must_be_left_curly_brace.(*token.LeftCurlyBrace)
	if !is_left_curly_brace {
		panic("expected left curly brace")
	}
	recv.increment(1)
	declaration.Statements = recv.handle_body()
	return declaration
}

func (recv *Ast) handle_return_statement() ReturnStatement {
	// skipping return keyword
	recv.increment(1)

	expression := recv.handle_expression()
	return ReturnStatement{Expression: expression}
}

func (recv *Ast) handle_loop_statement() LoopStatement {
	// skipping loop keyword and LeftCurlyBrace
	recv.increment(2)

	statements := recv.handle_body()
	return LoopStatement{Statements: statements}
}

func (recv *Ast) handle_break_statement() BreakStatement {
	// skipping break keyword
	recv.increment(1)

	return BreakStatement{}
}

func (recv *Ast) handle_if_expression() IfExpression {
	// skipping if keyword
	recv.increment(1)

	condition := recv.handle_expression()
	// we always enforce a block, I think this is necessary,
	// if we want to have if expressions, right?
	if _, is_right_curly_brace := recv.get_current_token().(*token.LeftCurlyBrace); !is_right_curly_brace {
		panic("if body needs to start with {")
	}
	block := recv.handle_block_expression()

	ifExpression := IfExpression{Condition: condition, Consequent: block}
	might_be_else := recv.get_current_token()
	next_keyword, next_is_keyword := might_be_else.(*token.Keyword)

	if next_is_keyword && next_keyword.KeywordVariant == token.KeywordVariant_Else {
		// skipping else keyword
		recv.increment(1)
		might_be_if := recv.get_current_token()
		next_keyword, next_is_keyword = might_be_if.(*token.Keyword)
		if next_is_keyword && next_keyword.KeywordVariant == token.KeywordVariant_If {
			// alternate := recv.handle_if_expession()
			alternate := recv.handle_expression()
			ifExpression.Alternate = &alternate
		} else {
			alternate := recv.handle_expression()
			ifExpression.Alternate = &alternate
		}
	}

	return ifExpression
}

func (recv *Ast) handle_keyword(keyword *token.Keyword) Statement {
	switch keyword.KeywordVariant {
	case token.KeywordVariant_Package:
		return recv.handle_package_statement()
	case token.KeywordVariant_Import:
		return recv.handle_import_statement()
	case token.KeywordVariant_Const:
		return recv.handle_value_variable_declaration(ValueDeclarationVariant_const)
	case token.KeywordVariant_Let:
		return recv.handle_value_variable_declaration(ValueDeclarationVariant_let)
	case token.KeywordVariant_Fn:
		// functions are expressions
		return recv.handle_function_declaration()
	case token.KeywordVariant_Return:
		return recv.handle_return_statement()
	case token.KeywordVariant_If:
		return recv.handle_if_expression()
	case token.KeywordVariant_Loop:
		return recv.handle_loop_statement()
	case token.KeywordVariant_Break:
		return recv.handle_break_statement()
	default:
		panic(fmt.Sprintf("unexpected token.KeywordVariant: %s", keyword.KeywordVariant))
	}

}
