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
	"errors"
	"fmt"
	"reflect"
	"sort"
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
	t := false
	u := User{Name: "John", age: 1, Married: &t}
	dbu := ToDBUser(u)
	panic(mustMap(u, dbu))
}

func getStructTagFields(s interface{}, tag string) (string, map[string]bool) {
	m := make(map[string]bool)
	v := reflect.Indirect(reflect.ValueOf(s))
	n := v.NumField()
	for i := 0; i < n; i++ {
		t := v.Type().Field(i)
		f := v.Field(i)
		k := t.Tag.Get(tag)
		m[k] = !f.IsZero()
	}
	return v.Type().Name(), m
}

func mustMap(a, b interface{}) error {
	lhsName, lhs := getStructTagFields(a, "map")
	rhsName, rhs := getStructTagFields(b, "map")

	all := make(map[string]bool)
	for k := range lhs {
		all[k] = true
	}
	for k := range rhs {
		all[k] = true
	}
	if len(all) == 0 {
		return errors.New("mapErr: no fields to map")
	}

	var lhsFields, rhsFields []string
	for k := range all {
		if !rhs[k] {
			lhsFields = append(lhsFields, k)
		}
		if !lhs[k] {
			rhsFields = append(rhsFields, k)
		}
	}
	if len(lhsFields) == 0 && len(rhsFields) == 0 {
		return nil
	}

	sort.Strings(lhsFields)
	sort.Strings(rhsFields)

	if len(lhsFields) == 0 {
		return fmt.Errorf("mapErr: fields not mapped for %s<%s>",
			rhsName, strings.Join(rhsFields, ", "),
		)
	} else if len(rhsFields) == 0 {
		return fmt.Errorf("mapErr: fields not mapped for %s<%s>",
			lhsName, strings.Join(lhsFields, ", "),
		)
	}

	return fmt.Errorf("mapErr: fields not mapped for %s<%s> and %s<%s>",
		lhsName, strings.Join(lhsFields, ", "),
		rhsName, strings.Join(rhsFields, ", "),
	)
}
```
