package object

type Environment struct {
	store map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
	}
}

func (e *Environment) Get(name string) (Object, bool) {
	o, exists := e.store[name]
	return o, exists
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
