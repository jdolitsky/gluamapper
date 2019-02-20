package goluamapper

import (
	"bytes"
	"github.com/Azure/golua/lua"
	"github.com/Azure/golua/std"

	"path/filepath"
	"runtime"
	"testing"
)

func errorIfNotEqual(t *testing.T, v1, v2 interface{}) {
	if v1 != v2 {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("%v line %v: '%v' expected, but got '%v'", filepath.Base(file), line, v1, v2)
	}
}

type testRole struct {
	Name string
}

type testPerson struct {
	Name      string
	Age       int
	WorkPlace string `goluamapper:"work_place"`
	Role      []*testRole
	X         int
}

type testStruct struct {
	Nil    interface{}
	Bool   bool
	String string
	Number int `goluamapper:"number_value"`
	Func   interface{}
}

func TestMap(t *testing.T) {
	debug := false
	opts := []lua.Option{lua.WithTrace(debug), lua.WithVerbose(debug)}
	state := lua.NewState(opts...)
	defer state.Close()
	std.Open(state)

	err := state.ExecFrom(bytes.NewReader([]byte(`
		person = {
      		name = "Michel",
      		age  = "31", -- weakly input
			x    = 100,
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
		t.Error(err)
	}

	var person testPerson

	state.GetGlobal("person")
	table := state.Pop().(lua.Table)

	if err := Map(table, &person); err != nil {
		t.Error(err)
	}
	errorIfNotEqual(t, "Michel", person.Name)
	errorIfNotEqual(t, 31, person.Age)
	errorIfNotEqual(t, 100, person.X)
	errorIfNotEqual(t, "San Jose", person.WorkPlace)
	/*errorIfNotEqual(t, 2, len(person.Role))
	errorIfNotEqual(t, "Administrator", person.Role[0].Name)
	errorIfNotEqual(t, "Operator", person.Role[1].Name)*/
}

func TestTypes(t *testing.T) {
	debug := false
	opts := []lua.Option{lua.WithTrace(debug), lua.WithVerbose(debug)}
	state := lua.NewState(opts...)
	defer state.Close()
	std.Open(state)

	err := state.ExecFrom(bytes.NewReader([]byte(`
        tbl = {
            ["Nil"] = nil,
            ["Bool"] = true,
            ["String"] = "string",
            ["Number_value"] = 10,
            ["Func"] = function() end
        }
	`)))
	if err != nil {
		t.Error(err)
	}

	var stct testStruct

	state.GetGlobal("tbl")
	table := state.Pop().(lua.Table)

	if err := NewMapper(Option{NameFunc: Id}).Map(table, &stct); err != nil {
		t.Error(err)
	}
	errorIfNotEqual(t, nil, stct.Nil)
	errorIfNotEqual(t, true, stct.Bool)
	errorIfNotEqual(t, "string", stct.String)
	errorIfNotEqual(t, 10, stct.Number)
}

func TestNameFunc(t *testing.T) {
	debug := true
	opts := []lua.Option{lua.WithTrace(debug), lua.WithVerbose(debug)}
	state := lua.NewState(opts...)
	defer state.Close()
	std.Open(state)

	err := state.ExecFrom(bytes.NewReader([]byte(`
person = {
	name = "Michel",
	age  = "31", -- weakly input
	x    = 100,
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
		t.Error(err)
	}

	var person testPerson

	state.GetGlobal("person")
	table := state.Pop().(lua.Table)

	mapper := NewMapper(Option{NameFunc: Id})
	if err := mapper.Map(table, &person); err != nil {
		t.Error(err)
	}
	errorIfNotEqual(t, "Michel", person.Name)
	errorIfNotEqual(t, 31, person.Age)
	errorIfNotEqual(t, 100, person.X)
	errorIfNotEqual(t, "San Jose", person.WorkPlace)
	/*errorIfNotEqual(t, 2, len(person.Role))
	errorIfNotEqual(t, "Administrator", person.Role[0].Name)
	errorIfNotEqual(t, "Operator", person.Role[1].Name)*/
}

/*

// TODO: fix and re-enable this test
func TestError(t *testing.T) {
	debug := true
	opts := []lua.Option{lua.WithTrace(debug), lua.WithVerbose(debug)}
	state := lua.NewState(opts...)
	defer state.Close()
	std.Open(state)

	tbl := lua.NewTable()
	L.SetField(tbl, "key", lua.LString("value"))
	err := Map(tbl, 1)
	if err.Error() != "result must be a pointer" {
		t.Error("invalid error message")
	}

	tbl = L.NewTable()
	tbl.Append(lua.LNumber(1))
	var person testPerson
	err = Map(tbl, &person)
	if err.Error() != "arguments #1 must be a table, but got an array" {
		t.Error("invalid error message")
	}
}
*/
