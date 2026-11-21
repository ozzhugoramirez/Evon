package compiler

import "fmt"

var NullValue = &Null{}

// --- LIBRERÍA ESTÁNDAR (Simulada) ---
// Aquí vivirán 'puts', 'sys.Imprimir', etc. Son funciones inyectadas
// en el entorno, separadas del núcleo, respetando la filosofía C.
type BuiltinFunction func(args ...Object) Object

var builtins = map[string]BuiltinFunction{
	"puts": func(args ...Object) Object {
		for _, arg := range args {
			if arg != nil && arg != NullValue {
				fmt.Println(arg.Inspect())
			}
		}
		return NullValue
	},
}

func Eval(node Node) Object {
	switch node := node.(type) {
	case *Program:
		return evalStatements(node.Statements)

	// Evaluación de valores puros
	case *IntegerLiteral:
		return &Integer{Value: node.Value}
	case *StringLiteral:
		return &String{Value: node.Value}
	case *Identifier:
		// Si es un identificador, buscamos si es una función de la librería estándar
		if _, ok := builtins[node.Value]; ok {
			return &Builtin{Name: node.Value}
		}
		// En el futuro, aquí buscaremos variables reales en la Bóveda .axm
		fmt.Printf("[SISTEMA INMUNE] Error: La referencia '%s' no existe en el contexto actual.\n", node.Value)
		return NullValue

	case *LetStatement:
		mut := "let"
		if node.IsMutable {
			mut = "mut"
		}
		val := Eval(node.Value)
		fmt.Printf("[MEMORIA TEMPORAL] Guardado -> [%s] %s : %s = %s\n", mut, node.Name.Value, node.ValueType.Literal, val.Inspect())
		return NullValue

	case *DenyStatement:
		fmt.Printf("[SISTEMA INMUNE] Evaluando Ley: %s\n", node.Condition)
		return NullValue

	case *ExpressionStatement:
		return Eval(node.Expression)

	case *CallExpression:
		// 1. Evaluamos la función (ej. buscando 'puts')
		function := Eval(node.Function)
		if function == NullValue {
			return NullValue // El error ya se reportó al evaluar el Identifier
		}

		// 2. Evaluamos los argumentos matemáticamente
		args := evalExpressions(node.Arguments)
		if len(args) == 1 && args[0] == NullValue {
			return NullValue
		}

		// 3. Ejecutamos (Aplicamos) la llamada
		return applyFunction(function, args)
	}
	return nil
}

func evalStatements(stmts []Statement) Object {
	var result Object
	for _, statement := range stmts {
		result = Eval(statement)
	}
	return result
}

func evalExpressions(exps []Expression) []Object {
	var result []Object
	for _, e := range exps {
		evaluated := Eval(e)
		result = append(result, evaluated)
	}
	return result
}

func applyFunction(fn Object, args []Object) Object {
	if builtin, ok := fn.(*Builtin); ok {
		// Ejecutamos la función puente de Go
		function := builtins[builtin.Name]
		return function(args...)
	}

	fmt.Printf("[SISTEMA INMUNE] Error Crítico: '%s' no es invocable (no es una función).\n", fn.Inspect())
	return NullValue
}
