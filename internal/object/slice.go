package object

import (
	"errors"
	"fmt"
	"reflect"
)

type Slice struct {
	reflect.Value
}

func MakeSlice(t reflect.Type, len, cap int) (Slice, error) {
	if t.Kind() != reflect.Slice {
		return Slice{}, errors.New("not a slice")
	}
	if t.Elem().Kind() != reflect.Struct {
		return Slice{}, errors.New("not a slice of structs")
	}
	return Slice{Value: reflect.MakeSlice(t, len, cap)}, nil
}

func MustMakeSlice(t reflect.Type, len, cap int) Slice {
	sl, err := MakeSlice(t, len, cap)
	if err != nil {
		panic(err)
	}
	return sl
}

func (sl Slice) Decode(objects ...Values) int {
	n := sl.Value.Len()
	if n > len(objects) {
		n = len(objects)
	}
	for i := 0; i < n; i++ {
		newStruct(sl.Index(i)).Decode(objects[i])
	}
	return n
}

func (sl Slice) StrictDecode(objects ...Values) (int, error) {
	n := sl.Value.Len()
	if n > len(objects) {
		n = len(objects)
	}
	for i := 0; i < n; i++ {
		err := newStruct(sl.Index(i)).StrictDecode(objects[i])
		if err != nil {
			return 0, fmt.Errorf("at %d: %w", i, err)
		}
	}
	return n, nil
}
