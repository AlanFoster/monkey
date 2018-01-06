package object

type Environment struct {
	values map[string]Object
	parent *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]Object),
		parent: nil,
	}
}

func NewClosedEnvironment(parent *Environment) *Environment {
	return &Environment{
		values: make(map[string]Object),
		parent: parent,
	}
}

func (e *Environment) Add(identifier string, o Object) Object {
	e.values[identifier] = o
	return o
}

func (e *Environment) Get(identifier string) (Object, bool) {
	obj, ok := e.values[identifier]
	if !ok && e.parent != nil {
		return e.parent.Get(identifier)
	}
	return obj, ok
}
