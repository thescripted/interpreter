package main

type Environment struct {
	mem map[string]interface{}
}

func (e Environment) get(name string) interface{} {
	return e.mem[name] // should type/error check.
}

func (e Environment) put(name string, value interface{}) {
	e.mem[name] = value
}
