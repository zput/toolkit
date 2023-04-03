package orm

import (
	"reflect"
	"strings"
	"unicode"
)

func Transfer(tag string, object interface{}) map[string]interface{} {
	return transfer(tag, object, ConvByGetFirstValue, ConvByCamel2Case)
}

type ConvertFunc = func(string) string

var ConvByGetFirstValue = func(in string) string { return strings.TrimSpace(strings.Split(in, ",")[0]) }
var ConvByCamel2Case = camel2Case

func transfer(tag string, Struct interface{}, converts ...ConvertFunc) map[string]interface{} {
	var ret = make(map[string]interface{})
	typeOfStruct := reflect.TypeOf(Struct)
	valueOfStruct := reflect.ValueOf(Struct)
	for i := typeOfStruct.NumField() - 1; i >= 0; i-- {
		fieldOfStruct := typeOfStruct.Field(i)
		if !fieldOfStruct.IsExported() {
			continue
		}
		var _tag = fieldOfStruct.Tag.Get(tag)
		if len(_tag) == 0 {
			_tag = fieldOfStruct.Name
		}
		for _, f := range converts {
			_tag = f(_tag)
		}
		ret[_tag] = valueOfStruct.Field(i).Interface()
	}
	return ret
}

func camel2Case(name string) string {
	buf := new(strings.Builder)
	for i, r := range name {
		if unicode.IsLower(r) {
			buf.WriteRune(r)
			continue
		}
		if i != 0 {
			buf.WriteByte('_')
		}
		buf.WriteRune(unicode.ToLower(r))
	}
	return buf.String()
}
