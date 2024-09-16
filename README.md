## golang Ask - the json query library

Ask provides a simple way of accessing nested properties in maps and slices. Works great in combination with encoding/json and other packages that "Unmarshal" arbitrary data into Go data-types. Inspired by the get function in the lodash javascript library.

This was originally a fork of github.com/simonnilsson/ask but I went a bit overboard with optimisations making it close to no allocation json parsing rocket.

## Use

```go
package main

import "json"
import "github.com/lukaszraczylo/ask"

func main() {

	// Use parsed JSON as source data
	var object map[string]interface{}
	json.Unmarshal([]byte(`{ "a": [{ "b": { "c": 3 } }] }`), &object)

	// Extract the 3
	res, ok := ask.For(object, "a[0].b.c").Int(0)

	fmt.Println(res, ok)
	// Output: 3 true

	// Attempt extracting a string at path .d that does not exist
	res2, ok := ask.For(object, "a[0].b.d").String("nothing")

	fmt.Println(res2, ok)
	// Output: nothing false

}
```

## Benchmarks

```
goos: darwin
goarch: arm64
pkg: github.com/lukaszraczylo/ask
cpu: Apple M1
BenchmarkForSimplePath
BenchmarkForSimplePath-8        	 3520928	       320.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkForNestedPath
BenchmarkForNestedPath-8        	 2763940	       437.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkForSliceIndexing
BenchmarkForSliceIndexing-8     	 1816672	       660.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkForComplexPath
BenchmarkForComplexPath-8       	 1000000	      1185 ns/op	      16 B/op	       1 allocs/op
BenchmarkForNonExistingPath
BenchmarkForNonExistingPath-8   	 2516452	       480.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkStringRetrieval
BenchmarkStringRetrieval-8      	 3494808	       340.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkIntRetrieval
BenchmarkIntRetrieval-8         	 3011366	       399.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkBoolRetrieval
BenchmarkBoolRetrieval-8        	 3562989	       328.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkSliceRetrieval
BenchmarkSliceRetrieval-8       	 3472946	       345.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkMapRetrieval
BenchmarkMapRetrieval-8         	 3512316	       342.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkExistsCheck
BenchmarkExistsCheck-8          	 3733477	       321.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkValueRetrieval
BenchmarkValueRetrieval-8       	 3737104	       343.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkTokenizePath
BenchmarkTokenizePath-8         	  886958	      1372 ns/op	     240 B/op	       8 allocs/op
BenchmarkAccessMap
BenchmarkAccessMap-8            	28041696	        42.21 ns/op	       0 B/op	       0 allocs/op
BenchmarkAccessSlice
BenchmarkAccessSlice-8          	 3647130	       332.5 ns/op	      24 B/op	       1 allocs/op
BenchmarkForWithCache
BenchmarkForWithCache-8         	 3623319	       320.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkForWithoutCache
BenchmarkForWithoutCache-8      	 3778798	       318.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkForConcurrent
BenchmarkForConcurrent-8        	 4503613	       255.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkForDifferentPaths
BenchmarkForDifferentPaths-8    	 2605047	       469.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkForLongPath
BenchmarkForLongPath-8          	 2215105	       542.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkForInvalidPath
BenchmarkForInvalidPath-8       	 3738764	       323.4 ns/op	      16 B/op	       1 allocs/op
```