//go:build !(go1.25 && goexperiment.jsonv2)

package poly_test

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestPoly_JSONOptionsAreNotReusedInV1(t *testing.T) {
	bIn := []byte(`{"type":"item-value-2","key":"k","nonexistingkey":"v"}`)

	t.Run("non-strict", func(t *testing.T) {
		var item ItemValue

		if err := json.Unmarshal(bIn, &item); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("strict", func(t *testing.T) {
		// in v1 the options are not reused and the unknown field is not rejected
		var item ItemValue

		dec := json.NewDecoder(bytes.NewReader(bIn))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&item); err != nil {
			t.Errorf(
				"unexpected error: %v",
				err,
			)
		}
	})
}
