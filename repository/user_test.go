package repository

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func validateStructFieldsMatch(src, tgt interface{}) error {
	sb, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		return err
	}

	tb, err := json.MarshalIndent(tgt, "", "  ")
	if err != nil {
		return err
	}

	splitStringsByNewline := func(b []byte) []string {
		return strings.Split(string(b), "\n")
	}
	opts := cmpopts.AcyclicTransformer("multiline", splitStringsByNewline)
	if diff := cmp.Diff(sb, tb, opts); diff != "" {
		return fmt.Errorf("want-, got+:\n%s", diff)
	}
	return nil
}

type AggregateError struct {
	Message string
	Fields  map[string]string
}

func NewAggregateError(message string) *AggregateError {
	return &AggregateError{
		Message: message,
		Fields:  make(map[string]string),
	}
}

func (e *AggregateError) Add(name, reason string) {
	e.Fields[name] = reason
}

func (e *AggregateError) Error() string {
	var b strings.Builder
	b.WriteString(e.Message + ": ")
	for k, v := range e.Fields {
		b.WriteString(fmt.Sprintf("%s %s\n", k, v))
	}
	return b.String()
}

// This solution might not be suitable for handling nested structs or slices.
func validateStructFieldsSet(s interface{}) error {
	v := reflect.Indirect(reflect.ValueOf(s))
	err := NewAggregateError(v.Type().Name() + "Error")
	for i := 0; i < v.Type().NumField(); i++ {
		t := v.Type().Field(i)
		if !t.IsExported() {
			continue
		}
		f := v.Field(i)
		if f.IsZero() {
			err.Add(t.Name, "not set")
		}
	}
	return err
}

func TestUserConstructor(t *testing.T) {
	u1 := User{}
	faker.FakeData(&u1)
	u2 := NewUser(u1)

	if err := validateStructFieldsMatch(u1, u2); err != nil {
		t.Error(err)
	}
}

func TestUserConstructorSet(t *testing.T) {
	u1 := User{}
	faker.FakeData(&u1)
	u2 := NewUser(u1)

	if err := validateStructFieldsSet(u2); err != nil {
		t.Error(err)
	}
}
