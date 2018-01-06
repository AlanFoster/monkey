package object

type Environment struct {
	values map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]Object),
	}
}

func (e *Environment) Add(identifier string, o Object) Object {
	e.values[identifier] = o
	return o
}

func (e *Environment) Get(identifier string) (Object, bool) {
	obj, ok := e.values[identifier]
	return obj, ok
}
