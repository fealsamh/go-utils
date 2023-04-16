package script

import (
	"fmt"
	"reflect"

	"github.com/fealsamh/go-utils/keyvalue"
	"github.com/fealsamh/go-utils/sexpr"
	"github.com/google/uuid"
)

// EvalContext is a node of a context in-tree.
type EvalContext struct {
	vars   map[string]any
	parent *EvalContext
}

// Set sets the value of a variable.
func (c *EvalContext) Set(name string, val any) {
	if c.vars == nil {
		c.vars = make(map[string]any)
	}
	c.vars[name] = val
}

// Get gets the value of a variable.
func (c *EvalContext) Get(name string) (any, error) {
	if val, ok := c.vars[name]; ok {
		return val, nil
	}
	if c.parent != nil {
		return c.parent.Get(name)
	}
	return nil, fmt.Errorf("unknown variable '%s'", name)
}

// Expr is a script expression.
type Expr interface {
	// Eval evaluates an expression.
	Eval(*EvalContext) (any, error)
}

// ConstantExpr is a constant literal.
type ConstantExpr[T any] struct {
	value T
}

// Eval evaluates an expression.
func (c *ConstantExpr[T]) Eval(ctx *EvalContext) (any, error) {
	return c.value, nil
}

// VariableExpr is a named variable.
type VariableExpr struct {
	name string
}

// Eval evaluates an expression.
func (v *VariableExpr) Eval(ctx *EvalContext) (any, error) {
	return ctx.Get(v.name)
}

// UUIDExpr is a `uuid` expression.
type UUIDExpr struct {
	arg Expr
}

// Eval evaluates an expression.
func (e *UUIDExpr) Eval(ctx *EvalContext) (any, error) {
	arg, err := e.arg.Eval(ctx)
	if err != nil {
		return nil, err
	}
	s, ok := arg.(string)
	if !ok {
		return nil, fmt.Errorf("argument of 'uuid' must be a string")
	}
	return uuid.Parse(s)
}

// NewExpr is a `new` expression.
type NewExpr struct {
	typ        reflect.Type
	properties map[string]Expr
}

// Eval evaluates an expression.
func (e *NewExpr) Eval(ctx *EvalContext) (any, error) {
	r := reflect.New(e.typ).Interface()
	a, err := keyvalue.NewObjectAdapter(r)
	if err != nil {
		return nil, err
	}
	for name, expr := range e.properties {
		arg, err := expr.Eval(ctx)
		if err != nil {
			return nil, err
		}
		if err := a.Put(name, arg); err != nil {
			return nil, err
		}
	}
	return r, nil
}

// WithExpr is a `with` expression.
type WithExpr struct {
	source     Expr
	properties map[string]Expr
}

// Eval evaluates an expression.
func (e *WithExpr) Eval(ctx *EvalContext) (any, error) {
	src, err := e.source.Eval(ctx)
	if err != nil {
		return nil, err
	}
	v1 := reflect.ValueOf(src)
	if v1.Kind() != reflect.Pointer || v1.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("first argument of 'with' must be an object, found '%v'", src)
	}
	typ := v1.Type().Elem()
	a1, err := keyvalue.NewObjectAdapter(src)
	if err != nil {
		return nil, err
	}
	r := reflect.New(typ).Interface()
	a2, err := keyvalue.NewObjectAdapter(r)
	if err != nil {
		return nil, err
	}
	a1.Pairs(func(key string, value interface{}) bool {
		if _, ok := e.properties[key]; !ok {
			if err := a2.Put(key, value); err != nil {
				panic(err) // this should never happen
			}
		}
		return true
	})
	for name, expr := range e.properties {
		arg, err := expr.Eval(ctx)
		if err != nil {
			return nil, err
		}
		if err := a2.Put(name, arg); err != nil {
			return nil, err
		}
	}
	return r, nil
}

// Translate translates an s-expression into an evaluatable AST.
func Translate(in any, types map[string]reflect.Type) (Expr, error) {
	switch in := in.(type) {
	case string:
		return &VariableExpr{name: in}, nil
	case int:
		return &ConstantExpr[int]{value: in}, nil
	case float64:
		return &ConstantExpr[float64]{value: in}, nil
	case sexpr.String:
		return &ConstantExpr[string]{value: string(in)}, nil
	case []any:
		if len(in) == 0 {
			return &ConstantExpr[[]any]{value: nil}, nil
		}
		fn, ok := in[0].(string)
		if !ok {
			return nil, fmt.Errorf("function name must be string, found '%v'", in[0])
		}
		switch fn {
		case "uuid":
			if len(in) != 2 {
				return nil, fmt.Errorf("expected one argument in function call '%s'", fn)
			}
			arg, err := Translate(in[1], types)
			if err != nil {
				return nil, err
			}
			return &UUIDExpr{arg: arg}, nil
		case "new":
			if len(in) < 2 {
				return nil, fmt.Errorf("too few arguments in function call '%s'", fn)
			}
			typeName, ok := in[1].(string)
			if !ok {
				return nil, fmt.Errorf("type name must be string, found '%v', in function call '%s'", in[1], fn)
			}
			typ, ok := types[typeName]
			if !ok {
				return nil, fmt.Errorf("unknown type '%s'", typeName)
			}
			props := make(map[string]Expr, len(in)-2)
			for _, pair := range in[2:] {
				pair, ok := pair.([]any)
				if !ok || len(pair) != 2 {
					return nil, fmt.Errorf("pair must be a list with two elements, found '%v'", pair)
				}
				name, ok := pair[0].(string)
				if !ok {
					return nil, fmt.Errorf("pair key must be string, found '%v'", pair[0])
				}
				expr, err := Translate(pair[1], types)
				if err != nil {
					return nil, err
				}
				props[name] = expr
			}
			return &NewExpr{typ: typ, properties: props}, nil
		case "with":
			if len(in) < 2 {
				return nil, fmt.Errorf("too few arguments in function call '%s'", fn)
			}
			src, err := Translate(in[1], types)
			if err != nil {
				return nil, err
			}
			props := make(map[string]Expr, len(in)-2)
			for _, pair := range in[2:] {
				pair, ok := pair.([]any)
				if !ok || len(pair) != 2 {
					return nil, fmt.Errorf("pair must be a list with two elements, found '%v'", pair)
				}
				name, ok := pair[0].(string)
				if !ok {
					return nil, fmt.Errorf("pair key must be string, found '%v'", pair[0])
				}
				expr, err := Translate(pair[1], types)
				if err != nil {
					return nil, err
				}
				props[name] = expr
			}
			return &WithExpr{source: src, properties: props}, nil
		default:
			return nil, fmt.Errorf("unknown function '%s'", fn)
		}
	default:
		return nil, fmt.Errorf("failed to translate '%v', unknown type '%T'", in, in)
	}
}
