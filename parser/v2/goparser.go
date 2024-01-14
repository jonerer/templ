package parser

import (
	"fmt"
	"strings"

	"github.com/a-h/parse"
	"github.com/a-h/templ/parser/v2/goexpression"
)

func parseGoFuncDecl(pi *parse.Input) (r Expression, err error) {
	from := pi.Position()
	src, _ := pi.Peek(-1)
	src = strings.TrimPrefix(src, "templ ")
	expr, err := goexpression.ParseFunc("func " + src)
	if err != nil {
		return r, parse.Error(fmt.Sprintf("invalid template declaration: %v", err.Error()), pi.Position())
	}
	pi.Take(len(expr))
	to := pi.Position()
	return NewExpression(expr, from, to), nil
}

func parseGoExpression(name string, pi *parse.Input) (r Expression, err error) {
	from := pi.Position()
	src, _ := pi.Peek(-1)
	expr, err := goexpression.ParseExpression(src)
	if err != nil {
		return r, parse.Error(fmt.Sprintf("%s: invalid go expression: %v", name, err.Error()), pi.Position())
	}
	pi.Take(len(expr))
	to := pi.Position()
	return NewExpression(expr, from, to), nil
}

