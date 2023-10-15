package evaluator

import "github.com/gkampitakis/monkey/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: len(arg.Value)}
			case *object.Array:
				return &object.Integer{Value: len(arg.Elements)}
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

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[0]
				}
				return NULL
			default:
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[len(arg.Elements)-1]
				}
				return NULL
			default:
				return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
			}
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					newElements := make([]object.Object, len(arg.Elements)-1)
					copy(newElements, arg.Elements[1:])
					return &object.Array{Elements: newElements}
				}
				return NULL
			default:
				return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) >= 0 {
					newElements := make([]object.Object, 0, len(arg.Elements)+len(args)-1)
					copy(newElements, arg.Elements)
					return &object.Array{Elements: append(newElements, args[1:]...)}
				}
				return NULL
			default:
				return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
			}
		},
	},
}
