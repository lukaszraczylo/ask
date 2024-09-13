package ask

import (
	"math"
	"reflect"
	"testing"
)

func TestFor(t *testing.T) {
	source := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{
				"b": 100,
				"c": -1031337,
			},
		},
	}

	tests := []struct {
		name   string
		source interface{}
		path   string
		want   interface{}
	}{
		{
			name:   "Valid path to positive integer",
			source: source,
			path:   "a[0].b",
			want:   100,
		},
		{
			name:   "Valid path to negative integer",
			source: source,
			path:   "a[0].c",
			want:   -1031337,
		},
		{
			name:   "Missing index in slice",
			source: source,
			path:   "a[1].b",
			want:   nil,
		},
		{
			name:   "Missing key in map",
			source: source,
			path:   "d[1]",
			want:   nil,
		},
		{
			name:   "Empty path returns source",
			source: source,
			path:   "",
			want:   source,
		},
		{
			name:   "Invalid path returns nil",
			source: source,
			path:   "---",
			want:   nil,
		},
		{
			name:   "Whitespace in path is trimmed",
			source: source,
			path:   " a[0] . b ",
			want:   100,
		},
		{
			name: "Access nested slice",
			source: map[string]interface{}{
				"list": [][]int{
					{1, 2, 3},
					{4, 5, 6},
				},
			},
			path: "list[1][2]",
			want: 6,
		},
		{
			name:   "Non-integer index",
			source: source,
			path:   "a[foo]",
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			answer := For(tt.source, tt.path)
			if !reflect.DeepEqual(answer.value, tt.want) {
				t.Errorf("For() = (%v); want (%v)", answer.value, tt.want)
			}
		})
	}
}

func TestPath(t *testing.T) {
	source := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{
				"b": 100,
			},
		},
	}

	tests := []struct {
		name   string
		answer *Answer
		path   string
		want   interface{}
	}{
		{
			name:   "Valid path after For",
			answer: For(source, "a[0]"),
			path:   "b",
			want:   100,
		},
		{
			name:   "Invalid path after For",
			answer: For(source, "a[1]"),
			path:   "b",
			want:   nil,
		},
		{
			name:   "Empty path returns current value",
			answer: For(source, ""),
			path:   "",
			want:   source,
		},
		{
			name:   "Nested Path calls",
			answer: For(source, "a[0]"),
			path:   "b",
			want:   100,
		},
		{
			name:   "Path with whitespace",
			answer: For(source, "a[0]"),
			path:   " b ",
			want:   100,
		},
		{
			name:   "Invalid intermediate path",
			answer: For(source, "invalid"),
			path:   "b",
			want:   nil,
		},
		{
			name:   "Accessing nil Answer",
			answer: &Answer{},
			path:   "b",
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			answer := tt.answer.Path(tt.path)
			if !reflect.DeepEqual(answer.value, tt.want) {
				t.Errorf("Path() = (%v); want (%v)", answer.value, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	source := map[string]interface{}{
		"string":  "test",
		"number":  100,
		"boolean": true,
		"nil":     nil,
	}

	tests := []struct {
		name   string
		path   string
		def    string
		want   string
		wantOK bool
	}{
		{
			name:   "Valid string value",
			path:   "string",
			def:    "default",
			want:   "test",
			wantOK: true,
		},
		{
			name:   "Non-string value",
			path:   "number",
			def:    "default",
			want:   "default",
			wantOK: false,
		},
		{
			name:   "Missing key",
			path:   "missing",
			def:    "default",
			want:   "default",
			wantOK: false,
		},
		{
			name:   "Nil value",
			path:   "nil",
			def:    "default",
			want:   "default",
			wantOK: false,
		},
		{
			name:   "Boolean value",
			path:   "boolean",
			def:    "default",
			want:   "default",
			wantOK: false,
		},
		{
			name:   "Empty string value",
			path:   "empty",
			def:    "default",
			want:   "default",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, ok := For(source, tt.path).String(tt.def)
			if res != tt.want || ok != tt.wantOK {
				t.Errorf(`String() = ("%s", %t); want ("%s", %t)`, res, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestBool(t *testing.T) {
	source := map[string]interface{}{
		"boolTrue":  true,
		"boolFalse": false,
		"number":    100,
		"string":    "test",
		"nil":       nil,
	}

	tests := []struct {
		name   string
		path   string
		def    bool
		want   bool
		wantOK bool
	}{
		{
			name:   "Valid true boolean",
			path:   "boolTrue",
			def:    false,
			want:   true,
			wantOK: true,
		},
		{
			name:   "Valid false boolean",
			path:   "boolFalse",
			def:    true,
			want:   false,
			wantOK: true,
		},
		{
			name:   "Non-boolean value (number)",
			path:   "number",
			def:    false,
			want:   false,
			wantOK: false,
		},
		{
			name:   "Non-boolean value (string)",
			path:   "string",
			def:    false,
			want:   false,
			wantOK: false,
		},
		{
			name:   "Missing key",
			path:   "missing",
			def:    false,
			want:   false,
			wantOK: false,
		},
		{
			name:   "Nil value",
			path:   "nil",
			def:    false,
			want:   false,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, ok := For(source, tt.path).Bool(tt.def)
			if res != tt.want || ok != tt.wantOK {
				t.Errorf(`Bool() = (%t, %t); want (%t, %t)`, res, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestInt(t *testing.T) {
	source := map[string]interface{}{
		"positive":        100,
		"negative":        -1001587205010,
		"unsigned":        uint(100),
		"float":           float64(100.9),
		"negative_float":  float64(-100.9),
		"string_number":   "123",
		"string_nonumber": "abc",
		"nil":             nil,
		"toobig":          uint64(math.MaxUint64),
		"toonegative":     int64(math.MinInt64),
	}

	tests := []struct {
		name   string
		path   string
		def    int64
		want   int64
		wantOK bool
	}{
		{
			name:   "Positive integer",
			path:   "positive",
			def:    5,
			want:   100,
			wantOK: true,
		},
		{
			name:   "Negative integer",
			path:   "negative",
			def:    5,
			want:   -1001587205010,
			wantOK: true,
		},
		{
			name:   "Unsigned integer",
			path:   "unsigned",
			def:    5,
			want:   100,
			wantOK: true,
		},
		{
			name:   "Float value",
			path:   "float",
			def:    5,
			want:   100,
			wantOK: true,
		},
		{
			name:   "Negative float value",
			path:   "negative_float",
			def:    5,
			want:   -100,
			wantOK: true,
		},
		{
			name:   "String representing number",
			path:   "string_number",
			def:    5,
			want:   5,
			wantOK: false,
		},
		{
			name:   "String not representing number",
			path:   "string_nonumber",
			def:    5,
			want:   5,
			wantOK: false,
		},
		{
			name:   "Nil value",
			path:   "nil",
			def:    5,
			want:   5,
			wantOK: false,
		},
		{
			name:   "Too big number",
			path:   "toobig",
			def:    5,
			want:   5,
			wantOK: false,
		},
		{
			name:   "Too negative number",
			path:   "toonegative",
			def:    5,
			want:   math.MinInt64,
			wantOK: true,
		},
		{
			name:   "Missing key",
			path:   "missing",
			def:    5,
			want:   5,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, ok := For(source, tt.path).Int(tt.def)
			if res != tt.want || ok != tt.wantOK {
				t.Errorf("Int() = (%d, %t); want (%d, %t)", res, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestUint(t *testing.T) {
	source := map[string]interface{}{
		"positive":       100,
		"negative":       -100,
		"unsigned":       uint(100),
		"float":          float64(100.9),
		"negative_float": float64(-100.9),
		"string_number":  "123",
		"string_text":    "text",
		"nil":            nil,
	}

	tests := []struct {
		name   string
		path   string
		def    uint64
		want   uint64
		wantOK bool
	}{
		{
			name:   "Positive integer",
			path:   "positive",
			def:    5,
			want:   100,
			wantOK: true,
		},
		{
			name:   "Unsigned integer",
			path:   "unsigned",
			def:    5,
			want:   100,
			wantOK: true,
		},
		{
			name:   "Float value",
			path:   "float",
			def:    5,
			want:   100,
			wantOK: true,
		},
		{
			name:   "Negative integer",
			path:   "negative",
			def:    5,
			want:   5,
			wantOK: false,
		},
		{
			name:   "Negative float value",
			path:   "negative_float",
			def:    5,
			want:   5,
			wantOK: false,
		},
		{
			name:   "String representing number",
			path:   "string_number",
			def:    5,
			want:   5,
			wantOK: false,
		},
		{
			name:   "String not representing number",
			path:   "string_text",
			def:    5,
			want:   5,
			wantOK: false,
		},
		{
			name:   "Nil value",
			path:   "nil",
			def:    5,
			want:   5,
			wantOK: false,
		},
		{
			name:   "Missing key",
			path:   "missing",
			def:    5,
			want:   5,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, ok := For(source, tt.path).Uint(tt.def)
			if res != tt.want || ok != tt.wantOK {
				t.Errorf("Uint() = (%d, %t); want (%d, %t)", res, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestFloat(t *testing.T) {
	source := map[string]interface{}{
		"positive":        100,
		"negative":        -100,
		"unsigned":        uint(100),
		"float32":         float32(100.1),
		"float64":         float64(100.1),
		"string_number":   "123.45",
		"string_nonumber": "abc",
		"nil":             nil,
	}

	tests := []struct {
		name   string
		path   string
		def    float64
		want   float64
		wantOK bool
	}{
		{
			name:   "Positive integer",
			path:   "positive",
			def:    5.0,
			want:   100.0,
			wantOK: true,
		},
		{
			name:   "Negative integer",
			path:   "negative",
			def:    5.0,
			want:   -100.0,
			wantOK: true,
		},
		{
			name:   "Unsigned integer",
			path:   "unsigned",
			def:    5.0,
			want:   100.0,
			wantOK: true,
		},
		{
			name:   "Float32 value",
			path:   "float32",
			def:    5.0,
			want:   100.1,
			wantOK: true,
		},
		{
			name:   "Float64 value",
			path:   "float64",
			def:    5.0,
			want:   100.1,
			wantOK: true,
		},
		{
			name:   "String representing number",
			path:   "string_number",
			def:    5.0,
			want:   5.0,
			wantOK: false,
		},
		{
			name:   "String not representing number",
			path:   "string_nonumber",
			def:    5.0,
			want:   5.0,
			wantOK: false,
		},
		{
			name:   "Nil value",
			path:   "nil",
			def:    5.0,
			want:   5.0,
			wantOK: false,
		},
		{
			name:   "Missing key",
			path:   "missing",
			def:    5.0,
			want:   5.0,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, ok := For(source, tt.path).Float(tt.def)
			if math.Abs(res-tt.want) > 1e-5 || ok != tt.wantOK {
				t.Errorf("Float() = (%f, %t); want (%f, %t)", res, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestSlice(t *testing.T) {
	def := []interface{}{"default"}

	source := map[string]interface{}{
		"slice":   []interface{}{1, 2, 3},
		"empty":   []interface{}{},
		"string":  "test",
		"nil":     nil,
		"invalid": 123,
	}

	tests := []struct {
		name    string
		path    string
		def     []interface{}
		want    []interface{}
		wantLen int
		wantOK  bool
	}{
		{
			name:    "Valid slice",
			path:    "slice",
			def:     def,
			want:    []interface{}{1, 2, 3},
			wantLen: 3,
			wantOK:  true,
		},
		{
			name:    "Empty slice",
			path:    "empty",
			def:     def,
			want:    []interface{}{},
			wantLen: 0,
			wantOK:  true,
		},
		{
			name:    "Non-slice value (string)",
			path:    "string",
			def:     def,
			want:    def,
			wantLen: len(def),
			wantOK:  false,
		},
		{
			name:    "Non-slice value (int)",
			path:    "invalid",
			def:     def,
			want:    def,
			wantLen: len(def),
			wantOK:  false,
		},
		{
			name:    "Nil value",
			path:    "nil",
			def:     def,
			want:    def,
			wantLen: len(def),
			wantOK:  false,
		},
		{
			name:    "Missing key",
			path:    "missing",
			def:     def,
			want:    def,
			wantLen: len(def),
			wantOK:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, ok := For(source, tt.path).Slice(tt.def)
			if len(res) != tt.wantLen || ok != tt.wantOK {
				t.Errorf("Slice() = (len=%d, %t); want (len=%d, %t)", len(res), ok, tt.wantLen, tt.wantOK)
			}
			if !reflect.DeepEqual(res, tt.want) {
				t.Errorf("Slice() content mismatch; got %v, want %v", res, tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	def := map[string]interface{}{"default": "value"}

	source := map[string]interface{}{
		"map":     map[string]interface{}{"key1": "value1", "key2": "value2"},
		"empty":   map[string]interface{}{},
		"string":  "test",
		"nil":     nil,
		"invalid": 123,
	}

	tests := []struct {
		name    string
		path    string
		def     map[string]interface{}
		want    map[string]interface{}
		wantLen int
		wantOK  bool
	}{
		{
			name:    "Valid map",
			path:    "map",
			def:     def,
			want:    map[string]interface{}{"key1": "value1", "key2": "value2"},
			wantLen: 2,
			wantOK:  true,
		},
		{
			name:    "Empty map",
			path:    "empty",
			def:     def,
			want:    map[string]interface{}{},
			wantLen: 0,
			wantOK:  true,
		},
		{
			name:    "Non-map value (string)",
			path:    "string",
			def:     def,
			want:    def,
			wantLen: len(def),
			wantOK:  false,
		},
		{
			name:    "Non-map value (int)",
			path:    "invalid",
			def:     def,
			want:    def,
			wantLen: len(def),
			wantOK:  false,
		},
		{
			name:    "Nil value",
			path:    "nil",
			def:     def,
			want:    def,
			wantLen: len(def),
			wantOK:  false,
		},
		{
			name:    "Missing key",
			path:    "missing",
			def:     def,
			want:    def,
			wantLen: len(def),
			wantOK:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, ok := For(source, tt.path).Map(tt.def)
			if len(res) != tt.wantLen || ok != tt.wantOK {
				t.Errorf("Map() = (len=%d, %t); want (len=%d, %t)", len(res), ok, tt.wantLen, tt.wantOK)
			}
			if !reflect.DeepEqual(res, tt.want) {
				t.Errorf("Map() content mismatch; got %v, want %v", res, tt.want)
			}
		})
	}
}

func TestExists(t *testing.T) {
	source := map[string]interface{}{
		"value1": "test",
		"value2": 0,
		"nil":    nil,
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "Existing string value",
			path: "value1",
			want: true,
		},
		{
			name: "Existing zero integer value",
			path: "value2",
			want: true,
		},
		{
			name: "Nil value",
			path: "nil",
			want: false,
		},
		{
			name: "Missing key",
			path: "missing",
			want: false,
		},
		{
			name: "Existing value in nested map",
			path: "value1.length",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := For(source, tt.path).Exists()
			if res != tt.want {
				t.Errorf("Exists() = (%t); want (%t)", res, tt.want)
			}
		})
	}
}

func TestValue(t *testing.T) {
	source := map[string]interface{}{
		"value1": "test",
		"value2": 0,
		"value3": nil,
	}

	tests := []struct {
		name string
		path string
		want interface{}
	}{
		{
			name: "Existing string value",
			path: "value1",
			want: "test",
		},
		{
			name: "Existing zero integer value",
			path: "value2",
			want: 0,
		},
		{
			name: "Nil value",
			path: "value3",
			want: nil,
		},
		{
			name: "Missing key",
			path: "missing",
			want: nil,
		},
		{
			name: "Nested missing key",
			path: "value1.missing",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := For(source, tt.path).Value()
			if !reflect.DeepEqual(res, tt.want) {
				t.Errorf("Value() = (%v); want (%v)", res, tt.want)
			}
		})
	}
}
