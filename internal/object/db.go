package object

import (
	"reflect"
	"sync/atomic"
)

type ID uint64

type Row struct {
	ID     ID
	Values Values
}

type Values map[reflect.Type]reflect.Value

func (vs Values) Encode(v interface{}) reflect.Type {
	rv := reflect.ValueOf(v)
	t := rv.Type()
	vs[t] = rv
	return t
}

func (vs Values) Clone() Values {
	out := make(Values, len(vs))
	for t, v := range vs {
		out[t] = v
	}
	return out
}

type DB struct {
	idState uint64
	objects map[ID]Values
	index   map[reflect.Type]map[ID]struct{}
}

func NewDB() *DB {
	return &DB{
		objects: map[ID]Values{},
		index:   map[reflect.Type]map[ID]struct{}{},
	}
}

func (db *DB) Get(id ID) (vs Values, ok bool) {
	vs, ok = db.objects[id]
	return vs, ok
}

func (db *DB) Create(vs Values) ID {
	id := db.nextID()
	idv := reflect.ValueOf(id)
	vs = vs.Clone()
	vs[idv.Type()] = idv
	db.objects[id] = vs
	for t := range vs {
		db.addIndex(t, id)
	}
	return id
}

func (db *DB) Update(id ID, add Values, remove ...reflect.Type) {
	vs, ok := db.objects[id]
	if !ok {
		return
	}
	for t, v := range add {
		vs[t] = v
		db.addIndex(t, id)
	}
	for _, t := range remove {
		delete(vs, t)
		db.removeIndex(t, id)
	}
}

func (db *DB) Delete(id ID) {
	vs := db.objects[id]
	delete(db.objects, id)
	for t := range vs {
		db.removeIndex(t, id)
	}
}

func (db *DB) List(types ...reflect.Type) []Row {
	set := map[ID]struct{}{}
	for i, t := range types {
		if i == 0 {
			for id := range db.index[t] {
				set[id] = struct{}{}
			}
			continue
		}
		oset := db.index[t]
		// intersection
		for id := range set {
			if _, ok := oset[id]; !ok {
				delete(set, id)
			}
		}
		if len(set) == 0 {
			break
		}
	}
	var out []Row
	for id := range set {
		out = append(out, Row{ID: id, Values: db.objects[id].Clone()})
	}
	return out
}

func (db *DB) nextID() ID {
	return ID(atomic.AddUint64(&db.idState, 1))
}

func (db *DB) addIndex(t reflect.Type, id ID) {
	set, ok := db.index[t]
	if !ok {
		set = map[ID]struct{}{}
		db.index[t] = set
	}
	set[id] = struct{}{}
}

func (db *DB) removeIndex(t reflect.Type, id ID) {
	set, ok := db.index[t]
	if !ok {
		return
	}
	delete(set, id)
}
