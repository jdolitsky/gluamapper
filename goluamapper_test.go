package goluamapper

import (
	"bytes"
	"fmt"
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
	WorkPlace string `goluamapper:"w"`
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
	state := lua.NewState()
	defer state.Close()
	std.Open(state)

	err := state.ExecFrom(bytes.NewReader([]byte(`
		person = {
      		name = "Michel",
      		age  = "31", -- weakly input
			x    = 100,
			w    = "San Jose",
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
	v := state.Pop()

	if err := Map(v, &person); err != nil {
		t.Error(err)
	}
	errorIfNotEqual(t, "Michel", person.Name)
	errorIfNotEqual(t, 31, person.Age)
	errorIfNotEqual(t, 100, person.X)
	errorIfNotEqual(t, "San Jose", person.WorkPlace)
	errorIfNotEqual(t, 2, len(person.Role))
	errorIfNotEqual(t, "Administrator", person.Role[0].Name)
	errorIfNotEqual(t, "Operator", person.Role[1].Name)
}

func TestTypes(t *testing.T) {
	state := lua.NewState()
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
	v := state.Pop()

	if err := NewMapper(Option{NameFunc: Id}).Map(v, &stct); err != nil {
		t.Error(err)
	}
	errorIfNotEqual(t, nil, stct.Nil)
	errorIfNotEqual(t, true, stct.Bool)
	errorIfNotEqual(t, "string", stct.String)
	errorIfNotEqual(t, 10, stct.Number)
}

func TestNameFunc(t *testing.T) {
	state := lua.NewState()
	defer state.Close()
	std.Open(state)

	err := state.ExecFrom(bytes.NewReader([]byte(`
		person = {
			name = "Michel",
			age  = "31", -- weakly input
			x    = 100,
			w    = "San Jose",
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
	v := state.Pop()

	mapper := NewMapper(Option{NameFunc: Id})
	if err := mapper.Map(v, &person); err != nil {
		t.Error(err)
	}
	errorIfNotEqual(t, "Michel", person.Name)
	errorIfNotEqual(t, 31, person.Age)
	errorIfNotEqual(t, 100, person.X)
	errorIfNotEqual(t, "San Jose", person.WorkPlace)
	errorIfNotEqual(t, 2, len(person.Role))
	errorIfNotEqual(t, "Administrator", person.Role[0].Name)
	errorIfNotEqual(t, "Operator", person.Role[1].Name)
}

func TestError(t *testing.T) {
	state := lua.NewState()
	defer state.Close()
	std.Open(state)

	tmpTable := make(map[interface{}]interface{})
	tmpTable["key"] = "value"
	v := lua.ValueOf(state, tmpTable)
	err := Map(v, 1)
	if err.Error() != "result must be a pointer" {
		t.Error("invalid error message")
	}

	var person testPerson
	err = Map(lua.ValueOf(state, []string{"hello"}), &person)
	if err.Error() != "arguments #1 must be a table, but got an array" {
		fmt.Println(err.Error())
		t.Error("invalid error message")
	}
}
