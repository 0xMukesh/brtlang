package runtime

import (
	"fmt"

	"github.com/0xmukesh/interpreter/internal/ast"
)

type RuntimeValue struct {
	Value interface{}
}

type RuntimeVarMapping = map[string]RuntimeValue
type FuncMapping struct {
	Node ast.AstNode
	Args []ast.AstNode
}
type RuntimeFuncMapping = map[string]FuncMapping

func NewRuntimeValue(value interface{}) *RuntimeValue {
	return &RuntimeValue{
		Value: value,
	}
}
func (e RuntimeValue) String() string {
	return fmt.Sprintf("%v", e.Value)
}

type Environment struct {
	Parent *Environment
	Vars   RuntimeVarMapping
	Funcs  RuntimeFuncMapping
}

func NewEnvironment(vars RuntimeVarMapping, funcs RuntimeFuncMapping, parent *Environment) *Environment {
	return &Environment{
		Parent: parent,
		Vars:   vars,
		Funcs:  funcs,
	}
}
func (e *Environment) GetVar(name string) (*RuntimeValue, *Environment) {
	val, ok := e.Vars[name]
	if !ok {
		if e.Parent == nil {
			return nil, nil
		}

		return e.Parent.GetVar(name)
	}

	return &val, e
}
func (e *Environment) GetFunc(name string) (*FuncMapping, *Environment) {
	val, ok := e.Funcs[name]
	if !ok {
		if e.Parent == nil {
			return nil, nil
		}

		return e.Parent.GetFunc(name)
	}

	return &val, e
}
func (e *Environment) SetVar(name string, value RuntimeValue) {
	e.Vars[name] = value
}
func (e *Environment) SetFunc(name string, node ast.AstNode, args []ast.AstNode) {
	e.Funcs[name] = FuncMapping{
		Node: node,
		Args: args,
	}
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
