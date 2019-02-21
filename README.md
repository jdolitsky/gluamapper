# goluamapper: maps an Azure/golua table to a Go struct

[![image](https://godoc.org/github.com/jdolitsky/goluamapper?status.svg)](http://godoc.org/github.com/jdolitsky/goluamapper)

goluamapper provides an easy way to map [Azure/golua](<https://github.com/Azure/golua>) tables to Go structs.

goluamapper converts an Azure/golua table to `map[string]interface{}`,
and then converts it to a Go struct using [mapstructure](https://github.com/mitchellh/mapstructure/).

## API

See [Go doc](http://godoc.org/github.com/jdolitsky/goluamapper).

## Usage

See the [source](./examples/readme/main.go) for the example below.

This example will evaluate a Lua script (as string), extract the resulting `person` global variable, and map it to our custom `Person` type:

``` go
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

```

## License

MIT

## Original Author

Yusuke Inuzuka

Source: <https://github.com/yuin/gluamapper>
