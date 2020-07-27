package ast

import ()

/*
type EvalContext struct {
	// All imported packages
	Packages map[string]*Program

	GlobalVariables map[string]Expression

	ScopeStack   []map[string]Expression
	currentScope int
}

func NewEvalContext() EvalContext {
	ctx := EvalContext{
		GlobalVariables: map[string]Expression{},
		ScopeStack:      []map[string]Expression{},
		currentScope:    -1, // Global Scope
	}
	return ctx
}

func (e *EvalContext) PushScope() {
	e.ScopeStack = append(e.ScopeStack, map[string]Expression{})
	e.currentScope += 1
}

func (e *EvalContext) PopScope() {
	if len(e.ScopeStack) > 0 {
		e.ScopeStack = e.ScopeStack[:len(e.ScopeStack)-1]
		e.currentScope -= 1
	}
}
*/
