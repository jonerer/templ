package parser

import (
	"strings"

	"github.com/a-h/parse"
)

var switchExpression parse.Parser[Node] = switchExpressionParser{}

type switchExpressionParser struct{}

func (_ switchExpressionParser) Parse(pi *parse.Input) (n Node, ok bool, err error) {
	// Check the prefix first.
	if _, ok, err = parse.String("switch ").Parse(pi); err != nil || !ok {
		return
	}

	// Once we have a prefix, we must have a Go expression.
	var r SwitchExpression
	if r.Expression, err = parseGoExpression("switch", pi); err != nil {
		return r, false, err
	}

	// Once we've had the start of a switch block, we must conclude the block.

	// Read the optional 'case' nodes.
	for {
		var ce CaseExpression
		ce, ok, err = caseExpressionParser.Parse(pi)
		if err != nil {
			return
		}
		if !ok {
			break
		}
		r.Cases = append(r.Cases, ce)
	}

	// Read the required closing brace.
	if _, ok, err = closeBraceWithOptionalPadding.Parse(pi); err != nil || !ok {
		err = parse.Error("switch: "+unterminatedMissingEnd, pi.Position())
		return
	}

	return r, true, nil
}

var caseExpressionStartParser = parse.Func(func(pi *parse.Input) (e Expression, ok bool, err error) {
	start := pi.Index()

	// Optional whitespace.
	if _, _, err = parse.OptionalWhitespace.Parse(pi); err != nil {
		return
	}

	// Look for a case or default statement.
	s, _ := pi.Peek(9)
	ok = (strings.HasPrefix(s, "case") || strings.HasPrefix(s, "default"))
	if !ok {
		pi.Seek(start)
		return
	}

	if e, err = parseGoExpression("case expression", pi); err != nil {
		return e, false, err
	}

	// Eat terminating newline.
	_, _, _ = parse.String("\n").Parse(pi)

	return
})

var caseExpressionParser = parse.Func(func(pi *parse.Input) (r CaseExpression, ok bool, err error) {
	if r.Expression, ok, err = caseExpressionStartParser.Parse(pi); err != nil || !ok {
		return
	}

	// Read until the next case statement, default, or end of the block.
	pr := newTemplateNodeParser(parse.Any(StripType(closeBraceWithOptionalPadding), StripType(caseExpressionStartParser)), "closing brace or case expression")
	var nodes Nodes
	if nodes, ok, err = pr.Parse(pi); err != nil || !ok {
		err = parse.Error("case: expected nodes, but none were found", pi.Position())
		return
	}
	r.Children = nodes.Nodes
	r.Diagnostics = nodes.Diagnostics

	// Optional whitespace.
	if _, ok, err = parse.OptionalWhitespace.Parse(pi); err != nil || !ok {
		return
	}

	return r, true, nil
})
