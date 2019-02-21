package main

import (
	"fmt"
	"os"

	"github.com/Azure/golua/lua"
	"github.com/Azure/golua/std"
	"github.com/jdolitsky/goluamapper"
	"gopkg.in/yaml.v2"
)

type (
	ExecMixin struct {
		Command   string
		Arguments []string
	}

	Action struct {
		Description string
		Exec        ExecMixin
	}

	Bundle struct {
		Name            string
		Version         string
		Description     string
		InvocationImage string
		Mixins          []string
		Install         []*Action
		Uninstall       []*Action
	}
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	state := lua.NewState()
	defer state.Close()
	std.Open(state)

	// Evaluate the Lua script
	path := os.Args[1]
	err := state.ExecFile(path)
	check(err)

	// Extract the "bundle" global var and map to Bundle type
	var bundle Bundle
	state.GetGlobal("bundle")
	err = goluamapper.Map(state.Pop(), &bundle)
	check(err)

	// Convert to YAML and print to stdout
	out, err := yaml.Marshal(bundle)
	check(err)
	fmt.Println(string(out))
}
