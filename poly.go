// Package poly allows polymorphic encoding and decoding of JSON in Go.
package poly

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// TypeName is an interface that types must implement to provide their unique name for polymorphic serialization.
type TypeName interface {
	TypeName() string
}

// Type holds the name and reflect.Type of a registered polymorphic type.
type Type struct {
	Name        string
	ReflectType reflect.Type
}

// NewType creates a new Type instance for a given TypeName.
func NewType[T TypeName]() Type {
	var t T

	return Type{
		Name:        t.TypeName(),
		ReflectType: reflect.TypeOf(t),
	}
}

// Types is an interface that provides a list of all registered polymorphic types.
type Types interface {
	Types() []Type
}

// Poly is a generic struct that wraps an interface and handles polymorphic JSON marshaling and unmarshaling.
// I is the interface type that the concrete types implement.
// T is a type that implements the Types interface, providing the list of known concrete types.
type Poly[I any, T Types] struct {
	Value I
}

// MarshalJSON implements the json.Marshaler interface for Poly.
// It marshals the underlying value along with its TypeName as a discriminator.
func (p Poly[I, T]) MarshalJSON() ([]byte, error) {
	implData, err := json.Marshal(p.Value)
	if err != nil {
		return nil, fmt.Errorf("poly: cannot marshal value %v: %w", p.Value, err)
	}

	if bytes.Equal(implData, []byte("null")) {
		return implData, nil
	}

	tnValue, ok := reflect.ValueOf(p.Value).Interface().(TypeName)
	if !ok {
		return nil, fmt.Errorf("poly: cannot get TypeName of %T to marshal", p.Value)
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
		return nil, fmt.Errorf("poly: unknown TypeName %s of %T to marshal", typeName, p.Value)
	}

	if bytes.Equal(implData, []byte("{}")) {
		return []byte(fmt.Sprintf(`{"type":"%s"}`, typeName)), nil
	}

	if len(implData) == 0 || implData[0] != '{' {
		return nil, fmt.Errorf("poly: expected JSON object for %T, got %s", p.Value, implData)
	}

	var buf bytes.Buffer

	buf.Grow(len(`{"type":"",`) + len(typeName) + len(implData) - 1)
	fmt.Fprintf(&buf, `{"type":"%s",`, typeName)
	buf.Write(implData[1:])

	return buf.Bytes(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for Poly.
// It unmarshals the JSON based on the 'type' discriminator field to the correct concrete type.
func (p *Poly[I, T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		var zero I

		p.Value = zero

		return nil
	}

	var typeName string

	reflectValue := reflect.ValueOf(p.Value)

	if reflectValue.IsValid() {
		if tnValue, ok := reflectValue.Interface().(TypeName); ok {
			typeName = tnValue.TypeName()
		} else {
			return fmt.Errorf("poly: cannot get TypeName of %T to unmarshal", p.Value)
		}
	}

	discriminator := struct {
		TypeName string `json:"type"`
	}{
		TypeName: typeName,
	}

	if err := json.Unmarshal(data, &discriminator); err != nil {
		return fmt.Errorf("poly: cannot unmarshal discriminator 'type': %w", err)
	}

	if discriminator.TypeName == "" {
		return errors.New("poly: missing discriminator 'type'")
	}

	var t T
	for _, typ := range t.Types() {
		if typ.Name != discriminator.TypeName {
			continue
		}

		// if there was no value yet or it's a new type, we create a new value
		if typeName != typ.Name {
			value, err := unmarshalNew(data, typ, false, p.Value)
			if err != nil {
				return err
			}

			p.Value = value

			return nil
		}

		// if there is a non-nil pointer to a struct, we can use it directly
		if reflectValue.Kind() == reflect.Pointer && !reflectValue.IsNil() {
			if err := json.Unmarshal(data, p.Value); err != nil {
				return fmt.Errorf("poly: cannot unmarshal '%s': %w", typ.ReflectType, err)
			}

			return nil
		}

		// otherwise we should create a pointer and copy the existing value there
		value, err := unmarshalNew(data, typ, true, p.Value)
		if err != nil {
			return err
		}

		p.Value = value

		return nil
	}

	return fmt.Errorf("poly: unknown TypeName %s to unmarshal", discriminator.TypeName)
}

func unmarshalNew[I any](data []byte, typ Type, useCurrent bool, current I) (I, error) {
	ptr := reflect.New(typ.ReflectType)

	if useCurrent {
		ptr.Elem().Set(reflect.ValueOf(current))
	}

	if err := json.Unmarshal(data, ptr.Interface()); err != nil {
		return current, fmt.Errorf("poly: cannot unmarshal '%s': %w", typ.ReflectType, err)
	}

	value, ok := ptr.Elem().Interface().(I)
	if !ok {
		return current, fmt.Errorf("poly: cannot use '%v' as I", ptr.Interface())
	}

	return value, nil
}
