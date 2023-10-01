package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	env := e

	for env != nil {
		o, exists := env.store[name]
		if exists {
			return o, true
		}

		env = e.outer
	}

	return nil, false
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
