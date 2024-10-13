package click

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Value struct {
	inner interface{}
}

func (v *Value) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &v.inner)
}

func (v Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.inner)
}

func (v Value) check(expectedType string) {
	if v.IsNil() {
		panic("cannot convert nil Value to " + expectedType)
	}
}

func (v Value) IsNil() bool {
	return v.inner == nil
}

func (v Value) AsFloat64(defaultValue ...float64) float64 {
	if v.IsNil() && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	v.check("Float64")
	rv := reflect.ValueOf(v.inner)
	if rv.CanFloat() {
		return rv.Float()
	} else if rv.CanInt() {
		return float64(rv.Int())
	} else if rv.CanUint() {
		return float64(rv.Uint())
	} else if rv.Kind() == reflect.Bool {
		if rv.Bool() {
			return 1
		}
		return 0
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	panic(fmt.Sprintf("cannot convert %v value to Float64", rv.Type()))
}

func (v Value) AsInt64(defaultValue ...int64) int64 {
	if v.IsNil() && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	v.check("Int64")
	rv := reflect.ValueOf(v.inner)
	if rv.CanFloat() {
		return int64(rv.Float())
	} else if rv.CanInt() {
		return rv.Int()
	} else if rv.CanUint() {
		return int64(rv.Uint())
	} else if rv.Kind() == reflect.Bool {
		if rv.Bool() {
			return 1
		}
		return 0
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	panic(fmt.Sprintf("cannot convert %v value to Int64", rv.Type()))
}

func (v Value) AsString(defaultValue ...string) string {
	if v.IsNil() && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	v.check("String")
	rv := reflect.ValueOf(v.inner)
	if rv.Kind() == reflect.String {
		return rv.String()
	}
	return fmt.Sprintf("%v", v.inner)
}

func (v Value) AsInterface(defaultValue ...interface{}) interface{} {
	if v.IsNil() && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return v.inner
}
