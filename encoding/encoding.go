/*
Package encoding defines types and functions used for dealing with stringified-JSON.
*/
package encoding

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// String is a string type to be used when the default json marshaler/unmarshaler cannot be avoided
type String string

func NewString(ss string) *String {
	s := String(ss)
	return &s
}

func (s *String) Value() *string {
	return (*string)(s)
}

func (s String) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

func (s *String) UnmarshalJSON(data []byte) error {
	var ss string
	err := json.Unmarshal(data, &ss)
	if err != nil {
		return err
	}

	*s = String(ss)
	return nil
}

// Bool is a bool type to be used when the default json marshaler/unmarshaler cannot be avoided
type Bool bool

func NewBool(bb bool) *Bool {
	b := Bool(bb)
	return &b
}

func (b *Bool) Value() *bool {
	return (*bool)(b)
}

func (b Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprint(bool(b)))
}

func (b *Bool) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	val, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}

	*b = Bool(val)
	return nil
}

// Int is an int type to be used when the default json marshaler/unmarshaler cannot be avoided
type Int int64

func NewInt(ii int64) *Int {
	i := Int(ii)
	return &i
}

func (i *Int) Value() *int64 {
	return (*int64)(i)
}

func (i Int) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprint(int64(i)))
}

func (i *Int) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	val, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}

	*i = Int(val)
	return nil
}

// Float is a float type to be used when the default json marshaler/unmarshaler cannot be avoided
type Float float64

func NewFloat(ff float64) *Float {
	f := Float(ff)
	return &f
}

func (f *Float) Value() *float64 {
	return (*float64)(f)
}

func (f Float) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprint(float64(f)))
}

func (f *Float) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*f = Float(val)
	return nil
}
