// goluamapper provides an easy way to map GopherLua tables to Go structs.
package goluamapper

import (
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/Azure/golua/lua"
	"github.com/mitchellh/mapstructure"
)

// Option is a configuration that is used to create a new mapper.
type Option struct {
	// Function to convert a lua table key to Go's one. This defaults to "ToUpperCamelCase".
	NameFunc func(string) string

	// Returns error if unused keys exist.
	ErrorUnused bool

	// A struct tag name for lua table keys . This defaults to "goluamapper"
	TagName string
}

// Mapper maps a lua table to a Go struct pointer.
type Mapper struct {
	Option Option
}

// NewMapper returns a new mapper.
func NewMapper(opt Option) *Mapper {
	if opt.NameFunc == nil {
		opt.NameFunc = ToUpperCamelCase
	}
	if opt.TagName == "" {
		opt.TagName = "goluamapper"
	}
	return &Mapper{opt}
}

// Map maps the lua table to the given struct pointer.
func (mapper *Mapper) Map(v lua.Value, st interface{}) error {
	opt := mapper.Option
	mp := ToGoValue(v, opt)
	if mp.Kind() != reflect.Map {
		return errors.New("arguments #1 must be a table, but got an array")
	}
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           st,
		TagName:          opt.TagName,
		ErrorUnused:      opt.ErrorUnused,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(mp.Interface())
}

// Map maps the lua table to the given struct pointer with default options.
func Map(v lua.Value, st interface{}) error {
	return NewMapper(Option{}).Map(v, st)
}

// Id is an Option.NameFunc that returns given string as-is.
func Id(s string) string {
	return s
}

var camelre = regexp.MustCompile(`_([a-z])`)

// ToUpperCamelCase is an Option.NameFunc that converts strings from snake case to upper camel case.
func ToUpperCamelCase(s string) string {
	return strings.ToUpper(string(s[0])) + camelre.ReplaceAllStringFunc(s[1:len(s)], func(s string) string { return strings.ToUpper(s[1:len(s)]) })
}

// ToGoValue converts the given LValue to a Go object.
// adapted form https://github.com/Azure/golua/blob/master/pkg/luautil/reflect.go
func ToGoValue(v lua.Value, opt Option) reflect.Value {
	switch v := v.(type) {
	case *lua.Object:
		return reflect.ValueOf(v.Value())
	case lua.Table:
		return tableToGo(v, opt)
	case lua.String:
		return reflect.ValueOf(string(v))
	case lua.Float:
		return reflect.ValueOf(float64(v))
	case lua.Int:
		return reflect.ValueOf(int64(v))
	case lua.Bool:
		return reflect.ValueOf(bool(v))
	}
	return reflect.ValueOf(nil)
}

// adapted form https://github.com/Azure/golua/blob/master/pkg/luautil/reflect.go
func tableToGo(table lua.Table, opt Option) reflect.Value {
	if length := table.Length(); length == 0 { // map
		gomap := make(map[interface{}]interface{})
		table.ForEach(func(key, val lua.Value) {
			v := ToGoValue(val, opt)
			if !isZeroOfUnderlyingType(v) {
				k := ToGoValue(key, opt)
				gomap[opt.NameFunc(k.String())] = v.Interface()
			}
		})
		return reflect.ValueOf(gomap)
	} else { // slice
		slice := make([]interface{}, 0, length)
		for i := 1; i <= length; i++ {
			elem := ToGoValue(table.Index(lua.Int(i)), opt)
			slice = append(slice, elem.Interface())
		}
		return reflect.ValueOf(slice)
	}
}

// adaped from https://stackoverflow.com/a/13906031
func isZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
