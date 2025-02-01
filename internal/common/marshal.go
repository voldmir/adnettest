package common

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type fieldParameters struct {
	optional     bool // true, если поле НЕОБЯЗАТЕЛЬНО
	omitEmpty    bool // true, если это значение следует опустить, если оно пусто при маршалинге.
	defaultValue *int64
	bigEndian    bool
}

type StructuralError struct {
	Msg string
}

func (e StructuralError) Error() string { return "structure error: " + e.Msg }

func canHaveDefaultValue(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	}

	return false
}

func parseFieldParameters(str string) (ret fieldParameters) {

	for _, part := range strings.Split(str, ",") {
		switch {
		case part == "optional":
			ret.optional = true

		case part == "omitempty":
			ret.omitEmpty = true
		case part == "be":
			ret.bigEndian = true

		case strings.HasPrefix(part, "default:"):
			i, err := strconv.ParseInt(part[8:], 10, 64)
			if err == nil {
				ret.defaultValue = new(int64)
				*ret.defaultValue = i
			}
		}
	}
	return
}

func marshalString(out *bytes.Buffer, s string) (err error) {
	_, err = out.WriteString(s)
	return
}

func MarshalNumberLE(out *bytes.Buffer, i reflect.Value) (err error) {
	n := uint8(i.Type().Bits() / 8)
	for ; n > 0; n-- {
		err = out.WriteByte(byte(i.Convert(reflect.TypeOf(uint(0))).Uint() >> uint((n-1)*8)))
		if err != nil {
			return
		}
	}
	return nil
}

func MarshalNumberBE(out *bytes.Buffer, i reflect.Value) (err error) {
	var n uint8
	for ; n < uint8(i.Type().Bits()/8); n++ {
		err = out.WriteByte(byte(i.Convert(reflect.TypeOf(uint(0))).Uint() >> uint((n)*8)))
		if err != nil {
			return
		}
	}
	return nil
}

func marshalBody(out *bytes.Buffer, value reflect.Value, params fieldParameters) (err error) {
	v := value
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if params.bigEndian {
			return MarshalNumberBE(out, v)
		} else {
			return MarshalNumberLE(out, v)
		}

	case reflect.Struct:
		t := v.Type()
		startingField := 0
		for i := startingField; i < t.NumField(); i++ {
			err = marshalField(out, v.Field(i), parseFieldParameters(t.Field(i).Tag.Get("tag")))
			if err != nil {
				return
			}
		}
		return

	case reflect.Slice, reflect.Array:
		sliceType := v.Type()
		if sliceType.Elem().Kind() == reflect.Uint8 {
			bytes := make([]byte, v.Len())
			for i := 0; i < v.Len(); i++ {
				bytes[i] = uint8(v.Index(i).Uint())
			}
			_, err = out.Write(bytes)
			return
		}

		for i := 0; i < v.Len(); i++ {
			err = marshalField(out, v.Index(i), params)
			if err != nil {
				return
			}
		}
		return

	case reflect.String:
		return marshalString(out, v.String())
	}

	return StructuralError{fmt.Sprintf("unknown Go type: %v %+v", v.Type(), v.Kind())}
}

func marshalField(out *bytes.Buffer, v reflect.Value, params fieldParameters) (err error) {
	if !v.IsValid() {
		return fmt.Errorf("cannot marshal nil value")
	}
	if v.Kind() == reflect.Interface && v.Type().NumMethod() == 0 {
		return marshalField(out, v.Elem(), params)
	}

	if (v.Kind() == reflect.Slice || v.Kind() == reflect.Array) && v.Len() == 0 && params.omitEmpty {
		return
	}

	if params.optional && params.defaultValue != nil && canHaveDefaultValue(v.Kind()) {
		defaultValue := reflect.New(v.Type()).Elem()
		defaultValue.SetInt(*params.defaultValue)

		if reflect.DeepEqual(v.Interface(), defaultValue.Interface()) {
			return
		}
	}

	if params.optional && params.defaultValue == nil {
		if reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface()) {
			return
		}
	}

	err = marshalBody(out, v, params)
	if err != nil {
		return
	}

	return err
}

func Marshal(val interface{}) ([]byte, error) {
	out := bytes.Buffer{}
	v := reflect.ValueOf(val)
	err := marshalField(&out, v, fieldParameters{})
	if err != nil {
		return nil, err
	}
	return out.Bytes(), err
}
