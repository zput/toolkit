package orm

import (
	"reflect"
	"strings"
	"unicode"
)

func TransferByDefaultConvAndTags(object interface{}) map[string]interface{} {
	return TransferByDefaultConv("orm", object)
}

// TransferByDefaultConv default sets convert functions including ConvByGetFirstValue etc.
func TransferByDefaultConv(tag string, object interface{}) map[string]interface{} {
	return NewTransfer(tag, ConvByGetFirstValue).transfer(object)
}

type ConvertFunc = func(string) string

var ConvByGetFirstValue = func(in string) string { return strings.TrimSpace(strings.Split(in, ",")[0]) }
var ConvByCamel2Case = camel2Case

func NewTransfer(tag string, converts ...ConvertFunc) *Transfer {
	return &Transfer{
		tag:      tag,
		converts: converts,
	}
}

type Transfer struct {
	tag      string
	converts []ConvertFunc
}

func (tran *Transfer) transfer(Struct interface{}) map[string]interface{} {
	var ret = make(map[string]interface{})
	t := reflect.TypeOf(Struct)
	v := reflect.ValueOf(Struct)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ret
	}
	tran.parse(ret, t, v)
	return ret
}

func (tran *Transfer) parse(storage map[string]interface{}, t reflect.Type, v reflect.Value) {
	for i := t.NumField() - 1; i >= 0; i-- {
		tt := t.Field(i)
		vv := v.Field(i)
		if !tt.IsExported() {
			continue
		}
		if tt.Anonymous {
			tran.parse(storage, tt.Type, vv)
			continue
		}
		var _tag = tt.Tag.Get(tran.tag)
		if len(_tag) == 0 {
			_tag = ConvByCamel2Case(tt.Name)
		}
		for _, f := range tran.converts {
			_tag = f(_tag)
		}
		storage[_tag] = vv.Interface()
	}
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
