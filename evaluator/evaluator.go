package evaluator

import (
	"bytes"
	"fmt"

	"github.com/labasubagia/interpreter/ast"
	"github.com/labasubagia/interpreter/object"
)

type ScopeType int

const (
	_ ScopeType = iota
	ScopeNone
	ScopeFunction
	ScopeLoop
)

var (
	NULL     = &object.Null{}
	TRUE     = &object.Boolean{Value: true}
	FALSE    = &object.Boolean{Value: false}
	BREAK    = &object.Break{}
	CONTINUE = &object.Continue{}
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Hash:
				return &object.Integer{Value: int64(len(arg.Pairs))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}

			return &object.Array{Elements: []object.Object{}}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			newElements := make([]object.Object, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		},
	},
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			var out bytes.Buffer
			for i, arg := range args {
				out.WriteString(arg.Inspect())
				if i < len(args)-1 {
					out.WriteString(" ")
				}
			}

			fmt.Println(out.String())
			return NULL
		},
	},
}

func Eval(node ast.Node, env *object.Environment, scope ScopeType) object.Object {

	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env, scope)
	case *ast.LetStatement:
		val := Eval(node.Value, env, scope)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.AssignExpression:
		return evalAssignExpression(node, env, scope)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env, scope)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env, scope)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env, scope)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env, scope)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.WhileStatement:
		return evalWhileStatement(node, env, scope)
	case *ast.BreakStatement:
		return BREAK
	case *ast.ContinueStatement:
		return CONTINUE
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.IfExpression:
		return evalIfExpression(node, env, scope)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env, scope)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env, scope)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.CallExpression:
		function := Eval(node.Function, env, scope)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env, scope)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env, scope)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env, scope)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env, scope)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env, scope)
	case *ast.Null:
		return NULL
	}

	return nil
}

func evalExpressions(exps []ast.Expression, env *object.Environment, scope ScopeType) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env, scope)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {

	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv, ScopeFunction)
		switch ev := evaluated.(type) {
		case *object.Break, *object.Continue:
			return newError("invalid keyword inside function: %s", ev.Type())
		}
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {

	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment, scope ScopeType) object.Object {
	condition := Eval(ie.Condition, env, scope)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env, scope)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env, scope)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalProgram(program *ast.Program, env *object.Environment, scope ScopeType) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env, scope)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment, scope ScopeType) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env, scope)

		if result != nil {
			rt := result.Type()
			switch rt {
			case object.RETURN_VALUE_OBJ, object.BREAK_OBJ, object.CONTINUE_OBJ, object.ERROR_OBJ:
				return result
			}
		}
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE, NULL:
		return TRUE
	default:
		if obj, ok := right.(*object.Integer); ok && obj.Value == 0 {
			return TRUE
		}
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalInfixIntegerExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalInfixStringExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalInfixIntegerExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+", "+=":
		return &object.Integer{Value: leftVal + rightVal}
	case "-", "-=":
		return &object.Integer{Value: leftVal - rightVal}
	case "*", "*=":
		return &object.Integer{Value: leftVal * rightVal}
	case "/", "/=":
		return &object.Integer{Value: leftVal / rightVal}
	case "%", "%=":
		return &object.Integer{Value: leftVal % rightVal}

	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalInfixStringExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)
	if idx < 0 || idx > max {
		return NULL
	}
	return arrayObject.Elements[idx]
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment, scope ScopeType) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env, scope)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env, scope)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalAssignExpression(node *ast.AssignExpression, env *object.Environment, scope ScopeType) object.Object {
	switch exp := node.Left.(type) {
	case *ast.Identifier:
		return evalIdentifierAssignExpression(exp, node.Operator, node.Value, env, scope)
	case *ast.IndexExpression:
		return evalIndexAssignExpression(exp, node.Operator, node.Value, env, scope)
	default:
		return newError("invalid identifier when assign value: %s", node.Left.String())
	}
}

func evalIdentifierAssignExpression(ident *ast.Identifier, operator string, value ast.Expression, env *object.Environment, scope ScopeType) object.Object {
	val := Eval(value, env, scope)
	if isError(val) {
		return val
	}

	cur, ok := env.Get(ident.Value)
	if !ok {
		return newError("identifier not found: %s", ident.Value)
	}

	if isCompoundAssignmentOperator(operator) {
		if !(cur.Type() == object.INTEGER_OBJ && val.Type() == object.INTEGER_OBJ) {
			return newError("unsupported assign %s %s %s", cur.Type(), operator, val.Type())
		}
		val = evalInfixExpression(operator, cur, val)
	}

	env.Assign(ident.Value, val)
	return val
}

func evalIndexAssignExpression(exp *ast.IndexExpression, operator string, value ast.Expression, env *object.Environment, scope ScopeType) object.Object {
	ident, ok := exp.Left.(*ast.Identifier)
	if !ok {
		return newError("invalid identifier using index")
	}

	cur, ok := env.Get(ident.Value)
	if !ok {
		return newError("identifier not found: %s", ident.Value)
	}

	index := Eval(exp.Index, env, scope)
	if isError(index) {
		return index
	}

	val := Eval(value, env, scope)
	if isError(val) {
		return val
	}

	switch {
	case cur.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexAssignExpression(ident, cur, index, operator, val, env)
	case cur.Type() == object.HASH_OBJ:
		return evalHashIndexAssignExpression(ident, cur, index, operator, val, env)
	default:
		return newError("index not supported: %s[%s]", cur.Type(), index.Type())
	}
}

func evalArrayIndexAssignExpression(
	ident *ast.Identifier,
	arr, index object.Object,
	operator string,
	val object.Object,
	env *object.Environment,
) object.Object {
	arrayObject := arr.(*object.Array)
	indexObject := index.(*object.Integer)

	n := len(arrayObject.Elements)
	if n == 0 {
		return newError("array is empty. cannot set at any index")
	}
	i := int(indexObject.Value)
	if i < 0 || i >= n {
		return newError("valid index range is 0 until %d. got=%d", n-1, i)
	}

	cur := arrayObject.Elements[i]
	if isCompoundAssignmentOperator(operator) {
		if !(cur.Type() == object.INTEGER_OBJ && val.Type() == object.INTEGER_OBJ) {
			return newError("unsupported assign %s[%s] -> %s %s %s", arr.Type(), index.Type(), cur.Type(), operator, val.Type())
		}
		val = evalInfixExpression(operator, cur, val)
	}

	arrayObject.Elements[i] = val
	env.Assign(ident.Value, arrayObject)
	return val
}

func evalHashIndexAssignExpression(
	ident *ast.Identifier,
	hash, index object.Object,
	operator string, val object.Object,
	env *object.Environment,
) object.Object {
	hashObject := hash.(*object.Hash)

	hashableKey, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}
	key := hashableKey.HashKey()

	if isCompoundAssignmentOperator(operator) {
		cur, ok := hashObject.Pairs[key]
		if !ok {
			return newError("cannot assign key not exist: %s[%s] %s %s", ident.Value, index.Inspect(), operator, val.Inspect())
		}
		if !(cur.Value.Type() == object.INTEGER_OBJ && val.Type() == object.INTEGER_OBJ) {
			return newError("unsupported assign %s[%s] -> %s %s %s", hashObject.Type(), cur.Key.Type(), cur.Value.Type(), operator, val.Type())
		}
		val = evalInfixExpression(operator, cur.Value, val)
	}

	hashObject.Pairs[key] = object.HashPair{
		Key:   index,
		Value: val,
	}
	env.Assign(ident.Value, hashObject)
	return val
}

func evalWhileStatement(node *ast.WhileStatement, env *object.Environment, scope ScopeType) object.Object {

	env = object.NewEnclosedEnvironment(env)

	condition := Eval(node.Condition, env, scope)
	if isError(condition) {
		return condition
	}
	for isTruthy(condition) {

		stmt := evalBlockStatement(node.Body, env, ScopeLoop)
		if stmt != nil {
			switch stmt.Type() {
			case object.BREAK_OBJ:
				return NULL
			case object.CONTINUE_OBJ:
				condition = Eval(node.Condition, env, ScopeLoop)
				if isError(condition) {
					return condition
				}
				continue
			case object.RETURN_VALUE_OBJ:
				if scope == ScopeFunction {
					return stmt
				} else {
					return newError("return statement unsupported if while-loop not inside a function")
				}
			case object.ERROR_OBJ:
				return stmt
			}
		}

		condition = Eval(node.Condition, env, scope)
		if isError(condition) {
			return condition
		}
	}
	return nil
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func isCompoundAssignmentOperator(operator string) bool {
	switch operator {
	case "+=", "-=", "/=", "*=", "%=":
		return true
	}
	return false
}
