package mineng

import (
	"reflect"

	"github.com/jncornett/mineng/internal/cache"
)

type dataEntry struct {
	cache cache.Cache
	typ   reflect.Type
	arg   reflect.Value
}

type systemEntry struct {
	def        systemDef
	args       []*dataEntry
	cachedArgs []reflect.Value
}

func newSystemEntry(def systemDef) *systemEntry {
	var args []*dataEntry
	for _, t := range def.inputs {
		args = append(args, &dataEntry{typ: t})
	}
	return &systemEntry{
		def:        def,
		args:       args,
		cachedArgs: make([]reflect.Value, len(def.inputs)),
	}
}

func (entry *systemEntry) next() *systemEntry {
	if entry.def.once {
		return nil
	}
	return entry
}

type systemDef struct {
	fn      CallableValue
	inputs  []reflect.Type
	outputs []reflect.Type
	once    bool
}

func newSystemDef(fn interface{}, once bool) systemDef {
	rv := reflect.ValueOf(fn)
	rt := rv.Type()
	var inputs []reflect.Type
	for i := 0; i < rt.NumIn(); i++ {
		inputs = append(inputs, rt.In(i))
	}
	var outputs []reflect.Type
	for i := 0; i < rt.NumOut(); i++ {
		outputs = append(outputs, rt.Out(i))
	}
	return systemDef{
		fn:      rv,
		inputs:  inputs,
		outputs: outputs,
		once:    once,
	}
}

type CallableValue interface {
	Type() reflect.Type
	Call(in []reflect.Value) []reflect.Value
}
