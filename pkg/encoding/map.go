package encoding

import (
	"errors"
	"fmt"
	"path"
	"reflect"
	"strings"
)

type OverrideFunc func(interface{}) interface{}

var (
	string_t = reflect.TypeOf(string(""))
	bool_t   = reflect.TypeOf(true)
)

// Common Go idiom of implementing set<x> as map[*x]bool.  This is problematic for JSON decoding and encoding.
func IsValueASet(fv reflect.Value) bool {
	if fv.Type().Kind() != reflect.Map {
		return false
	}
	return !fv.Type().Key().ConvertibleTo(string_t) && fv.Type().Elem().Kind() == reflect.Bool
}

type Codec struct {
	marshalOverrides   map[string]OverrideFunc
	unmarshalOverrides map[string]OverrideFunc
}

func NewCodec() *Codec {
	return &Codec{
		marshalOverrides:   map[string]OverrideFunc{},
		unmarshalOverrides: map[string]OverrideFunc{},
	}
}

// Override function when marshalling value of the given type with given field name.
func (this *Codec) OverrideMarshal(s interface{}, fieldName string, f OverrideFunc) *Codec {
	ft := reflect.TypeOf(s)
	if _, has := ft.FieldByName(fieldName); !has {
		panic(fmt.Errorf("no such field:%v, %s", ft, fieldName))
	}
	this.marshalOverrides[path.Join(ft.PkgPath(), ft.Name(), fieldName)] = f
	return this
}

// Override function when unmarshalling value of the given type with given field name.
func (this *Codec) OverrideUnmarshal(s interface{}, fieldName string, f OverrideFunc) *Codec {
	ft := reflect.TypeOf(s)
	if _, has := ft.FieldByName(fieldName); !has {
		panic(fmt.Errorf("no such field:%v, %s", ft, fieldName))
	}
	this.unmarshalOverrides[path.Join(ft.PkgPath(), ft.Name(), fieldName)] = f
	return this
}

func (this *Codec) MarshalMap(s interface{}) map[string]interface{} {
	return MarshalMap(s, this.marshalOverrides)
}

func (this *Codec) UnmarshalMap(input map[string]interface{}, s interface{}) {
	UnmarshalMap(input, s, this.unmarshalOverrides)
}

// Function to transform a struct to a JSON-friendly map (map[string]interface{}).
// This will convert some problematic structures like map[t]interface{} where t is not a string.
// In that case the structure will be transformed to a slice to model a set.
func MarshalMap(this interface{}, overrides map[string]OverrideFunc) map[string]interface{} {
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

			switch {
			case ft.Type.Kind() == reflect.Ptr && ft.Type.Elem().Kind() == reflect.Struct:
				if !fv.IsNil() {
					next := MarshalMap(fv.Elem().Interface(), overrides)
					if jf == "" && ft.Anonymous {
						m = next
					} else {
						m[jf] = next
					}
				}

			case ft.Type.Kind() == reflect.Struct:
				next := MarshalMap(fv.Interface(), overrides)
				if jf == "" && ft.Anonymous {
					m = next
				} else {
					m[jf] = next
				}

			default:
				if jf != "" {
					value := fv.Interface()
					k := path.Join(thisType.PkgPath(), thisType.Name(), ft.Name)
					if override, has := overrides[k]; has {
						value = override(value)
					} else if IsValueASet(fv) {
						// this is a type that cannot be marshaled by json
						sl := reflect.MakeSlice(reflect.SliceOf(fv.Type().Key()), 0, 0)
						for _, mk := range fv.MapKeys() {
							sl = reflect.Append(sl, mk)
						}
						value = sl.Interface()
					}
					m[jf] = value
				}
			}
		}
	}
	return m
}

// Function to unmarshal the content of a JSON-friendly map onto the provided struct.  This assumes that
// the structure of the map conforms to that generated by MarshalMap where Set<X> idioms are transformed into slices.
func UnmarshalMap(input map[string]interface{}, this interface{}, overrides map[string]OverrideFunc) {
	thisValue := reflect.ValueOf(this)
	thisType := reflect.TypeOf(this)
	if thisType.Kind() == reflect.Ptr {
		thisValue = reflect.Indirect(thisValue)
		thisType = thisValue.Type()
	} else {
		panic(errors.New("Input must be a pointer."))
	}

	if thisType.Kind() == reflect.Struct {
		for i := 0; i < thisValue.NumField(); i++ {
			ft := thisType.Field(i)
			fv := thisValue.Field(i)
			jf := strings.Split(ft.Tag.Get("json"), ",")[0] // json field name

			switch {
			case ft.Type.Kind() == reflect.Ptr && ft.Type.Elem().Kind() == reflect.Struct:
				if fv.IsNil() {
					newVal := reflect.New(ft.Type.Elem())
					fv.Set(newVal)
					if jf == "" && ft.Anonymous {
						UnmarshalMap(input, newVal.Interface(), overrides)
					} else if next, ok := input[jf].(map[string]interface{}); ok {
						UnmarshalMap(next, newVal.Interface(), overrides)
					}
				}

			case ft.Type.Kind() == reflect.Struct:
				if jf == "" && ft.Anonymous {
					UnmarshalMap(input, fv.Addr().Interface(), overrides)
				} else if next, ok := input[jf].(map[string]interface{}); ok {
					UnmarshalMap(next, fv.Addr().Interface(), overrides)
				}

			default:
				if jf != "" {
					value := input[jf]
					k := path.Join(thisType.PkgPath(), thisType.Name(), ft.Name)
					if override, has := overrides[k]; has {
						value = override(value)
					} else if IsValueASet(fv) {

						// Set<X> implemented as map[X]bool

						if reflect.TypeOf(value).Kind() == reflect.Slice {

							vv := reflect.ValueOf(value)
							mc := reflect.MakeMap(fv.Type())

							for i := 0; i < vv.Len(); i++ {
								k := vv.Index(i).Convert(fv.Type().Key())
								mc.SetMapIndex(k, reflect.ValueOf(true))
							}
							value = mc.Interface()
						}
					}
					fv.Set(reflect.ValueOf(value))
				}
			}
		}
	}
}
