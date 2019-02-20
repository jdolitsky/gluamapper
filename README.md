# goluamapper: maps an Azure/golua table to a Go struct

[![image](https://godoc.org/github.com/jdolitsky/goluamapper?status.svg)](http://godoc.org/github.com/jdolitsky/goluamapper)

goluamapper provides an easy way to map [Azure/golua](<https://github.com/Azure/golua>) tables to Go structs.

goluamapper converts an Azure/golua table to `map[string]interface{}`,
and then converts it to a Go struct using [mapstructure](https://github.com/mitchellh/mapstructure/).

## API

See [Go doc](http://godoc.org/github.com/jdolitsky/goluamapper).

## Usage

``` go
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
```

## License

MIT

## Original Author

Yusuke Inuzuka

Source: <https://github.com/yuin/goluamapper>