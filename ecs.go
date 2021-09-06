package mineng

import (
	"reflect"

	"sync"

	"github.com/jncornett/mineng/internal/dep"
	"github.com/jncornett/mineng/internal/object"
)

type ECS struct {
	objects *object.DB
	assets  object.Values
	closure dep.Closure
	pending []systemDef
	systems []*systemEntry
	mu      sync.Mutex
}

func NewECS() *ECS {
	return &ECS{
		objects: object.NewDB(),
		assets:  object.Values{},
		closure: dep.NewClosure(),
	}
}

func (ecs *ECS) Asset(asset interface{}) {
	key := ecs.assets.Encode(asset)
	ecs.closure.Reset(key)
}

func (ecs *ECS) InitSystem(system interface{}) {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()
	ecs.pending = append(ecs.pending, newSystemDef(system, true))
}

func (ecs *ECS) System(system interface{}) {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()
	ecs.pending = append(ecs.pending, newSystemDef(system, false))
}

func (ecs *ECS) Spawn(st interface{}) object.ID {
	os := object.MustParseStruct(st)
	id := ecs.objects.Create(os.Values())
	return id
}

func (ecs *ECS) Attach(id object.ID, st interface{}) error {
	os, err := object.ParseStruct(st)
	if err != nil {
		return err
	}
	ecs.objects.Update(id, os.Values())
	ecs.closure.Reset(id)
	return nil
}

func (ecs *ECS) Step() {
	ecs.flushPending()
	var keep []*systemEntry
	for _, entry := range ecs.systems {
		for i, de := range entry.args {
			de.cache.Do(func() {
				if de.typ.Kind() == reflect.Slice {
					keys := object.FieldTypes(de.typ.Elem())
					rows := ecs.objects.List(keys...)
					objects := make([]object.Values, 0, len(rows))
					for _, row := range rows {
						ecs.closure.Link(row.ID, &de.cache)
						objects = append(objects, row.Values)
					}
					arg := object.MustMakeSlice(de.typ, len(objects), len(objects))
					n := arg.Decode(objects...)
					de.arg = arg.Value.Slice(0, n)
				} else {
					v, ok := ecs.assets[de.typ]
					if !ok {
						panic("not ok")
					}
					de.arg = v
					ecs.closure.Link(de.typ, &de.cache)
				}
				entry.cachedArgs[i] = de.arg
			})
		}
		_ = entry.def.fn.Call(entry.cachedArgs)
		if next := entry.next(); next != nil {
			keep = append(keep, next)
		} else {
			for _, de := range entry.args {
				ecs.closure.Forget(&de.cache)
			}
		}
	}
	ecs.systems = keep
}

func (ecs *ECS) flushPending() {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()
	for _, def := range ecs.pending {
		ecs.systems = append(ecs.systems, newSystemEntry(def))
	}
	ecs.pending = nil
}
