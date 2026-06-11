//go:build go1.25 && goexperiment.jsonv2

package poly

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"encoding/json/jsontext"
	"encoding/json/v2"
)

// MarshalJSONTo implements the json.MarshalerTo interface for Poly.
// It marshals the underlying value along with its TypeName as a discriminator.
func (p Poly[I, T]) MarshalJSONTo(enc *jsontext.Encoder) error {
	implData, err := json.Marshal(p.Value, enc.Options())
	if err != nil {
		return fmt.Errorf("poly: cannot marshal value %v: %w", p.Value, err)
	}

	if bytes.Equal(implData, []byte("null")) {
		return enc.WriteValue(implData)
	}

	tnValue, ok := reflect.ValueOf(p.Value).Interface().(TypeName)
	if !ok {
		return fmt.Errorf("poly: cannot get TypeName of %T to marshal", p.Value)
	}

	typeName := tnValue.TypeName()

	var (
		t     T
		found bool
	)

	for _, typ := range t.Types() {
		if typeName == typ.Name {
			found = true

			break
		}
	}

	if !found {
		return fmt.Errorf("poly: unknown TypeName %s of %T to marshal", typeName, p.Value)
	}

	if bytes.Equal(implData, []byte("{}")) {
		return enc.WriteValue([]byte(fmt.Sprintf(`{"type":"%s"}`, typeName)))
	}

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf(`{"type":"%s",`, typeName))
	buf.Write(implData[1:])

	return enc.WriteValue(buf.Bytes())
}

// UnmarshalJSONFrom implements the json.UnmarshalerFrom interface for Poly.
// It unmarshals the JSON based on the 'type' discriminator field to the correct concrete type.
func (p *Poly[I, T]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	var typeName string

	if reflectValue := reflect.ValueOf(p.Value); reflectValue.IsValid() {
		if tnValue, ok := reflectValue.Interface().(TypeName); ok {
			typeName = tnValue.TypeName()
		} else {
			return fmt.Errorf("poly: cannot get TypeName of %T to unmarshal", p.Value)
		}
	}

	fullData := &struct {
		TypeName string         `json:"type"`
		Data     jsontext.Value `json:",unknown"`
	}{
		TypeName: typeName,
		Data:     []byte("{}"), // default to empty object to avoid unmarshaling into nil
	}

	if err := json.UnmarshalDecode(dec, &fullData, json.RejectUnknownMembers(false)); err != nil {
		return fmt.Errorf("poly: cannot unmarshal discriminator 'type': %w", err)
	}

	if fullData == nil {
		return nil
	}

	if fullData.TypeName == "" {
		return errors.New("poly: missing discriminator 'type'")
	}

	var t T
	for _, typ := range t.Types() {
		if typ.Name != fullData.TypeName {
			continue
		}

		// if there was no value yet or it's a new type, we create a new value
		if typeName != typ.Name {
			value, err := unmarshalV2New(fullData.Data, typ, false, p.Value, dec.Options())
			if err != nil {
				return err
			}

			p.Value = value

			return nil
		}

		// if there is a pointer to a struct, we can use it directly
		if reflect.ValueOf(p.Value).Kind() == reflect.Pointer {
			if err := json.Unmarshal(fullData.Data, p.Value, dec.Options()); err != nil {
				return fmt.Errorf("poly: cannot unmarshal '%s': %w", typ.ReflectType, err)
			}

			return nil
		}

		// otherwise we should create a pointer and copy the existing value there
		value, err := unmarshalV2New(fullData.Data, typ, true, p.Value, dec.Options())
		if err != nil {
			return err
		}

		p.Value = value

		return nil
	}

	return fmt.Errorf("poly: unknown TypeName %s to unmarshal", fullData.TypeName)
}

func unmarshalV2New[I any](
	data []byte,
	typ Type,
	useCurrent bool,
	current I,
	opts json.Options,
) (I, error) {
	ptr := reflect.New(typ.ReflectType)

	if useCurrent {
		ptr.Elem().Set(reflect.ValueOf(current))
	}

	if err := json.Unmarshal(data, ptr.Interface(), opts); err != nil {
		return current, fmt.Errorf("poly: cannot unmarshal '%s': %w", typ.ReflectType, err)
	}

	value, ok := ptr.Elem().Interface().(I)
	if !ok {
		return current, fmt.Errorf("poly: cannot use '%v' as I", ptr.Interface())
	}

	return value, nil
}
