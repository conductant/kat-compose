package encoding

import (
	"reflect"
	"strings"
)

type OverrideFunc func(interface{}) interface{}

var (
	string_t = reflect.TypeOf(string(""))
)

// Function to transform a struct to a JSON-friendly map (map[string]interface{}).
// This will convert some problematic structures like map[t]interface{} where t is not a string.
// In that case the structure will be transformed to a slice to model a set.
func MarshalMap(this interface{}, tag string, overrides map[string]OverrideFunc) map[string]interface{} {
	m := map[string]interface{}{}
	thisValue := reflect.ValueOf(this)
	thisType := reflect.TypeOf(this)
	if thisType.Kind() == reflect.Ptr {
		thisValue = reflect.Indirect(thisValue)
		thisType = thisValue.Type()
	}
	if thisType.Kind() == reflect.Struct {
		for i := 0; i < thisValue.NumField(); i++ {
			ft := thisType.Field(i)
			fv := thisValue.Field(i)
			jf := strings.Split(ft.Tag.Get("json"), ",")[0] // json field name
			path := tag + "." + jf

			switch {
			case ft.Type.Kind() == reflect.Ptr:
				if ft.Type.Elem().Kind() == reflect.Struct {
					v := fv.Elem()
					if !fv.IsNil() {
						next := MarshalMap(v.Interface(), path, overrides)
						if jf == "" && ft.Anonymous {
							m = next
						} else {
							m[jf] = next
						}
					}
				}
			case ft.Type.Kind() == reflect.Struct:
				next := MarshalMap(fv.Interface(), path, overrides)
				if jf == "" && ft.Anonymous {
					m = next
				} else {
					m[jf] = next
				}
			default:
				if jf != "" {
					value := fv.Interface()

					if override, has := overrides[path]; has {

						value = override(value)

					} else if fv.Type().Kind() == reflect.Map && !fv.Type().Key().ConvertibleTo(string_t) {
						// this is a type that cannot be marshaled by json
						value = []interface{}{}
						for _, mk := range fv.MapKeys() {
							value = append(value.([]interface{}), mk.Interface())
						}
					}
					m[jf] = value
				}
			}
		}
	}
	return m
}
