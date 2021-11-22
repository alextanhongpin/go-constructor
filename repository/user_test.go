package repository

import (
	"encoding/json"
	"fmt"
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

func TestUserConstructor(t *testing.T) {
	u1 := User{}
	faker.FakeData(&u1)
	u2 := NewUser(u1)

	if err := validateStructFieldsMatch(u1, u2); err != nil {
		t.Error(err)
	}
}
