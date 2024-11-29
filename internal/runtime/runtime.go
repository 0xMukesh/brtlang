package runtime

import (
	"fmt"
)

type RuntimeValue struct {
	Value interface{}
}

type RuntimeVarMapping = map[string]RuntimeValue

func NewRuntimeValue(value interface{}) *RuntimeValue {
	return &RuntimeValue{
		Value: value,
	}
}
func (e RuntimeValue) String() string {
	return fmt.Sprintf("%v", e.Value)
}

type Environment struct {
	Vars RuntimeVarMapping
}

func NewEnvironment(vars RuntimeVarMapping) *Environment {
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
func (e *Environment) SetVar(name string, value RuntimeValue) {
	e.Vars[name] = value
}

type Runtime struct {
	Envs *[]Environment
}

func NewRuntime(envs *[]Environment) *Runtime {
	return &Runtime{
		Envs: envs,
	}
}
func (r *Runtime) AddNewEnv(env Environment) {
	if r.Envs != nil {
		*r.Envs = append(*r.Envs, env)
	}
}
func (r *Runtime) RemoveLastEnv() {
	if r.Envs != nil {
		*r.Envs = (*r.Envs)[:len(*r.Envs)-1]
	}
}
func (r *Runtime) CurrEnv() *Environment {
	if r.Envs != nil {
		return &(*r.Envs)[len(*r.Envs)-1]
	}
	return nil
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
