// Package ask provides a simple way of accessing nested properties in maps and arrays.
// Works great in combination with encoding/json and other packages that "Unmarshal" arbitrary data into Go data-types.
// Inspired by the get function in the lodash javascript library.
package ask

import (
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var tokenMatcher = regexp.MustCompile("([^[]+)?(?:\\[(\\d+)])?")
var mapType = reflect.TypeOf(map[string]interface{}{})
var sliceType = reflect.TypeOf([]interface{}{})
var splitCache = make(map[string][]string)

// Answer holds result of call to For, use one of its methods to extract a value.
type Answer struct {
	value interface{}
}

// For is used to select a path from source to return as answer.
func For(source interface{}, path string) *Answer {
	parts, ok := splitCache[path]
	if !ok {
		parts = strings.Split(path, ".")
		splitCache[path] = parts
	}

	current := source

	for _, part := range parts {
		match := tokenMatcher.FindStringSubmatch(strings.TrimSpace(part))
		if len(match) == 3 {
			if match[1] != "" {
				current = accessMap(current, match[1])
				if current == nil {
					return &Answer{}
				}
			}

			if match[2] != "" {
				current = accessSlice(current, match[2])
				if current == nil {
					return &Answer{}
				}
			}
		}
	}

	return &Answer{value: current}
}

func accessMap(source interface{}, key string) interface{} {
	val := reflect.ValueOf(source)
	if val.IsValid() && val.Type().ConvertibleTo(mapType) {
		return val.Convert(mapType).Interface().(map[string]interface{})[key]
	}
	return nil
}

func accessSlice(source interface{}, indexStr string) interface{} {
	val := reflect.ValueOf(source)
	if val.IsValid() && val.Type().ConvertibleTo(sliceType) {
		s := val.Convert(sliceType).Interface().([]interface{})
		index, _ := strconv.Atoi(indexStr)
		if index >= 0 && index < len(s) {
			return s[index]
		}
	}
	return nil
}

// Path does the same thing as For but uses existing answer as source.
func (a *Answer) Path(path string) *Answer {
	return For(a.value, path)
}

// Exists returns a boolean indicating if the answer exists (not nil).
func (a *Answer) Exists() bool {
	return a.value != nil
}

// Value returns the raw value as type interface{}, can be nil if no value is available.
func (a *Answer) Value() interface{} {
	return a.value
}

// Slice attempts asserting answer as a []interface{}.
// The first return value is the result, and the second indicates if the operation was successful.
// If not successful the first return value will be set to the d parameter.
func (a *Answer) Slice(d []interface{}) ([]interface{}, bool) {
	val := reflect.ValueOf(a.value)
	if val.IsValid() && val.CanConvert(sliceType) {
		return val.Convert(sliceType).Interface().([]interface{}), true
	}
	return d, false
}

// Map attempts asserting answer as a map[string]interface{}.
// The first return value is the result, and the second indicates if the operation was successful.
// If not successful the first return value will be set to the d parameter.
func (a *Answer) Map(d map[string]interface{}) (map[string]interface{}, bool) {
	val := reflect.ValueOf(a.value)
	if val.IsValid() && val.CanConvert(mapType) {
		return val.Convert(mapType).Interface().(map[string]interface{}), true
	}
	return d, false
}

// String attempts asserting answer as a string.
// The first return value is the result, and the second indicates if the operation was successful.
// If not successful the first return value will be set to the d parameter.
func (a *Answer) String(d string) (string, bool) {
	res, ok := a.value.(string)
	if ok {
		return res, ok
	}
	return d, false
}

// Bool attempts asserting answer as a bool.
// The first return value is the result, and the second indicates if the operation was successful.
// If not successful the first return value will be set to the d parameter.
func (a *Answer) Bool(d bool) (bool, bool) {
	res, ok := a.value.(bool)
	if ok {
		return res, ok
	}
	return d, false
}

// Int attempts asserting answer as a int64. Casting from other number types will be done if necessary.
// The first return value is the result, and the second indicates if the operation was successful.
// If not successful the first return value will be set to the d parameter.
func (a *Answer) Int(d int64) (int64, bool) {
	switch vt := a.value.(type) {
	case int:
		return int64(vt), true
	case int8:
		return int64(vt), true
	case int16:
		return int64(vt), true
	case int32:
		return int64(vt), true
	case int64:
		return vt, true
	case uint:
		if vt <= uint(math.MaxInt64) {
			return int64(vt), true
		}
	case uint8:
		return int64(vt), true
	case uint16:
		return int64(vt), true
	case uint32:
		return int64(vt), true
	case uint64:
		if vt <= uint64(math.MaxInt64) {
			return int64(vt), true
		}
	case float32:
		if vt >= float32(math.MinInt64) && vt <= float32(math.MaxInt64) {
			return int64(vt), true
		}
	case float64:
		if vt >= float64(math.MinInt64) && vt <= float64(math.MaxInt64) {
			return int64(vt), true
		}
	}
	return d, false
}

// Uint attempts asserting answer as a uint64. Casting from other number types will be done if necessary.
// The first return value is the result, and the second indicates if the operation was successful.
// If not successful the first return value will be set to the d parameter.
func (a *Answer) Uint(d uint64) (uint64, bool) {
	switch vt := a.value.(type) {
	case int:
		if vt >= 0 {
			return uint64(vt), true
		}
	case int8:
		if vt >= 0 {
			return uint64(vt), true
		}
	case int16:
		if vt >= 0 {
			return uint64(vt), true
		}
	case int32:
		if vt >= 0 {
			return uint64(vt), true
		}
	case int64:
		if vt >= 0 {
			return uint64(vt), true
		}
	case uint:
		return uint64(vt), true
	case uint8:
		return uint64(vt), true
	case uint16:
		return uint64(vt), true
	case uint32:
		return uint64(vt), true
	case uint64:
		return vt, true
	case float32:
		if vt >= 0 {
			if vt > float32(math.MaxUint64) {
				return math.MaxUint64, true
			}
			return uint64(vt), true
		}
	case float64:
		if vt >= 0 {
			if vt > float64(math.MaxUint64) {
				return math.MaxUint64, true
			}
			return uint64(vt), true
		}
	}
	return d, false
}

// Float attempts asserting answer as a float64. Casting from other number types will be done if necessary.
// The first return value is the result, and the second indicates if the operation was successful.
// If not successful the first return value will be set to the d parameter.
func (a *Answer) Float(d float64) (float64, bool) {
	switch vt := a.value.(type) {
	case int:
		return float64(vt), true
	case int8:
		return float64(vt), true
	case int16:
		return float64(vt), true
	case int32:
		return float64(vt), true
	case int64:
		return float64(vt), true
	case uint:
		return float64(vt), true
	case uint8:
		return float64(vt), true
	case uint16:
		return float64(vt), true
	case uint32:
		return float64(vt), true
	case uint64:
		return float64(vt), true
	case float32:
		return float64(vt), true
	case float64:
		return vt, true
	}
	return d, false
}
