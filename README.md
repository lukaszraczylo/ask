## golang Ask - the json query library

Ask provides a simple way of accessing nested properties in maps and slices. Works great in combination with encoding/json and other packages that "Unmarshal" arbitrary data into Go data-types. Inspired by the get function in the lodash javascript library.

This was originally a fork of github.com/simonnilsson/ask but I went a bit overboard with optimisations.

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