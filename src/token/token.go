package token

import (
	"fmt"
	"unicode"
)

type Span struct {
	StartIndex       uint
	ExcludedEndIndex uint
	StartRowIndex    uint
	StartColumnIndex uint
	EndRowIndex      uint
	EndColumnIndex   uint
}

type Token interface {
	isToken()
	GetSpan() *Span
	String() string
}

type Identifier struct {
	Span
	Name string
}

func (w Identifier) isToken() {}
func (recv *Identifier) GetSpan() *Span {
	return &recv.Span
}
func (recv *Identifier) String() string {
	return fmt.Sprintf("{kind: Identifier, value: %s, span: %+v}", recv.Name, recv.Span)
}

type KeywordVariant int

const (
	KeywordVariant_Package KeywordVariant = iota
	KeywordVariant_Import
	KeywordVariant_Fn
	KeywordVariant_Const
	KeywordVariant_Let
	KeywordVariant_Return
	KeywordVariant_If
	KeywordVariant_Else
	KeywordVariant_Loop
	KeywordVariant_Break
)

var keywords = []string{
	"package",
	"import",
	"fn",
	"const",
	"let",
	"return",
	"if",
	"else",
	"loop",
	"break",
}

func (recv KeywordVariant) String() string {
	if recv < 0 || int(recv) >= len(keywords) {
		panic(fmt.Sprintf("unexpected token.KeywordVariant: %#v", recv))
	}
	return keywords[recv]
}

func isKeyword(word string) (KeywordVariant, bool) {
	for i, keyword := range keywords {
		if keyword == word {
			return KeywordVariant(i), true
		}
	}
	return -1, false
}

type Keyword struct {
	Span
	KeywordVariant
}

func (w Keyword) isToken() {}
func (recv *Keyword) GetSpan() *Span {
	return &recv.Span
}
func (recv *Keyword) String() string {
	return fmt.Sprintf("{kind: Keyword, value: %s, span: %+v}", recv.KeywordVariant, recv.Span)
}

type OperatorVariant int

const (
	OperatorVariant_Plus OperatorVariant = iota
	OperatorVariant_Minus
	OperatorVariant_Multiply
	OperatorVariant_Divide
	OperatorVariant_PowerOf
	OperatorVariant_Modulo
	// this is also used for creating pointers, so the name should
	// probably be something like 'SingleAnd'. Oh well, I'll keep that
	// in mind for my next language :D
	OperatorVariant_BinaryAnd
	OperatorVariant_BinaryOr
	OperatorVariant_LogicalAnd
	OperatorVariant_LogicalOr
	OperatorVariant_Not
	OperatorVariant_NotEquals
	OperatorVariant_Equals
	OperatorVariant_LowerThan
	OperatorVariant_LowerThanOrEqual
	OperatorVariant_GreaterThan
	OperatorVariant_GreaterThanOrEqual
)

func (recv OperatorVariant) HasHigherPrecedenceThan(other OperatorVariant) bool {
	precedences := map[OperatorVariant]int{
		OperatorVariant_Plus:        1,
		OperatorVariant_Minus:       1,
		OperatorVariant_Multiply:    2,
		OperatorVariant_Divide:      2,
		OperatorVariant_Modulo:      2,
		OperatorVariant_PowerOf:     3,
		OperatorVariant_LogicalAnd:  -1,
		OperatorVariant_GreaterThan: 0,
	}
	_, recv_ok := precedences[recv]
	_, other_ok := precedences[other]
	if !recv_ok || !other_ok {
		panic(fmt.Sprintf("precendence checking has not been implement for %s and %s", recv.String(), other.String()))
	}
	return precedences[recv] > precedences[other]
}

var operators = []string{
	"+",
	"-",
	"*",
	"/",
	"**",
	"%",
	"&",
	"|",
	"&&",
	"||",
	"!",
	"!=",
	"==",
	"<",
	"<=",
	">",
	">=",
}

func (recv OperatorVariant) String() string {
	if recv < 0 || int(recv) >= len(operators) {
		panic(fmt.Sprintf("unexpected token.OperatorVariant: %#v", recv))
	}
	return operators[recv]
}

type Operator struct {
	Span
	OperatorVariant
}

func (w Operator) isToken() {}
func (recv *Operator) GetSpan() *Span {
	return &recv.Span
}
func (recv *Operator) String() string {
	return fmt.Sprintf("{kind: Operator, value: %s, span: %+v}", recv.OperatorVariant, recv.Span)
}

type NumericLiteral struct {
	Span
	Value string
}

func (w NumericLiteral) isToken() {}
func (recv *NumericLiteral) GetSpan() *Span {
	return &recv.Span
}
func (recv *NumericLiteral) String() string {
	return fmt.Sprintf("{kind: NumericLiteral, value: %s, span: %+v}", recv.Value, recv.Span)
}

type StringLiteral struct {
	Span
	Value string
}

func (w StringLiteral) isToken() {}
func (recv *StringLiteral) GetSpan() *Span {
	return &recv.Span
}
func (recv *StringLiteral) String() string {
	return fmt.Sprintf("{kind: StringLiteral, value: %s, span: %+v}", recv.Value, recv.Span)
}

type EqualAssignment struct {
	Span
}

func (w EqualAssignment) isToken() {}
func (recv *EqualAssignment) GetSpan() *Span {
	return &recv.Span
}
func (recv *EqualAssignment) String() string {
	return fmt.Sprintf("{kind: EqualAssignment, span: %+v}", recv.Span)
}

type Colon struct {
	Span
}

func (w Colon) isToken() {}
func (recv *Colon) GetSpan() *Span {
	return &recv.Span
}
func (recv *Colon) String() string {
	return fmt.Sprintf("{kind: Colon, span: %+v}", recv.Span)
}

type Comma struct {
	Span
}

func (w Comma) isToken() {}
func (recv *Comma) GetSpan() *Span {
	return &recv.Span
}
func (recv *Comma) String() string {
	return fmt.Sprintf("{kind: Comma, span: %+v}", recv.Span)
}

type Dot struct {
	Span
}

func (w Dot) isToken() {}
func (recv *Dot) GetSpan() *Span {
	return &recv.Span
}
func (recv *Dot) String() string {
	return fmt.Sprintf("{kind: Dot, span: %+v}", recv.Span)
}

type Dollar struct {
	Span
}

func (w Dollar) isToken() {}
func (recv *Dollar) GetSpan() *Span {
	return &recv.Span
}
func (recv *Dollar) String() string {
	return fmt.Sprintf("{kind: Dollar, span: %+v}", recv.Span)
}

type LeftParenthesis struct {
	Span
}

func (w LeftParenthesis) isToken() {}
func (recv *LeftParenthesis) GetSpan() *Span {
	return &recv.Span
}
func (recv *LeftParenthesis) String() string {
	return fmt.Sprintf("{kind: LeftParenthesis, span: %+v}", recv.Span)
}

type RightParenthesis struct {
	Span
}

func (w RightParenthesis) isToken() {}
func (recv *RightParenthesis) GetSpan() *Span {
	return &recv.Span
}
func (recv *RightParenthesis) String() string {
	return fmt.Sprintf("{kind: RightParenthesis, span: %+v}", recv.Span)
}

type LeftCurlyBrace struct {
	Span
}

func (w LeftCurlyBrace) isToken() {}
func (recv *LeftCurlyBrace) GetSpan() *Span {
	return &recv.Span
}
func (recv *LeftCurlyBrace) String() string {
	return fmt.Sprintf("{kind: LeftCurlyBraces, span: %+v}", recv.Span)
}

type RightCurlyBrace struct {
	Span
}

func (w RightCurlyBrace) isToken() {}
func (recv *RightCurlyBrace) GetSpan() *Span {
	return &recv.Span
}
func (recv *RightCurlyBrace) String() string {
	return fmt.Sprintf("{kind: RightCurlyBraces, span: %+v}", recv.Span)
}

type NewLine struct {
	Span
}

func (w NewLine) isToken() {}
func (recv *NewLine) GetSpan() *Span {
	return &recv.Span
}
func (recv *NewLine) String() string {
	return fmt.Sprintf("{kind: NewLine, span: %+v}", recv.Span)
}

type UnexpectedCharacter struct {
	Span
	rune
}

func (w UnexpectedCharacter) isToken() {}
func (recv *UnexpectedCharacter) GetSpan() *Span {
	return &recv.Span
}
func (recv *UnexpectedCharacter) String() string {
	return fmt.Sprintf("{kind: UnexpectedCharacter, value: %c, span: %+v}", recv.rune, recv.Span)
}

type lexer struct {
	input                string
	current_index        uint
	current_char         rune
	tokens               []Token
	current_row_index    uint
	current_column_index uint
}

func lexerNew(input string) lexer {
	return lexer{
		input:  input,
		tokens: []Token{},
	}
}

func (recv *lexer) increment(by uint) {
	recv.current_index += by
	recv.current_column_index += by
}

func (recv *lexer) new_line() {
	recv.current_row_index += 1
	recv.current_column_index = 0
}

func (recv *lexer) get_nth_char(i uint) (rune, bool) {
	runes := []rune(recv.input)
	if i >= uint(len(runes)) {
		return ' ', false
	}
	return []rune(recv.input)[i], true
}

func (recv *lexer) lex() {
	for {
		if recv.current_index >= uint(len(recv.input)) {
			break
		}
		current_char := []rune(recv.input)[recv.current_index]
		recv.current_char = current_char
		var token Token
		switch current_char {
		case '=':
			token = recv.lex_equals()
		case ':':
			token = recv.lex_simple(&Colon{})
		case ',':
			token = recv.lex_simple(&Comma{})
		case '.':
			token = recv.lex_simple(&Dot{})
		case '$':
			token = recv.lex_simple(&Dollar{})
		case '(':
			token = recv.lex_simple(&LeftParenthesis{})
		case ')':
			token = recv.lex_simple(&RightParenthesis{})
		case '{':
			token = recv.lex_simple(&LeftCurlyBrace{})
		case '}':
			token = recv.lex_simple(&RightCurlyBrace{})
		case '"':
			token = recv.lex_string()
		case '+':
			token = recv.lex_simple(&Operator{OperatorVariant: OperatorVariant_Plus})
		case '-':
			token = recv.lex_simple(&Operator{OperatorVariant: OperatorVariant_Minus})
		case '*':
			token = recv.lex_multiply()
		case '/':
			token = recv.lex_simple(&Operator{OperatorVariant: OperatorVariant_Divide})
		case '%':
			token = recv.lex_simple(&Operator{OperatorVariant: OperatorVariant_Modulo})
		case '&':
			token = recv.lex_and()
		case '|':
			token = recv.lex_or()
		case '!':
			token = recv.lex_not()
		case '<':
			token = recv.lex_lower_than()
		case '>':
			token = recv.lex_greater_than()
		default:
			{
				if current_char >= '0' && current_char <= '9' || current_char == '-' {
					token = recv.lex_number()
				} else if unicode.IsSpace(current_char) {
					if current_char == '\n' {
						token = recv.lex_simple(&NewLine{})
						recv.new_line()
					} else {
						// all other whitespaces are skipped
						recv.increment(1)
						continue
					}
				} else if unicode.IsLetter(current_char) || current_char == '_' {
					token = recv.lex_word()
				} else {
					token = recv.lex_simple(&UnexpectedCharacter{rune: current_char})
				}
			}
		}
		recv.tokens = append(recv.tokens, token)
	}
}

func (recv *lexer) lex_simple(token Token) Token {
	span := token.GetSpan()
	span.StartIndex = recv.current_index
	span.StartRowIndex = recv.current_row_index
	span.StartColumnIndex = recv.current_column_index
	span.ExcludedEndIndex = recv.current_index + 1
	span.EndRowIndex = recv.current_row_index
	span.EndColumnIndex = recv.current_column_index + 1
	recv.increment(1)
	return token
}

func (recv *lexer) lex_multiply() Token {
	recv.increment(1)
	current_rune, ok := recv.get_nth_char(recv.current_index)
	if !ok {
		return &Operator{OperatorVariant: OperatorVariant_Multiply}
	}
	if current_rune == '*' {
		recv.increment(1)
		return &Operator{OperatorVariant: OperatorVariant_PowerOf}
	}
	return &Operator{OperatorVariant: OperatorVariant_Multiply}
}

func (recv *lexer) lex_and() Token {
	recv.increment(1)
	current_rune, ok := recv.get_nth_char(recv.current_index)
	if !ok {
		return &Operator{OperatorVariant: OperatorVariant_BinaryAnd}
	}
	if current_rune == '&' {
		recv.increment(1)
		return &Operator{OperatorVariant: OperatorVariant_LogicalAnd}
	}
	return &Operator{OperatorVariant: OperatorVariant_BinaryAnd}
}

func (recv *lexer) lex_or() Token {
	recv.increment(1)
	current_rune, ok := recv.get_nth_char(recv.current_index)
	if !ok {
		return &Operator{OperatorVariant: OperatorVariant_BinaryOr}
	}
	if current_rune == '|' {
		recv.increment(1)
		return &Operator{OperatorVariant: OperatorVariant_LogicalOr}
	}
	return &Operator{OperatorVariant: OperatorVariant_BinaryOr}
}

func (recv *lexer) lex_equals() Token {
	recv.increment(1)
	current_rune, ok := recv.get_nth_char(recv.current_index)
	if !ok {
		return &EqualAssignment{}
	}
	if current_rune == '=' {
		recv.increment(1)
		return &Operator{OperatorVariant: OperatorVariant_Equals}
	}
	return &EqualAssignment{}
}

func (recv *lexer) lex_not() Token {
	recv.increment(1)
	current_rune, ok := recv.get_nth_char(recv.current_index)
	if !ok {
		return &Operator{OperatorVariant: OperatorVariant_Not}
	}
	if current_rune == '=' {
		recv.increment(1)
		return &Operator{OperatorVariant: OperatorVariant_NotEquals}
	}
	return &Operator{OperatorVariant: OperatorVariant_Not}
}

func (recv *lexer) lex_lower_than() Token {
	recv.increment(1)
	current_rune, ok := recv.get_nth_char(recv.current_index)
	if !ok {
		return &Operator{OperatorVariant: OperatorVariant_LowerThan}
	}
	if current_rune == '=' {
		recv.increment(1)
		return &Operator{OperatorVariant: OperatorVariant_LowerThanOrEqual}
	}
	return &Operator{OperatorVariant: OperatorVariant_LowerThan}
}

func (recv *lexer) lex_greater_than() Token {
	recv.increment(1)
	current_rune, ok := recv.get_nth_char(recv.current_index)
	if !ok {
		return &Operator{OperatorVariant: OperatorVariant_GreaterThan}
	}
	if current_rune == '=' {
		recv.increment(1)
		return &Operator{OperatorVariant: OperatorVariant_GreaterThanOrEqual}
	}
	return &Operator{OperatorVariant: OperatorVariant_GreaterThan}
}

func (recv *lexer) lex_string() Token {
	string_literal := StringLiteral{}

	span := string_literal.GetSpan()
	span.StartIndex = recv.current_index
	span.StartRowIndex = recv.current_row_index
	span.StartColumnIndex = recv.current_column_index

	str := ""
	if recv.current_char != '"' {
		panic("called lex_string, even though not string")
	}
	recv.increment(1)
	// TODO: consider escaping
	for {
		var c rune
		c, ok := recv.get_nth_char(recv.current_index)
		if !ok {
			break
		}
		recv.increment(1)
		if c == '"' {
			break
		}
		str += string(c)
	}

	span.ExcludedEndIndex = recv.current_index
	span.EndRowIndex = recv.current_row_index
	span.EndColumnIndex = recv.current_column_index

	string_literal.Value = str
	return &string_literal
}

func (recv *lexer) lex_word() Token {
	span := Span{}
	span.StartIndex = recv.current_index
	span.StartRowIndex = recv.current_row_index
	span.StartColumnIndex = recv.current_column_index

	str := ""
	if !(recv.current_char == '_' || unicode.IsLetter(recv.current_char)) {
		panic("called lex_string, even though not string")
	}
	for {
		var c rune
		c, ok := recv.get_nth_char(recv.current_index)
		if !ok {
			break
		}
		if !(c == '_' || unicode.IsLetter(c) || unicode.IsNumber(c)) {
			break
		}
		recv.increment(1)
		str += string(c)
	}

	span.ExcludedEndIndex = recv.current_index
	span.EndRowIndex = recv.current_row_index
	span.EndColumnIndex = recv.current_column_index

	if keywordVariant, ok := isKeyword(str); ok {
		return &Keyword{KeywordVariant: keywordVariant, Span: span}
	}

	return &Identifier{Name: str, Span: span}
}

func (recv *lexer) lex_number() Token {
	span := Span{}
	span.StartIndex = recv.current_index
	span.StartRowIndex = recv.current_row_index
	span.StartColumnIndex = recv.current_column_index

	// TODO: consider allowing underscore notation 1_000_000

	str := ""
	// TODO: since there now is a '-' token, this is probably not
	// possible anymore -> remove
	if recv.current_char == '-' {
		str += "-"
		recv.increment(1)
	}
	// TODO: this only works if the checking for '-' has been removed
	// TODO: the check for octal should probably be done in ast
	// or rather its conversion
	started_with_zero := recv.current_char == '0'
	if started_with_zero {
		str += "0"
		recv.increment(1)
	}
	for {
		var c rune
		c, ok := recv.get_nth_char(recv.current_index)
		if !ok {
			break
		}

		// for now no e notation and only decimal numbers
		// e notation exponent might be negative, too
		// if c == 'e' || c == 'E' || c == '.' || c == '-' {
		if c == '.' {
			str += string(c)
			recv.increment(1)
			started_with_zero = false
			continue
		}
		if !(unicode.IsDigit(c)) {
			break
		}
		if started_with_zero {
			panic("octal values are currently not supported")
		}
		recv.increment(1)
		str += string(c)
	}

	span.ExcludedEndIndex = recv.current_index
	span.EndRowIndex = recv.current_row_index
	span.EndColumnIndex = recv.current_column_index

	return &NumericLiteral{Value: str, Span: span}
}

func Tokenize(input string) []Token {
	lexer := lexerNew(input)
	lexer.lex()
	return lexer.tokens
}
