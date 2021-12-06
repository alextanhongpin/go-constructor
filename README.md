# go-constructor

Why constructors is not that simple.


## How to ensure struct fields are all set

Constructor
- use private fields, and allow setting them in constructor only
- use type alias for constructor args


Struct
- validate all fields to ensure they are set
- for boolean fields, use pointer, and validate they are not nil
- use unit test to fake all input data, and check if they are mapped correctly
- use reflection to check fields are mapped (example below)


```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"reflect"
	"strings"
)

type User struct {
	Name    string `map:"name"`
	age     int64  `map:"age"`
	Married *bool  `map:"married"`
}

type DBUser struct {
	Name    string `map:"name"`
	Age     int64  `map:"age"`
	Married *bool  `map:"married"`
}

func ToDBUser(u User) *DBUser {
	return &DBUser{
		Name:    u.Name,
		Age:     u.age,
		Married: u.Married,
	}
}

func main() {
	u := User{Name: "John", age: 1, Married: nil}
	dbu := ToDBUser(u)
	mustMap(u, dbu)
}

func getTagMap(s interface{}) map[string]bool {
	m := make(map[string]bool)
	v := reflect.Indirect(reflect.ValueOf(s))
	n := v.NumField()
	for i := 0; i < n; i++ {
		t := v.Type().Field(i)
		f := v.Field(i)
		tag := t.Tag.Get("map")
		m[tag] = !f.IsZero()
	}
	return m
}

func mustMap(a, b interface{}) {
	lhsName := reflect.Indirect(reflect.ValueOf(a)).Type().Name()
	rhsName := reflect.Indirect(reflect.ValueOf(b)).Type().Name()
	lhs := getTagMap(a)
	rhs := getTagMap(b)
	if len(lhs) > len(rhs) {
		var fields []string
		for k := range lhs {
			_, ok := rhs[k]
			if !ok {
				fields = append(fields, k)
			}
		}
		panic(
			fmt.Errorf("mapErr: missing fields %s<%s>",
				rhsName, strings.Join(fields, ", "),
			),
		)
	} else if len(lhs) == len(rhs) {
		var lhsFields, rhsFields []string
		for k, v := range lhs {
			if !v {
				lhsFields = append(lhsFields, k)
			}
			if !rhs[k] {
				rhsFields = append(rhsFields, k)
			}
		}
		if len(lhsFields) == 0 && len(rhsFields) == 0 {
			return
		}
		if len(lhsFields) == 0 {
			panic(
				fmt.Errorf("mapErr: fields not set fields for %s<%s>",
					rhsName, strings.Join(rhsFields, ", "),
				),
			)
		} else if len(rhsFields) == 0 {
			panic(
				fmt.Errorf("mapErr: fields not set fields for %s<%s>",
					lhsName, strings.Join(lhsFields, ", "),
				),
			)
		} else {
			panic(
				fmt.Errorf("mapErr: fields not set fields for %s<%s> and %s<%s>",
					lhsName, strings.Join(lhsFields, ", "),
					rhsName, strings.Join(rhsFields, ", "),
				),
			)
		}
	} else {
		mustMap(b, a)
	}
}
```
