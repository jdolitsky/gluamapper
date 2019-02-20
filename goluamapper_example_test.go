package goluamapper

import (
	"bytes"
	"fmt"
	"github.com/Azure/golua/lua"
	"github.com/Azure/golua/std"
)

func ExampleMap() {
	type Role struct {
		Name string
	}

	type Person struct {
		Name      string
		Age       int
		WorkPlace string
		Role      []*Role
	}

	debug := true
	opts := []lua.Option{lua.WithTrace(debug), lua.WithVerbose(debug)}
	state := lua.NewState(opts...)
	defer state.Close()
	std.Open(state)

	err := state.ExecFrom(bytes.NewReader([]byte(`
		person = {
      		name = "Michel",
      		age  = "31", -- weakly input
			work_place = "San Jose",
      		role = {
        		{
          			name = "Administrator"
        		},
        		{
          			name = "Operator"
        		}
      		}
    	}
	`)))
	if err != nil {
		panic(err)
	}

	var person Person

	state.GetGlobal("person")
	table := state.Pop().(lua.Table)

	if err := Map(table, &person); err != nil {
		panic(err)
	}
	fmt.Printf("%s %d", person.Name, person.Age)
	// Output:
	// Michel 31
}
