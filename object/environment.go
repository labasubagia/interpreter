package object

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// Assign only if exists anywhere in inner or outer env
func (e *Environment) Assign(name string, val Object) (value Object, found bool) {
	env := e
	for env != nil {
		if _, ok := env.store[name]; ok {
			env.store[name] = val
			return val, true
		}
		env = env.outer
	}
	return val, false
}
