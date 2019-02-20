// goluamapper provides an easy way to map GopherLua tables to Go structs.
package goluamapper

import (
	"fmt"
	"regexp"
	"strconv"
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
func (mapper *Mapper) Map(table lua.Table, st interface{}) error {
	opt := mapper.Option
	mp := ToGoValue(table, opt)
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
	return decoder.Decode(mp)
}

// Map maps the lua table to the given struct pointer with default options.
func Map(table lua.Table, st interface{}) error {
	return NewMapper(Option{}).Map(table, st)
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
func ToGoValue(table lua.Table, opt Option) map[interface{}]interface{} {
	m := map[interface{}]interface{}{}
	table.ForEach(func(key lua.Value, value lua.Value) {
		fmt.Println(opt.NameFunc(key.String()))
		switch value.Type() {
		case lua.NilType:
			m[opt.NameFunc(key.String())] = nil
		case lua.BoolType:
			m[opt.NameFunc(key.String())] = true
		case lua.StringType:
			fmt.Println("hi", value.String())
			m[opt.NameFunc(key.String())] = value.String()
		case lua.NumberType:
			v, err := strconv.ParseInt(value.String(), 10, 64)
			if err == nil {
				m[opt.NameFunc(key.String())] = v
			}
		}
	})
	return m
}

/*

	switch lv.Type() {
	case lua.NilType:
		return nil
	case lua.BoolType:
		return true //bool(lv)
	case lua.StringType:
		return "" //string(v)
	case lua.NumberType:
		return int64(0) //v
	case lua.TableType:
		maxn := 0 //v.Length()
		if maxn == 0 { // table
			ret := make(map[interface{}]interface{})
			/*
			lv.ForEach(func(key, value lua.Value) {
				keystr := fmt.Sprint(ToGoValue(key, opt))
				ret[opt.NameFunc(keystr)] = ToGoValue(value, opt)
			})

			return ret
		} else { // array
			ret := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				ret = append(ret, ToGoValue(lua.Int(i), opt))
			}
			return ret
		}
	default:
		return m
	}
}
			*/