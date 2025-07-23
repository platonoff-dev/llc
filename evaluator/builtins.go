package evaluator

import (
	"fmt"

	"anubis/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=%d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": {
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=%d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[0]
				}

				return NULL
			default:
				return newError("argument to `first` not supported, got %s", args[0].Type())
			}
		},
	},
	"last": {
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=%d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[len(arg.Elements)-1]
				}

				return NULL
			default:
				return newError("argument to `first` not supported, got %s", args[0].Type())
			}
		},
	},
	"rest": {
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=%d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return &object.Array{Elements: arg.Elements[1:len(arg.Elements)]}
				}

				return NULL
			default:
				return newError("argument to `first` not supported, got %s", args[0].Type())
			}
		},
	},
	"push": {
		Function: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=%d", len(args), 1)
			}

			if args[0].Type() != object.ArrayObj {
				return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
			}
			arr, _ := args[0].(*object.Array)
			length := len(arr.Elements)
			newElements := make([]object.Object, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		},
	},
	"print": {
		Function: func(args ...object.Object) object.Object {
			for _, a := range args {
				fmt.Println(a.Inspect())
			}

			return NULL
		},
	},
}
