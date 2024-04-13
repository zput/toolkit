package assert

import (
	"database/sql"
	"fmt"
	"github.com/shopspring/decimal"
	"reflect"
	"time"
)

var (
	decimalType     = reflect.TypeOf(decimal.Decimal{}).Name()
	nullDecimalType = reflect.TypeOf(decimal.NullDecimal{}).Name()
	timeType        = reflect.TypeOf(time.Time{}).Name()
	nullTimeType    = reflect.TypeOf(sql.NullTime{}).Name()
)

// compareInt 比较int类型
func compareInt(v1, v2 reflect.Value) bool {
	if v1.CanInt() && v2.CanInt() {
		return v1.Int() == v2.Int()
	}
	if v1.CanUint() && v2.CanUint() {
		return v1.Uint() == v2.Uint()
	}
	return false
}

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

// compareDecimal 比较decimal
func compareDecimal(v1, v2 interface{}) (eq bool) {
	if value1, ok := v1.(decimal.Decimal); ok {
		value2, _ := v2.(decimal.Decimal)
		return value1.Equal(value2)
	}
	if value1, ok := v1.(decimal.NullDecimal); ok {
		value2, _ := v2.(decimal.NullDecimal)
		return value1.Valid == value2.Valid && value1.Decimal.Equal(value2.Decimal)
	}
	return
}
func isSelfType(t reflect.Type) bool {
	tp := t.Name()
	switch tp {
	case decimalType, nullDecimalType:
		fallthrough
	case timeType, nullTimeType:
		return true
	default:
		return false
	}
}

func compareSelfType(field string, data1, data2 reflect.Value) (ret bool) {
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
		fmt.Printf("断言失败，字段：%s, 实际：%v，期望：%v", field, data1, data2)
	}
	return
}
