package toml

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

var typeOfStringSlice = reflect.TypeOf([]string(nil))
var typeOfIntSlice = reflect.TypeOf([]int(nil))

func isPointer(v interface{}) bool {
	vkind := reflect.ValueOf(v).Type().Kind()
	return vkind == reflect.Ptr
}

// Same as PrimitiveDecode but adds a strict verification.
func PrimitiveDecodeStrict(primValue Primitive,
	v interface{},
	ignore_fields map[string]interface{}) (err error) {

	// Only accept pointer types.
	if !isPointer(v) {
		return fmt.Errorf("Must use pointer type for strict decoding: [%s]", v)
	}

	err = PrimitiveDecode(primValue, v)
	if err != nil {
		return
	}

	thestruct := reflect.ValueOf(v).Elem().Interface()
	return CheckType(primValue, thestruct, ignore_fields)
}

// The same as Decode, except that parsed data that cannot be mapped will
// throw an error.
func DecodeStrict(data string,
	v interface{},
	ignore_fields map[string]interface{}) (m MetaData, err error) {

	// Only accept pointer types.
	if !isPointer(v) {
		err = fmt.Errorf("Must use pointer type for strict decoding: [%s]", v)
		return
	}

	m, err = Decode(data, v)
	if err != nil {
		return
	}

	thestruct := reflect.ValueOf(v).Elem().Interface()
	err = CheckType(m.mapping, thestruct, ignore_fields)
	return
}

func Contains(list []string, elem string) bool {
	for _, t := range list {
		if t == elem {
			return true
		}
	}
	return false
}

func CheckType(data interface{},
	thestruct interface{},
	ignore_fields map[string]interface{}) (err error) {

	var dType reflect.Type
	var structAsType reflect.Type
	var structAsTypeOk bool
	var structAsValue reflect.Value
	var structAsValueType reflect.Type

	dType = reflect.TypeOf(data)

	structAsType, structAsTypeOk = thestruct.(reflect.Type)

	structAsValue = reflect.ValueOf(thestruct)
	structAsValueType = structAsValue.Type()

	// Special case. Go's `time.Time` is a struct, which we don't want
	// to confuse with a user struct.
	timeType := rvalue(time.Time{}).Type()

	if dType == timeType && thestruct == timeType {
		return nil
	}

	if structAsTypeOk {
		return checkTypeStructAsType(data,
			structAsType,
			ignore_fields)
	} else {
		return checkTypeStructAsType(data,
			structAsValueType,
			ignore_fields)
	}
}

func checkTypeStructAsType(data interface{},
	structAsType reflect.Type,
	ignore_fields map[string]interface{}) (err error) {

	dType := reflect.ValueOf(data).Type()
	dKind := dType.Kind()

	// Handle all the int types
	dIsInt := (dKind >= reflect.Int && dKind <= reflect.Uint64)
	sIsInt := (structAsType.Kind() >= reflect.Int && structAsType.Kind() <= reflect.Uint64)
	if dIsInt && sIsInt {
		return nil
	}

	structKind := structAsType.Kind()
	switch structKind {
	case reflect.Map:
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Expected data to be a map: [%s]", data)
		}

		// Check the elem, which is the type inside the structAsType
		// container
		structMapElem := structAsType.Elem()

		for _, v := range dataMap {
			// Check each of the items in our dataMap against the
			// underlying type of the slice type we are mapping onto
			elemType := structMapElem.(reflect.Type)
			if err = CheckType(v, elemType, ignore_fields); err != nil {
				return err
			}
		}
		return nil
	case reflect.Slice:
		dataSlice := data.([]interface{})
		// Get the underlying type of the slice in the struct
		structSliceElem := structAsType.Elem()
		for _, v := range dataSlice {
			// Check each of the items in our dataslice against the
			// underlying type of the slice type we are mapping onto
			elemType := structSliceElem.(reflect.Type)
			if err = CheckType(v, elemType, ignore_fields); err != nil {
				return err
			}
		}
		return nil
	case reflect.String:
		_, ok := data.(string)
		if ok {
			return nil
		}
		return fmt.Errorf("Incoming type didn't match gotype string")
	case reflect.Bool:
		_, ok := data.(bool)
		if ok {
			return nil
		}
		return fmt.Errorf("Incoming type didn't match gotype bool")
	case reflect.Interface:
		if structAsType.NumMethod() == 0 {
			return nil
		} else {
			return fmt.Errorf("We don't write data to non-empty interfaces around here")
		}
	case reflect.Float32, reflect.Float64:
		var ok bool
		_, ok = data.(float32)
		if ok {
			return nil
		}
		_, ok = data.(float64)
		if ok {
			return nil
		}
		return fmt.Errorf("Incoming type didn't match gotype float32/float64")
	case reflect.Array:
		return fmt.Errorf("*** This shouldn't happen")
	case reflect.Struct:
		dataMap := data.(map[string]interface{})
		// need to iterate over each key in the data to make
		// sure it exists in structAsType
		mapKeys := make([]string, 0)
		for k, _ := range dataMap {
			mapKeys = append(mapKeys, strings.ToLower(k))
		}
		structKeys := make([]string, 0)
		var fieldName string
		for i := 0; i < structAsType.NumField(); i++ {
			f := structAsType.Field(i)

			fieldName = f.Tag.Get("toml")
			if len(fieldName) == 0 {
				fieldName = f.Name
			}
			structKeys = append(structKeys, strings.ToLower(fieldName))
		}

		for _, k := range mapKeys {
			if !Contains(structKeys, k) {
				if _, ok := insensitiveGet(ignore_fields, k); !ok {
					return e("Configuration contains key [%s] "+
						"which doesn't exist in struct", k)
				}
			}
		}

		// Check each struct field against incoming data if
		// available
		for i := 0; i < structAsType.NumField(); i++ {
			f := structAsType.Field(i)
			fieldName := f.Name
			mapdata, ok := insensitiveGet(dataMap, fieldName)
			if ok {
				err = CheckType(mapdata, f.Type, ignore_fields)
				if err != nil {
					return err
				}
			}
		}
		return nil
	default:
		return fmt.Errorf("Unrecognized struct kind: [%s]", structKind)
	}

	return nil
}
