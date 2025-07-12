package poly_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/ykalchevskiy/poly"
)

type IsItemValue interface {
	IsItemValue()
}

type ItemValue1 struct{}

func (ItemValue1) IsItemValue() {}

func (ItemValue1) TypeName() string {
	return "item-value-1"
}

type ItemValue2 struct {
	Key  string `json:"key,omitempty"`
	Key2 string `json:"key2,omitempty"`
}

func (ItemValue2) IsItemValue() {}

func (ItemValue2) TypeName() string {
	return "item-value-2"
}

type ItemValueUnknown struct{}

func (ItemValueUnknown) IsItemValue() {}

func (ItemValueUnknown) TypeName() string {
	return "item-value-unknown"
}

type ItemValueUnimplementedTypeName struct{}

func (ItemValueUnimplementedTypeName) IsItemValue() {}

type ItemValueUnimplementedIs struct{}

func (ItemValueUnimplementedIs) TypeName() string {
	return "item-value-unimplemented-is"
}

type ItemValue = poly.Poly[IsItemValue, poly.Types2[ItemValue1, ItemValue2]]

type ItemValueInner struct {
	ItemV ItemValue
	ItemP *ItemValue
}

type IsItemPointer interface {
	IsItemPointer()
}

type ItemPointer1 struct{}

func (*ItemPointer1) IsItemPointer() {}

func (*ItemPointer1) TypeName() string {
	return "item-pointer-1"
}

type ItemPointer2 struct {
	Key  string `json:"key,omitempty"`
	Key2 string `json:"key2,omitempty"`
}

func (*ItemPointer2) IsItemPointer() {}

func (*ItemPointer2) TypeName() string {
	return "item-pointer-2"
}

type ItemPointerUnknown struct{}

func (*ItemPointerUnknown) IsItemPointer() {}

func (*ItemPointerUnknown) TypeName() string {
	return "item-pointer-unknown"
}

type ItemPointerUnimplementedTypeName struct{}

func (*ItemPointerUnimplementedTypeName) IsItemPointer() {}

type ItemPointerUnimplementedIs struct{}

func (*ItemPointerUnimplementedIs) TypeName() string {
	return "item-pointer-unimplemented-is"
}

type ItemPointer = poly.Poly[IsItemPointer, poly.Types2[*ItemPointer1, *ItemPointer2]]

type ItemPointerInner struct {
	ItemV ItemPointer
	ItemP *ItemPointer
}

func TestPoly_MarshalJSON_null(t *testing.T) {
	t.Run("value", func(t *testing.T) {
		var item ItemValue

		b, err := json.Marshal(item)
		if err != nil {
			t.Fatalf("marshaling error: %v", err)
		}

		if !bytes.Equal([]byte("null"), b) {
			t.Fatalf("expected null, got %s", b)
		}
	})

	t.Run("value with null pointer", func(t *testing.T) {
		item := ItemValue{Value: (*ItemValue1)(nil)}

		b, err := json.Marshal(item)
		if err != nil {
			t.Fatalf("marshaling error: %v", err)
		}

		if !bytes.Equal([]byte("null"), b) {
			t.Fatalf("expected null, got %s", b)
		}
	})

	t.Run("pointer", func(t *testing.T) {
		var item ItemPointer

		b, err := json.Marshal(item)
		if err != nil {
			t.Fatalf("marshaling error: %v", err)
		}

		if !bytes.Equal([]byte("null"), b) {
			t.Fatalf("expected null, got %s", b)
		}
	})

	t.Run("pointer with null pointer", func(t *testing.T) {
		item := ItemPointer{Value: (*ItemPointer1)(nil)}

		b, err := json.Marshal(item)
		if err != nil {
			t.Fatalf("marshaling error: %v", err)
		}

		if !bytes.Equal([]byte("null"), b) {
			t.Fatalf("expected null, got %s", b)
		}
	})
}

func TestPoly_UnmarshalJSON_null(t *testing.T) {
	t.Run("value", func(t *testing.T) {
		var item ItemValue

		if err := json.Unmarshal([]byte("null"), &item); err != nil {
			t.Fatalf("unmarshaling error: %v", err)
		}

		if item.Value != nil {
			t.Fatalf("expected nil, got %v", item.Value)
		}
	})

	t.Run("pointer", func(t *testing.T) {
		var item ItemPointer

		if err := json.Unmarshal([]byte("null"), &item); err != nil {
			t.Fatalf("unmarshaling error: %v", err)
		}

		if item.Value != nil {
			t.Fatalf("expected nil, got %v", item.Value)
		}
	})
}

func TestPoly_UnmarshalJSON_discriminator(t *testing.T) {
	t.Run("value", func(t *testing.T) {
		var item ItemValue

		err := json.Unmarshal([]byte("{}"), &item)
		if err == nil || !strings.Contains(err.Error(), "missing discriminator") {
			t.Fatalf("expected missing discriminator error, got %v", err)
		}
	})

	t.Run("pointer", func(t *testing.T) {
		var item ItemPointer

		err := json.Unmarshal([]byte("{}"), &item)
		if err == nil || !strings.Contains(err.Error(), "missing discriminator") {
			t.Fatalf("expected missing discriminator error, got %v", err)
		}
	})
}

func TestPoly_ItemValue(t *testing.T) {
	t.Run("value1", func(t *testing.T) {
		var item ItemValue
		bIn := []byte(`{"type":"item-value-1"}`)

		if err := json.Unmarshal(bIn, &item); err != nil {
			t.Fatalf("unmarshaling error: %v", err)
		}

		if got, ok := item.Value.(ItemValue1); !ok {
			t.Fatalf("expected ItemValue1, got %T", got)
		}

		bOut, err := json.Marshal(item)
		if err != nil {
			t.Fatalf("marshaling error: %v", err)
		}

		if !bytes.Equal(bIn, bOut) {
			t.Fatalf("expected %s, got %s", bIn, bOut)
		}
	})

	t.Run("value2", func(t *testing.T) {
		var item ItemValue
		bIn := []byte(`{"type":"item-value-2","key":"k"}`)

		err := json.Unmarshal(bIn, &item)
		if err != nil {
			t.Fatalf("unmarshaling error: %v", err)
		}

		if got, ok := item.Value.(ItemValue2); !ok {
			t.Fatalf("expected ItemValue1, got %T", got)
		}

		bOut, err := json.Marshal(item)
		if err != nil {
			t.Fatalf("marshaling error: %v", err)
		}

		if !bytes.Equal(bIn, bOut) {
			t.Fatalf("expected %s, got %s", bIn, bOut)
		}

		if err := json.Unmarshal([]byte(`{"key":"k2"}`), &item); err != nil {
			t.Fatalf("unmarshaling patch error: %v", err)
		}

		bPatchOut, err := json.Marshal(item)
		if err != nil {
			t.Fatalf("marshaling patch error: %v", err)
		}

		if bIn := []byte(`{"type":"item-value-2","key":"k2"}`); !bytes.Equal(bIn, bPatchOut) {
			t.Fatalf("expected %s, got %s", bIn, bPatchOut)
		}
	})

	t.Run("unknown unmarshal", func(t *testing.T) {
		var item ItemValue

		err := json.Unmarshal([]byte(`{"type":"item-value-unknown"}`), &item)
		if err == nil {
			t.Fatal("unmarshaling unknown should err")
		}
		if !strings.Contains(
			err.Error(),
			"poly: unknown TypeName item-value-unknown to unmarshal",
		) {
			t.Fatalf("expected unknown TypeName error, got %v", err)
		}
	})

	t.Run("unknown marshal", func(t *testing.T) {
		item := ItemValue{Value: ItemValueUnknown{}}

		_, err := json.Marshal(item)
		if err == nil {
			t.Fatal("unmarshaling unknown should err")
		}
		if !strings.Contains(err.Error(), "poly: unknown TypeName item-value-unknown") {
			t.Fatalf("expected unknown TypeName error, got %v", err)
		}
	})

	t.Run("unimplemented unmarshal (is)", func(t *testing.T) {
		var item ItemValue

		err := json.Unmarshal([]byte(`{"type":"item-value-unimplemented-is"}`), &item)
		if err == nil {
			t.Fatal("unmarshaling unimplemented should err")
		}
		if !strings.Contains(
			err.Error(),
			"poly: unknown TypeName item-value-unimplemented-is to unmarshal",
		) {
			t.Fatalf("expected unknown TypeName error, got %v", err)
		}
	})

	t.Run("unimplemented unmarshal (type)", func(t *testing.T) {
		item := ItemValue{Value: ItemValueUnimplementedTypeName{}}

		err := json.Unmarshal([]byte(`{"type":"item-value-1"}`), &item)
		if err == nil {
			t.Fatal("unmarshaling unknown should err")
		}
		if !strings.Contains(
			err.Error(),
			"poly: cannot get TypeName of poly_test.ItemValueUnimplementedTypeName to unmarshal",
		) {
			t.Fatalf("expected unknown type error, got %v", err)
		}
	})

	t.Run("unimplemented marshal", func(t *testing.T) {
		item := ItemValue{Value: ItemValueUnimplementedTypeName{}}

		_, err := json.Marshal(item)
		if err == nil {
			t.Fatal("unmarshaling unknown should err")
		}
		if !strings.Contains(
			err.Error(),
			"poly: cannot get TypeName of poly_test.ItemValueUnimplementedTypeName to marshal",
		) {
			t.Fatalf("expected unknown type error, got %v", err)
		}
	})

	t.Run("type change", func(t *testing.T) {
		item := ItemValue{Value: ItemValue1{}}

		err := json.Unmarshal([]byte(`{"type":"item-value-2","key":"k"}`), &item)
		if err != nil {
			t.Fatalf("unmarshaling error: %v", err)
		}

		if got, ok := item.Value.(ItemValue2); !ok {
			t.Fatalf("expected ItemValue1, got %T", got)
		}
	})

	t.Run("settability", func(t *testing.T) {
		item := ItemValue{Value: ItemValue2{}}

		rValue := reflect.ValueOf(item)
		if rValue.Field(0).Elem().Field(0).CanSet() {
			t.Fatalf("value cannot be settable")
		}
	})

	t.Run("inner", func(t *testing.T) {
		bIn := []byte(
			`{"ItemV":{"type":"item-value-1"},"ItemP":{"type":"item-value-2","key2":"k2"}}`,
		)

		var inner ItemValueInner

		if err := json.Unmarshal(bIn, &inner); err != nil {
			t.Fatalf("unmarshaling error: %v", err)
		}

		if got, ok := inner.ItemV.Value.(ItemValue1); !ok {
			t.Fatalf("expected ItemValue1, got %T", got)
		}

		if got, ok := inner.ItemP.Value.(ItemValue2); !ok {
			t.Fatalf("expected ItemValue2, got %T", got)
		}

		bOut, err := json.Marshal(inner)
		if err != nil {
			t.Fatalf("marshaling error: %v", err)
		}

		if !bytes.Equal(bIn, bOut) {
			t.Fatalf("expected %s, got %s", bIn, bOut)
		}
	})
}

func TestPoly_ItemPointer(t *testing.T) {
	t.Run("pointer1", func(t *testing.T) {
		var item ItemPointer
		bIn := []byte(`{"type":"item-pointer-1"}`)

		if err := json.Unmarshal(bIn, &item); err != nil {
			t.Fatalf("unmarshaling error: %v", err)
		}

		if got, ok := item.Value.(*ItemPointer1); !ok {
			t.Fatalf("expected ItemPointer1, got %T", got)
		}

		bOut, err := json.Marshal(item)
		if err != nil {
			t.Fatalf("marshaling error: %v", err)
		}

		if !bytes.Equal(bIn, bOut) {
			t.Fatalf("expected %s, got %s", bIn, bOut)
		}
	})

	t.Run("pointer2", func(t *testing.T) {
		var item ItemPointer
		bIn := []byte(`{"type":"item-pointer-2","key":"k"}`)

		err := json.Unmarshal(bIn, &item)
		if err != nil {
			t.Fatalf("unmarshaling error: %v", err)
		}

		if got, ok := item.Value.(*ItemPointer2); !ok {
			t.Fatalf("expected ItemPointer1, got %T", got)
		}

		bOut, err := json.Marshal(item)
		if err != nil {
			t.Fatalf("marshaling error: %v", err)
		}

		if !bytes.Equal(bIn, bOut) {
			t.Fatalf("expected %s, got %s", bIn, bOut)
		}

		if err := json.Unmarshal([]byte(`{"key":"k2"}`), &item); err != nil {
			t.Fatalf("unmarshaling patch error: %v", err)
		}

		bPatchOut, err := json.Marshal(item)
		if err != nil {
			t.Fatalf("marshaling patch error: %v", err)
		}

		if bIn := []byte(`{"type":"item-pointer-2","key":"k2"}`); !bytes.Equal(bIn, bPatchOut) {
			t.Fatalf("expected %s, got %s", bIn, bPatchOut)
		}
	})

	t.Run("unknown unmarshal", func(t *testing.T) {
		var item ItemPointer

		err := json.Unmarshal([]byte(`{"type":"item-pointer-unknown"}`), &item)
		if err == nil {
			t.Fatal("unmarshaling unknown should err")
		}
		if !strings.Contains(
			err.Error(),
			"poly: unknown TypeName item-pointer-unknown to unmarshal",
		) {
			t.Fatalf("expected unknown TypeName error, got %v", err)
		}
	})

	t.Run("unknown marshal", func(t *testing.T) {
		item := ItemPointer{Value: &ItemPointerUnknown{}}

		_, err := json.Marshal(item)
		if err == nil {
			t.Fatal("unmarshaling unknown should err")
		}
		if !strings.Contains(err.Error(), "poly: unknown TypeName item-pointer-unknown") {
			t.Fatalf("expected unknown TypeName error, got %v", err)
		}
	})

	t.Run("unimplemented unmarshal (is)", func(t *testing.T) {
		var item ItemPointer

		err := json.Unmarshal([]byte(`{"type":"item-pointer-unimplemented-is"}`), &item)
		if err == nil {
			t.Fatal("unmarshaling unimplemented should err")
		}
		if !strings.Contains(
			err.Error(),
			"poly: unknown TypeName item-pointer-unimplemented-is to unmarshal",
		) {
			t.Fatalf("expected unknown TypeName error, got %v", err)
		}
	})

	t.Run("unimplemented unmarshal (type)", func(t *testing.T) {
		item := ItemPointer{Value: &ItemPointerUnimplementedTypeName{}}

		err := json.Unmarshal([]byte(`{"type":"item-pointer-1"}`), &item)
		if err == nil {
			t.Fatal("unmarshaling unknown should err")
		}
		if !strings.Contains(
			err.Error(),
			"poly: cannot get TypeName of *poly_test.ItemPointerUnimplementedTypeName to unmarshal",
		) {
			t.Fatalf("expected unknown type error, got %v", err)
		}
	})

	t.Run("unimplemented marshal", func(t *testing.T) {
		item := ItemPointer{Value: &ItemPointerUnimplementedTypeName{}}

		_, err := json.Marshal(item)
		if err == nil {
			t.Fatal("unmarshaling unknown should err")
		}
		if !strings.Contains(
			err.Error(),
			"poly: cannot get TypeName of *poly_test.ItemPointerUnimplementedTypeName to marshal",
		) {
			t.Fatalf("expected unknown type error, got %v", err)
		}
	})

	t.Run("type change", func(t *testing.T) {
		item := ItemPointer{Value: &ItemPointer1{}}

		err := json.Unmarshal([]byte(`{"type":"item-pointer-2","key":"k"}`), &item)
		if err != nil {
			t.Fatalf("unmarshaling error: %v", err)
		}

		if got, ok := item.Value.(*ItemPointer2); !ok {
			t.Fatalf("expected ItemPointer1, got %T", got)
		}
	})

	t.Run("settability", func(t *testing.T) {
		item := ItemPointer{Value: &ItemPointer2{}}

		rValue := reflect.ValueOf(item)
		if !rValue.Field(0).Elem().Elem().Field(0).CanSet() {
			t.Fatalf("pointer must be settable")
		}
	})

	t.Run("inner", func(t *testing.T) {
		bIn := []byte(
			`{"ItemV":{"type":"item-pointer-1"},"ItemP":{"type":"item-pointer-2","key2":"k2"}}`,
		)

		var inner ItemPointerInner

		if err := json.Unmarshal(bIn, &inner); err != nil {
			t.Fatalf("unmarshaling error: %v", err)
		}

		if got, ok := inner.ItemV.Value.(*ItemPointer1); !ok {
			t.Fatalf("expected ItemPointer1, got %T", got)
		}

		if got, ok := inner.ItemP.Value.(*ItemPointer2); !ok {
			t.Fatalf("expected ItemPointer2, got %T", got)
		}

		bOut, err := json.Marshal(inner)
		if err != nil {
			t.Fatalf("marshaling error: %v", err)
		}

		if !bytes.Equal(bIn, bOut) {
			t.Fatalf("expected %s, got %s", bIn, bOut)
		}
	})
}
