package main

import (
	"fmt"

	"github.com/Azure/golua/lua"
	"github.com/Azure/golua/std"
	"github.com/jdolitsky/goluamapper"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	state := lua.NewState()
	defer state.Close()
	std.Open(state)

	// Evaluate the Lua script.
	// We can also use state.ExecFile(<filepath>)
	err := state.ExecText(`
		local name = "Fred"
		local age = 42

		person = {
			name = name,
			age = age
		}
	`)
	if err != nil {
		panic(err)
	}

	// Extract the "person" global var and map to Person type
	var person Person
	state.GetGlobal("person")
	err = goluamapper.Map(state.Pop(), &person)
	if err != nil {
		panic(err)
	}

	// Should print "Fred 42"
	fmt.Printf("%s %d\n", person.Name, person.Age)
}
