package ask

import (
	"testing"
)

func BenchmarkForSimplePath(b *testing.B) {
	source := map[string]interface{}{
		"key": "value",
	}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = For(source, "key")
	}
}

func BenchmarkForNestedPath(b *testing.B) {
	source := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "value",
			},
		},
	}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = For(source, "a.b.c")
	}
}

func BenchmarkForSliceIndexing(b *testing.B) {
	source := map[string]interface{}{
		"list": []interface{}{1, 2, 3, 4, 5},
	}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = For(source, "list[3]")
	}
}

func BenchmarkForComplexPath(b *testing.B) {
	source := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{
				"b": []interface{}{
					map[string]interface{}{
						"c": "value",
					},
				},
			},
		},
	}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = For(source, "a[0].b[0].c")
	}
}

func BenchmarkForNonExistingPath(b *testing.B) {
	source := map[string]interface{}{
		"a": map[string]interface{}{
			"b": "value",
		},
	}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = For(source, "a.b.c.d")
	}
}

func BenchmarkStringRetrieval(b *testing.B) {
	source := map[string]interface{}{
		"key": "value",
	}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		For(source, "key").String("")
	}
}

func BenchmarkIntRetrieval(b *testing.B) {
	source := map[string]interface{}{
		"key": 12345,
	}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		For(source, "key").Int(0)
	}
}

func BenchmarkBoolRetrieval(b *testing.B) {
	source := map[string]interface{}{
		"key": true,
	}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		For(source, "key").Bool(false)
	}
}

func BenchmarkSliceRetrieval(b *testing.B) {
	source := map[string]interface{}{
		"list": []interface{}{1, 2, 3, 4, 5},
	}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		For(source, "list").Slice(nil)
	}
}

func BenchmarkMapRetrieval(b *testing.B) {
	source := map[string]interface{}{
		"map": map[string]interface{}{
			"key": "value",
		},
	}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		For(source, "map").Map(nil)
	}
}

func BenchmarkExistsCheck(b *testing.B) {
	source := map[string]interface{}{
		"key": "value",
	}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		For(source, "key").Exists()
	}
}

func BenchmarkValueRetrieval(b *testing.B) {
	source := map[string]interface{}{
		"key": "value",
	}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		For(source, "key").Value()
	}
}

func BenchmarkTokenizePath(b *testing.B) {
	path := "a[0].b[1].c.d[2]"

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = tokenizePath(path)
	}
}

func BenchmarkAccessMap(b *testing.B) {
	source := map[string]interface{}{
		"key": "value",
	}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = accessMap(source, "key")
	}
}

func BenchmarkAccessSlice(b *testing.B) {
	source := []interface{}{1, 2, 3, 4, 5}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = accessSlice(source, 3)
	}
}

func BenchmarkForWithCache(b *testing.B) {
	source := map[string]interface{}{
		"key": "value",
	}
	path := "key"
	// Warm up the cache
	For(source, path)

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = For(source, path)
	}
}

func BenchmarkForWithoutCache(b *testing.B) {
	source := map[string]interface{}{
		"key": "value",
	}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = For(source, "key")
	}
}

func BenchmarkForConcurrent(b *testing.B) {
	source := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "value",
			},
		},
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = For(source, "a.b.c")
		}
	})
}

func BenchmarkForDifferentPaths(b *testing.B) {
	source := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "value1",
				"d": "value2",
				"e": "value3",
			},
		},
	}

	paths := []string{"a.b.c", "a.b.d", "a.b.e"}
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		path := paths[n%len(paths)]
		_ = For(source, path)
	}
}

func BenchmarkForLongPath(b *testing.B) {
	source := map[string]interface{}{
		"level1": map[string]interface{}{
			"level2": map[string]interface{}{
				"level3": map[string]interface{}{
					"level4": map[string]interface{}{
						"level5": "deep_value",
					},
				},
			},
		},
	}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = For(source, "level1.level2.level3.level4.level5")
	}
}

func BenchmarkForInvalidPath(b *testing.B) {
	source := map[string]interface{}{
		"key": "value",
	}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = For(source, "invalid[")
	}
}
