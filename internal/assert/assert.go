package assert

import (
	"database/sql"
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	"reflect"
	"time"
)

// compareTime 比较time
func compareTime(v1, v2 interface{}) (eq bool) {
	if value1, ok := v1.(time.Time); ok {
		value2, _ := v2.(time.Time)
		return value1.Equal(value2)
	}
	if value1, ok := v1.(sql.NullTime); ok {
		value2, _ := v2.(sql.NullTime)
		return value1.Valid == value2.Valid && value1.Time.Equal(value2.Time)
	}
	return
}

func equal(field string, data1, data2 reflect.Value) (ret bool) {
	// 0值存在一些特殊性，可能类型不一样
	if data1.IsZero() && data2.IsZero() {
		return true
	}
	if data1.Type() != data2.Type() {
		_, _ = convey.Println("断言失败，类型不一致", field, data1.Type(), data2.Type())
		return false
	}
	//比较int类型
	if compareInt(data1, data2) {
		return true
	}
	tp := data1.Type().Name()
	value1, value2 := data1.Interface(), data2.Interface()
	switch tp {
	case decimalType, nullDecimalType:
		ret = compareDecimal(value1, value2)
	case timeType, nullTimeType:
		ret = compareTime(value1, value2)
	default:
		ret = reflect.DeepEqual(value1, value2)
	}
	if !ret {
		msg := fmt.Sprintf("断言失败，字段：%s, 实际：%v，期望：%v", field, data1, data2)
		_, _ = convey.Println(msg)
	}
	return
}

func getTheRealVal(val reflect.Value) (isZero bool, ret reflect.Value) {
	if val.IsZero() {
		isZero = true
		return
	}
	return
}

func Assert(cur, exp interface{}) {
	data1Val := reflect.ValueOf(cur)
	data2Val := reflect.ValueOf(exp)

	if data1Val.Type() != data2Val.Type() {
		return
	}

	for i := 0; i < data1Val.NumField(); i++ {
		field := data2Val.Type().Field(i).Name
		data2Field := data2Val.Field(i)
		isZero, theVal := getTheRealVal(data2Field)
		if isZero {
			continue
		}
		data1Field := data1Val.Field(i)
		convey.So(equal(field, data1Field, theVal), convey.ShouldBeTrue)
	}
}
