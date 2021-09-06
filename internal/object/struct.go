package object

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type Struct struct {
	reflect.Value
	StructType
}

func NewStruct(rv reflect.Value) (Struct, error) {
	rv = reflect.Indirect(rv)
	if rv.Kind() != reflect.Struct {
		return Struct{}, errors.New("not a struct")
	}
	return newStruct(rv), nil
}

func newStruct(rv reflect.Value) Struct {
	return Struct{Value: rv, StructType: StructType{Type: rv.Type()}}
}

func MustNewStruct(rv reflect.Value) Struct {
	st, err := NewStruct(rv)
	if err != nil {
		panic(err)
	}
	return st
}

func ParseStruct(v interface{}) (Struct, error) {
	return NewStruct(reflect.ValueOf(v))
}

func MustParseStruct(v interface{}) Struct {
	st, err := ParseStruct(v)
	if err != nil {
		panic(err)
	}
	return st
}

func (st Struct) Decode(vs Values) {
	for _, sf := range st.StructType.Fields() {
		v, ok := vs[sf.Type]
		if !ok {
			continue
		}
		st.FieldByIndex(sf.Index).Set(v)
	}
}

func (st Struct) StrictDecode(vs Values) error {
	for _, sf := range st.StructType.Fields() {
		v, ok := vs[sf.Type]
		if !ok {
			return fmt.Errorf("not found: %v", sf.Type)
		}
		st.FieldByIndex(sf.Index).Set(v)
	}
	return nil
}

func (st Struct) Values() Values {
	out := Values{}
	for _, sf := range st.StructType.Fields() {
		out[sf.Type] = st.Value.FieldByIndex(sf.Index)
	}
	return out
}

type StructType struct {
	reflect.Type
}

func (st StructType) Fields() []reflect.StructField {
	fields, err := fields(st.Type)
	if err != nil {
		panic(err)
	}
	return fields
}

func FieldTypes(t reflect.Type) []reflect.Type {
	return StructType{Type: t}.FieldTypes()
}

func (st StructType) FieldTypes() []reflect.Type {
	types, err := fieldTypes(st.Type)
	if err != nil {
		panic(err)
	}
	return types
}

func fieldTypes(t reflect.Type) ([]reflect.Type, error) {
	if v, ok := fieldTypeCache.Load(t); ok {
		return v.([]reflect.Type), nil
	}
	fields, err := fields(t)
	if err != nil {
		return nil, err
	}
	var out []reflect.Type
	for _, f := range fields {
		out = append(out, f.Type)
	}
	fieldTypeCache.Store(t, out)
	return out, nil
}

func fields(t reflect.Type) ([]reflect.StructField, error) {
	if v, ok := fieldCache.Load(t); ok {
		return v.([]reflect.StructField), nil
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a type of kind %v, got %v", reflect.Struct, t)
	}
	fields := reflect.VisibleFields(t)
	fieldCache.Store(t, fields)
	return fields, nil
}

var (
	fieldCache     sync.Map
	fieldTypeCache sync.Map
)
