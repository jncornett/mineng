package dep

// Dep represents a data dependency that can be invalidated by a call to Reset().
type Dep interface {
	// Reset invalidates a data dependency.
	// Calls to Reset are assumed to be idempotent.
	Reset()
}

// Graph represents a set of data dependencies.
type Graph map[Dep]struct{}

// Reset resets all the data dependencies in g.
func (g Graph) Reset() {
	for dep := range g {
		dep.Reset()
	}
}

// Add associates a data dependency to g.
// Add is idempotent.
func (g Graph) Add(d Dep) {
	g[d] = struct{}{}
}

// Remove removes an association of a data dependency to g.
// Remove is idempotent.
func (g Graph) Remove(d Dep) {
	delete(g, d)
}

// Closure represents a set data dependency graphs.
type Closure struct {
	forest map[interface{}]Graph
	index  map[Dep]map[interface{}]struct{}
}

func NewClosure() Closure {
	return Closure{
		forest: map[interface{}]Graph{},
		index:  map[Dep]map[interface{}]struct{}{},
	}
}

func (cl Closure) Link(key interface{}, d Dep) {
	g, ok := cl.forest[key]
	if !ok {
		g = Graph{}
		cl.forest[key] = g
	}
	g.Add(d)
	keys, ok := cl.index[d]
	if !ok {
		keys = map[interface{}]struct{}{}
		cl.index[d] = keys
	}
	keys[key] = struct{}{}
}

func (cl Closure) Forget(d Dep) {
	for key := range cl.index[d] {
		cl.Unlink(key, d)
	}
	delete(cl.index, d)
}

func (cl Closure) Unlink(key interface{}, d Dep) {
	g, ok := cl.forest[key]
	if !ok {
		return
	}
	g.Remove(d)
}

func (cl Closure) Flush(key interface{}) {
	delete(cl.forest, key)
}

func (cl Closure) Reset(key interface{}) {
	cl.forest[key].Reset()
}
