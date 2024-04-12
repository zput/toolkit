package assert

import (
	"database/sql"
	"github.com/shopspring/decimal"
	"reflect"
	"time"
)

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
