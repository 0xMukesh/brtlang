package runtime

import (
	"fmt"
)

type RuntimeValue struct {
	Value interface{}
}

func NewRuntimeValue(value interface{}) *RuntimeValue {
	return &RuntimeValue{
		Value: value,
	}
}
func (e RuntimeValue) String() string {
	return fmt.Sprintf("%v", e.Value)
}

type Environment struct {
	Vars map[string]RuntimeValue
}

func NewEnvironment(vars map[string]RuntimeValue) *Environment {
	return &Environment{
		Vars: vars,
	}
}
func (e *Environment) GetVar(name string) *RuntimeValue {
	val, ok := e.Vars[name]
	if !ok {
		return nil
	}

	return &val
}

type Runtime struct {
	Envs []Environment
}

func NewRuntime(envs []Environment) *Runtime {
	return &Runtime{
		Envs: envs,
	}
}

type RuntimeError struct {
	Message string
	At      string
	Line    int
}

func NewRuntimeError(msg string, at string, line int) *RuntimeError {
	return &RuntimeError{
		Message: msg,
		At:      at,
		Line:    line,
	}
}
func (e RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] Error at '%s': %s", e.Line, e.At, e.Message)
}
