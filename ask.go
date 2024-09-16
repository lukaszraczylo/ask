package ask

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var splitCache sync.Map // concurrent-safe map

// Answer holds result of call to For, use one of its methods to extract a value.
type Answer struct {
	value any
}

// For is used to select a path from source to return as answer.
func For(source any, path string) *Answer {
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
		if strings.HasPrefix(token, "[") && strings.HasSuffix(token, "]") {
			// Handle array index
			indexStr := strings.TrimSpace(token[1 : len(token)-1])
			if index, err := strconv.Atoi(indexStr); err == nil {
				current = accessSlice(current, index)
			} else {
				return &Answer{}
			}
		} else {
			// Handle map key
			current = accessMap(current, token)
		}
		if current == nil {
			return &Answer{}
		}
	}

	return &Answer{value: current}
}

func accessMap(source any, key string) any {
	switch m := source.(type) {
	case map[string]any:
		return m[key]
	case map[string]string:
		return m[key]
	case map[string]int:
		return m[key]
	}
	// Use reflect as last resort
	val := reflect.ValueOf(source)
	if val.Kind() == reflect.Map {
		keyVal := reflect.ValueOf(key)
		valueVal := val.MapIndex(keyVal)
		if valueVal.IsValid() {
			return valueVal.Interface()
		}
	}
	return nil
}

func accessSlice(source any, index int) any {
	val := reflect.ValueOf(source)
	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		if index >= 0 && index < val.Len() {
			return val.Index(index).Interface()
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

// Value returns the raw value as type any, can be nil if no value is available.
func (a *Answer) Value() any {
	return a.value
}

// String attempts to retrieve the answer as a string.
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
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(v).Int(), true
	case uint, uint8, uint16, uint32, uint64:
		uv := reflect.ValueOf(v).Uint()
		if uv <= uint64(^uint64(0)>>1) {
			return int64(uv), true
		}
	case float32, float64:
		return int64(reflect.ValueOf(v).Float()), true
	}
	return def, false
}

// Uint attempts to retrieve the answer as uint64.
func (a *Answer) Uint(def uint64) (uint64, bool) {
	if a.value == nil {
		return def, false
	}
	switch v := a.value.(type) {
	case int, int8, int16, int32, int64:
		iv := reflect.ValueOf(v).Int()
		if iv >= 0 {
			return uint64(iv), true
		}
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(v).Uint(), true
	case float32, float64:
		fv := reflect.ValueOf(v).Float()
		if fv >= 0 {
			return uint64(fv), true
		}
	}
	return def, false
}

// Float attempts to retrieve the answer as float64.
// Float attempts to retrieve the answer as float64.
func (a *Answer) Float(def float64) (float64, bool) {
	if a.value == nil {
		return def, false
	}
	switch v := a.value.(type) {
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(v).Int()), true
	case uint, uint8, uint16, uint32, uint64:
		return float64(reflect.ValueOf(v).Uint()), true
	case float32, float64:
		return reflect.ValueOf(v).Float(), true
	}
	return def, false
}

// Slice attempts to retrieve the answer as []any.
func (a *Answer) Slice(def []any) ([]any, bool) {
	if a.value == nil {
		return def, false
	}
	if s, ok := a.value.([]any); ok {
		return s, true
	}
	val := reflect.ValueOf(a.value)
	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		length := val.Len()
		result := make([]any, length)
		for i := 0; i < length; i++ {
			result[i] = val.Index(i).Interface()
		}
		return result, true
	}
	return def, false
}

// Map attempts to retrieve the answer as map[string]any.
func (a *Answer) Map(def map[string]any) (map[string]any, bool) {
	if a.value == nil {
		return def, false
	}
	if m, ok := a.value.(map[string]any); ok {
		return m, true
	}
	val := reflect.ValueOf(a.value)
	if val.Kind() == reflect.Map {
		result := make(map[string]any)
		iter := val.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			if key, ok := k.Interface().(string); ok {
				result[key] = v.Interface()
			}
		}
		return result, true
	}
	return def, false
}

func tokenizePath(path string) []string {
	tokens := make([]string, 0, 8) // Pre-allocate for small paths
	var token strings.Builder
	token.Grow(len(path)) // Pre-allocate builder capacity
	inBracket := false

	for i := 0; i < len(path); i++ {
		ch := path[i]
		switch {
		case ch <= ' ':
			if inBracket {
				token.WriteByte(ch)
			}
		case ch == '.':
			if inBracket {
				token.WriteByte(ch)
			} else if token.Len() > 0 {
				tokens = append(tokens, trimSpaceASCII(token.String()))
				token.Reset()
			}
		case ch == '[':
			if token.Len() > 0 {
				tokens = append(tokens, trimSpaceASCII(token.String()))
				token.Reset()
			}
			token.WriteByte(ch)
			inBracket = true
		case ch == ']':
			token.WriteByte(ch)
			if inBracket {
				tokens = append(tokens, trimSpaceASCII(token.String()))
				token.Reset()
				inBracket = false
			}
		default:
			token.WriteByte(ch)
		}
	}

	if token.Len() > 0 {
		tokens = append(tokens, trimSpaceASCII(token.String()))
	}

	return tokens
}

// trimSpaceASCII is a faster version of strings.TrimSpace for ASCII strings
func trimSpaceASCII(s string) string {
	start := 0
	for ; start < len(s); start++ {
		if s[start] > ' ' {
			break
		}
	}
	end := len(s)
	for ; end > start; end-- {
		if s[end-1] > ' ' {
			break
		}
	}
	if start == 0 && end == len(s) {
		return s
	}
	return s[start:end]
}
