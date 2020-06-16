package db

import (
	"awans.org/aft/internal/datatypes"
	"github.com/google/uuid"
)

type Code struct {
	ID                uuid.UUID
	Name              string
	Runtime           Runtime
	FunctionSignature FunctionSignature
	Code              string
	Function          func(interface{}) (interface{}, error)
	executor          CodeExecutor
}

type FunctionSignature int64

const (
	FromJSON FunctionSignature = iota
	RPC
)

type Runtime int64

const (
	Golang Runtime = iota
	Javascript
	Starlark
)

type CodeExecutor interface {
	Invoke(Code, interface{}) (interface{}, error)
}

type bootstrapCodeExecutor struct{}

func (*bootstrapCodeExecutor) Invoke(c Code, args interface{}) (interface{}, error) {
	fh := datatypes.GoFunctionHandle{Function: c.Function}
	return fh.Invoke(args)
}

var codeMap map[uuid.UUID]Code = map[uuid.UUID]Code{
	boolValidator.ID:   boolValidator,
	intValidator.ID:    intValidator,
	enumValidator.ID:   enumValidator,
	stringValidator.ID: stringValidator,
	textValidator.ID:   textValidator,
	uuidValidator.ID:   uuidValidator,
	floatValidator.ID:  floatValidator,
}
