package ast

import (
	"fmt"
	"strings"

	"ddbt/compilerInterface"
	"ddbt/jinja/lexer"
)

type List struct {
	position lexer.Position
	items    []AST
}

var _ AST = &List{}

func NewList(token *lexer.Token) *List {
	return &List{
		position: token.Start,
		items:    make([]AST, 0),
	}
}

func (l *List) Position() lexer.Position {
	return l.position
}

func (l *List) Execute(ec compilerInterface.ExecutionContext) (compilerInterface.AST, error) {
	return nil, nil
}

func (l *List) String() string {
	var builder strings.Builder

	for i, item := range l.items {
		if i > 0 {
			builder.WriteString(",\n\t\t")
		}

		builder.WriteString(item.String())
	}

	return fmt.Sprintf("[\n\t\t%s\n\t\t]", builder.String())
}

func (l *List) Append(item AST) {
	l.items = append(l.items, item)
}
