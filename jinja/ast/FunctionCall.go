package ast

import (
	"fmt"
	"strings"

	"ddbt/compilerInterface"
	"ddbt/jinja/lexer"
)

type FunctionCall struct {
	position  lexer.Position
	name      string
	arguments funcCallArgs
}

type funcCallArg struct {
	name string
	arg  AST
}

type funcCallArgs []funcCallArg

var _ AST = &FunctionCall{}

func NewFunctionCall(token *lexer.Token, funcName string) *FunctionCall {
	return &FunctionCall{
		position:  token.Start,
		name:      funcName,
		arguments: make(funcCallArgs, 0),
	}
}

func (fc *FunctionCall) Position() lexer.Position {
	return fc.position
}

func (fc *FunctionCall) Execute(ec compilerInterface.ExecutionContext) (*compilerInterface.Value, error) {
	args, err := fc.arguments.Execute(ec)
	if err != nil {
		return nil, err
	}

	function := ec.GetVariable(fc.name)
	if function.IsUndefined {
		return nil, ec.ErrorAt(fc, fmt.Sprintf("function `%s` not found", fc.name))
	}

	if function.Type() != compilerInterface.FunctionalVal && function.Function == nil {
		return nil, ec.ErrorAt(fc, fmt.Sprintf("expected `%s` to be a function, got %s", fc.name, function.Type()))
	}

	if function.Function == nil {
		return nil, ec.ErrorAt(fc, "function is nil!")
	}

	result, err := function.Function(ec.PushState(), fc, args)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (fc *FunctionCall) String() string {
	var builder strings.Builder

	for i, arg := range fc.arguments {
		if i > 0 {
			builder.WriteString(", ")
		}

		if arg.name != "" {
			builder.WriteString(arg.name)
			builder.WriteRune('=')
		}

		builder.WriteString(arg.arg.String())
	}

	return fmt.Sprintf("{{ %s(%s) }}", fc.name, builder.String())
}

func (fc *FunctionCall) AddArgument(argName string, node AST) {
	fc.arguments = append(fc.arguments, funcCallArg{argName, node})
}

func (fca funcCallArgs) Execute(ec compilerInterface.ExecutionContext) (compilerInterface.Arguments, error) {
	arguments := make(compilerInterface.Arguments, 0, len(fca))

	for _, arg := range fca {
		result, err := arg.arg.Execute(ec)
		if err != nil {
			return nil, err
		}

		if result == nil {
			return nil, ec.NilResultFor(arg.arg)
		}

		arguments = append(arguments, compilerInterface.Argument{Name: arg.name, Value: result})
	}

	return arguments, nil
}
