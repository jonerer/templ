package parser

import (
	"github.com/a-h/parse"
)

// StripType takes the parser and throws away the return value.
func StripType[T any](p parse.Parser[T]) parse.Parser[any] {
	return parse.Func(func(in *parse.Input) (out any, ok bool, err error) {
		return p.Parse(in)
	})
}

func ExpressionOf(p parse.Parser[string]) parse.Parser[Expression] {
	return parse.Func(func(in *parse.Input) (out Expression, ok bool, err error) {
		from := in.Position()

		var exp string
		if exp, ok, err = p.Parse(in); err != nil || !ok {
			return
		}

		return NewExpression(exp, from, in.Position()), true, nil
	})
}

var lt = parse.Rune('<')
var gt = parse.Rune('>')
var openBrace = parse.String("{")
var optionalSpaces = parse.StringFrom(parse.Optional(
	parse.AtLeast(1, parse.Rune(' '))))
var openBraceWithPadding = parse.StringFrom(optionalSpaces,
	openBrace,
	optionalSpaces)
var openBraceWithOptionalPadding = parse.Any(openBraceWithPadding, openBrace)

var closeBrace = parse.String("}")
var closeBraceWithPadding = parse.String(" }")
var closeBraceWithOptionalPadding = parse.Any(closeBraceWithPadding, closeBrace)

var openBracket = parse.String("(")
var closeBracket = parse.String(")")
var closeBracketWithOptionalPadding = parse.StringFrom(optionalSpaces, closeBracket)

var stringUntilNewLine = parse.StringUntil[string](parse.NewLine)
var newLineOrEOF = parse.Or(parse.NewLine, parse.EOF[string]())
var stringUntilNewLineOrEOF = parse.StringUntil(newLineOrEOF)
