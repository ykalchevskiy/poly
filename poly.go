// Package poly allows polymorphic serialization.
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
		return nil, fmt.Errorf("poly: cannot marshal value %v: %v", p.Value, err)
	}

	if bytes.Equal(implData, []byte("null")) {
		return implData, nil
	}

	tnValue, ok := reflect.ValueOf(p.Value).Interface().(TypeName)
	if !ok {
		return nil, fmt.Errorf("poly: cannot get TypeName of %T to marshal", p.Value)
	}

	typeName := tnValue.TypeName()

	var found bool
	var t T
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

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf(`{"type":"%s",`, typeName))
	buf.Write(implData[1:])

	return buf.Bytes(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for Poly.
// It unmarshals the JSON based on the 'type' discriminator field to the correct concrete type.
func (p *Poly[I, T]) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		return nil
	}

	var typeName string

	if reflectValue := reflect.ValueOf(p.Value); reflectValue.IsValid() {
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

	if err := json.Unmarshal(b, &discriminator); err != nil {
		return fmt.Errorf("poly: cannot unmarshal discriminator 'type': %v", err)
	}

	if discriminator.TypeName == "" {
		return errors.New("poly: missing discriminator 'type'")
	}

	var t T
	for _, typ := range t.Types() {
		if typ.Name == discriminator.TypeName {
			if typeName == "" || typeName != typ.Name {
				ptr := reflect.New(typ.ReflectType)
				if err := json.Unmarshal(b, ptr.Interface()); err != nil {
					return fmt.Errorf("poly: cannot unmarshal '%s': %v", typ.ReflectType, err)
				}
				value, ok := ptr.Elem().Interface().(I)
				if !ok {
					return fmt.Errorf("poly: cannon use '%v' as I", ptr.Interface())
				}
				p.Value = value
				return nil
			} else {
				if reflect.ValueOf(p.Value).Kind() == reflect.Pointer {
					// if there is a pointer to a struct, we can use it directly
					if err := json.Unmarshal(b, p.Value); err != nil {
						return fmt.Errorf("poly: cannot unmarshal '%s': %v", typ.ReflectType, err)
					}
					return nil
				} else {
					// otherwise we should create a pointer and copy the existing value there
					ptr := reflect.New(typ.ReflectType)
					ptr.Elem().Set(reflect.ValueOf(p.Value))
					if err := json.Unmarshal(b, ptr.Interface()); err != nil {
						return fmt.Errorf("poly: cannot unmarshal '%s': %v", typ.ReflectType, err)
					}
					value, ok := ptr.Elem().Interface().(I)
					if !ok {
						return fmt.Errorf("poly: cannon use '%v' as I", ptr.Interface())
					}
					p.Value = value
					return nil
				}
			}
		}
	}

	return fmt.Errorf("poly: unknown TypeName %s to unmarshal", discriminator.TypeName)
}
