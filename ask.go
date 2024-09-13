// Package ask provides a simple way of accessing nested properties in maps and arrays.
// Works great in combination with encoding/json and other packages that "Unmarshal" arbitrary data into Go data-types.
// Inspired by the get function in the lodash javascript library.
package ask

import (
	"reflect"
	"strconv"
	"sync"
)

var splitCache sync.Map // concurrent-safe map

// Answer holds result of call to For, use one of its methods to extract a value.
type Answer struct {
	value interface{}
}

// For is used to select a path from source to return as answer.
func For(source interface{}, path string) *Answer {
	partsInterface, ok := splitCache.Load(path)
	var parts []string
	if ok {
		parts = partsInterface.([]string)
	} else {
		parts = tokenizePath(path)
		splitCache.Store(path, parts)
	}

	current := source

	for _, token := range parts {
		if index, err := strconv.Atoi(token); err == nil {
			current = accessSlice(current, index)
		} else {
			current = accessMap(current, token)
		}
		if current == nil {
			return &Answer{}
		}
	}

	return &Answer{value: current}
}

func accessMap(source interface{}, key string) interface{} {
	switch m := source.(type) {
	case map[string]interface{}:
		return m[key]
	case map[interface{}]interface{}:
		return m[key]
	case map[string]string:
		return m[key]
	case map[string]int:
		return m[key]
	// Add more cases as needed
	default:
		// Use reflect as last resort
		val := reflect.ValueOf(source)
		if val.Kind() == reflect.Map {
			keyVal := reflect.ValueOf(key)
			valueVal := val.MapIndex(keyVal)
			if valueVal.IsValid() {
				return valueVal.Interface()
			}
		}
	}
	return nil
}

func accessSlice(source interface{}, index int) interface{} {
	switch s := source.(type) {
	case []interface{}:
		if index >= 0 && index < len(s) {
			return s[index]
		}
	case []int:
		if index >= 0 && index < len(s) {
			return s[index]
		}
	case [][]int:
		if index >= 0 && index < len(s) {
			return s[index]
		}
	default:
		val := reflect.ValueOf(source)
		if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
			if index >= 0 && index < val.Len() {
				return val.Index(index).Interface()
			}
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

// String attempts to retrieve the answer as a string.
// The first return value is the result, and the second indicates if the operation was successful.
// If not successful, the first return value will be set to the default value provided.
func (a *Answer) String(def string) (string, bool) {
	if a.value == nil {
		return def, false
	}
	if res, ok := a.value.(string); ok {
		return res, true
	}
	return def, false
}

// Bool attempts to retrieve the answer as a bool.
// The first return value is the result, and the second indicates if the operation was successful.
// If not successful, the first return value will be set to the default value provided.
func (a *Answer) Bool(def bool) (bool, bool) {
	if a.value == nil {
		return def, false
	}
	if res, ok := a.value.(bool); ok {
		return res, true
	}
	return def, false
}

// Int attempts to retrieve the answer as int64.
func (a *Answer) Int(def int64) (int64, bool) {
	if a.value == nil {
		return def, false
	}
	switch v := a.value.(type) {
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case uint:
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		if v <= uint64(^uint64(0)>>1) {
			return int64(v), true
		}
	case float32:
		return int64(v), true
	case float64:
		return int64(v), true
	default:
		val := reflect.ValueOf(a.value)
		if val.Kind() >= reflect.Int && val.Kind() <= reflect.Int64 {
			return val.Int(), true
		}
		if val.Kind() >= reflect.Uint && val.Kind() <= reflect.Uint64 {
			u := val.Uint()
			if u <= uint64(^uint64(0)>>1) {
				return int64(u), true
			}
		}
		if val.Kind() == reflect.Float32 || val.Kind() == reflect.Float64 {
			return int64(val.Float()), true
		}
	}
	return def, false
}

// Uint attempts to retrieve the answer as uint64.
func (a *Answer) Uint(def uint64) (uint64, bool) {
	if a.value == nil {
		return def, false
	}
	switch v := a.value.(type) {
	case int:
		if v >= 0 {
			return uint64(v), true
		}
	case int8:
		if v >= 0 {
			return uint64(v), true
		}
	case int16:
		if v >= 0 {
			return uint64(v), true
		}
	case int32:
		if v >= 0 {
			return uint64(v), true
		}
	case int64:
		if v >= 0 {
			return uint64(v), true
		}
	case uint:
		return uint64(v), true
	case uint8:
		return uint64(v), true
	case uint16:
		return uint64(v), true
	case uint32:
		return uint64(v), true
	case uint64:
		return v, true
	case float32:
		if v >= 0 {
			return uint64(v), true
		}
	case float64:
		if v >= 0 {
			return uint64(v), true
		}
	default:
		val := reflect.ValueOf(a.value)
		if val.Kind() >= reflect.Int && val.Kind() <= reflect.Int64 {
			i := val.Int()
			if i >= 0 {
				return uint64(i), true
			}
		}
		if val.Kind() >= reflect.Uint && val.Kind() <= reflect.Uint64 {
			return val.Uint(), true
		}
		if val.Kind() == reflect.Float32 || val.Kind() == reflect.Float64 {
			f := val.Float()
			if f >= 0 {
				return uint64(f), true
			}
		}
	}
	return def, false
}

// Float attempts to retrieve the answer as float64.
func (a *Answer) Float(def float64) (float64, bool) {
	if a.value == nil {
		return def, false
	}
	switch v := a.value.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		val := reflect.ValueOf(a.value)
		if val.Kind() >= reflect.Int && val.Kind() <= reflect.Int64 {
			return float64(val.Int()), true
		}
		if val.Kind() >= reflect.Uint && val.Kind() <= reflect.Uint64 {
			return float64(val.Uint()), true
		}
		if val.Kind() == reflect.Float32 || val.Kind() == reflect.Float64 {
			return val.Float(), true
		}
	}
	return def, false
}

// Slice attempts to retrieve the answer as []interface{}.
func (a *Answer) Slice(def []interface{}) ([]interface{}, bool) {
	if a.value == nil {
		return def, false
	}
	if s, ok := a.value.([]interface{}); ok {
		return s, true
	}
	val := reflect.ValueOf(a.value)
	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		length := val.Len()
		result := make([]interface{}, length)
		for i := 0; i < length; i++ {
			result[i] = val.Index(i).Interface()
		}
		return result, true
	}
	return def, false
}

// Map attempts to retrieve the answer as map[string]interface{}.
func (a *Answer) Map(def map[string]interface{}) (map[string]interface{}, bool) {
	if a.value == nil {
		return def, false
	}
	if m, ok := a.value.(map[string]interface{}); ok {
		return m, true
	}
	val := reflect.ValueOf(a.value)
	if val.Kind() == reflect.Map {
		result := make(map[string]interface{})
		iter := val.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			var key string
			if k.Kind() == reflect.String {
				key = k.String()
			} else if k.CanInterface() {
				if ks, ok := k.Interface().(string); ok {
					key = ks
				} else {
					continue // skip non-string keys
				}
			} else {
				continue
			}
			result[key] = v.Interface()
		}
		return result, true
	}
	return def, false
}

func tokenizePath(path string) []string {
	tokens := make([]string, 0, 8) // Preallocate with expected capacity
	i := 0
	n := len(path)
	for i < n {
		switch path[i] {
		case '.':
			i++
		case '[':
			i++ // skip '['
			// Skip leading spaces
			for i < n && path[i] == ' ' {
				i++
			}
			start := i
			for i < n && path[i] != ']' {
				i++
			}
			end := i
			// Trim trailing spaces
			for end > start && path[end-1] == ' ' {
				end--
			}
			if start < end {
				token := path[start:end]
				tokens = append(tokens, token)
			}
			if i < n && path[i] == ']' {
				i++ // skip ']'
			}
		default:
			start := i
			for i < n && path[i] != '.' && path[i] != '[' {
				i++
			}
			end := i
			// Trim trailing spaces
			for end > start && path[end-1] == ' ' {
				end--
			}
			// Trim leading spaces
			for start < end && path[start] == ' ' {
				start++
			}
			if start < end {
				token := path[start:end]
				tokens = append(tokens, token)
			}
		}
	}
	return tokens
}
