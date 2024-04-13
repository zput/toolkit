package assert

import (
	"fmt"
	"reflect"
)

func Equal(cur, exp interface{}) bool {
	var isEqual = reflect.DeepEqual(cur, exp)
	if isEqual {
		return true
	}
	return equal(cur, exp)
}

func equal(cur, exp interface{}) bool {
	v1 := reflect.ValueOf(cur)
	v2 := reflect.ValueOf(exp)
	if v1.Type() != v2.Type() {
		return false
	}
	if v1.Kind() == reflect.Ptr {
		v1 = v1.Elem()
		v2 = v2.Elem()
	}

	// 提前判断自定义类型
	if isSelfType(v1.Type()) {
		return compareSelfType(v1.Type().Name(), v1, v2)
	}

	switch v2.Kind() {
	case reflect.Slice:
		if v1.Len() != v2.Len() {
			fmt.Println("断言失败，字段长度不一致", v1.Len(), v2.Len())
			return false
		}
		for j := 0; j < v2.Len(); j++ {
			if !Equal(v1.Index(j).Interface(), v2.Index(j).Interface()) {
				return false
			}
		}
		return true
	case reflect.Struct:
		return Struct(v1.Interface(), v2.Interface())
	}

	return reflect.DeepEqual(v1.Interface(), v2.Interface())
}

func Struct(cur, exp interface{}) bool {
	v1 := reflect.ValueOf(cur)
	v2 := reflect.ValueOf(exp)
	if v1.Type() != v2.Type() || v2.Kind() != reflect.Struct {
		return false
	}

	fmt.Printf("compare [%s] type, num field：%d\n", v2.Type().String(), v2.Type().NumField())

	for i := 0; i < v2.NumField(); i++ {
		fieldName := v2.Type().Field(i).Name
		fmt.Printf("\t %d: %s\n", i, fieldName)

		// 1.判断v1的字段值
		vByV1I := v1.Field(i)
		// 2.判断v2的字段是否为零值
		vByV2I := v2.Field(i)
		isZero, theVal := getTheRealVal(vByV2I)
		if isZero {
			continue
		}
		// 3.获取v2字段名
		if !equal(vByV1I.Interface(), theVal.Interface()) {
			return false
		}
	}
	return true
}

func getTheRealVal(val reflect.Value) (isZero bool, ret reflect.Value) {
	if val.IsZero() {
		isZero = true
		return
	}
	if val.CanInterface() {
		ret = val
	}
	return
}
