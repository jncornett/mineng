package object

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStruct(t *testing.T) {
	type A struct{}
	t.Run("ok", func(t *testing.T) {
		_, err := NewStruct(reflect.ValueOf(A{}))
		require.NoError(t, err)
	})
	t.Run("ok pointer", func(t *testing.T) {
		_, err := NewStruct(reflect.ValueOf(&A{}))
		require.NoError(t, err)
	})
	t.Run("nil pointer", func(t *testing.T) {
		_, err := NewStruct(reflect.ValueOf((*A)(nil)))
		require.Error(t, err)
	})
	t.Run("unstruct", func(t *testing.T) {
		_, err := NewStruct(reflect.ValueOf(0))
		require.Error(t, err)
	})
}

func TestStruct_Decode(t *testing.T) {
	type A struct {
		X int
		Y *float64
	}
	var (
		x = 42
		y = new(float64)
	)
	*y = 3.14
	var (
		xv = reflect.ValueOf(x)
		yv = reflect.ValueOf(y)
	)
	t.Run("complete", func(t *testing.T) {
		vs := Values{xv.Type(): xv, yv.Type(): yv}
		var a A
		st := MustParseStruct(&a)
		st.Decode(vs)
		assert.Equal(t, A{X: x, Y: y}, a)
	})
	t.Run("subset", func(t *testing.T) {
		vs := Values{xv.Type(): xv}
		var a A
		st := MustParseStruct(&a)
		st.Decode(vs)
		assert.Equal(t, A{X: x}, a)
	})
}

func TestStruct_StrictDecode(t *testing.T) {
	type A struct {
		X int
		Y *float64
	}
	var (
		x = 42
		y = new(float64)
	)
	*y = 3.14
	var (
		xv = reflect.ValueOf(x)
		yv = reflect.ValueOf(y)
	)
	t.Run("complete", func(t *testing.T) {
		vs := Values{xv.Type(): xv, yv.Type(): yv}
		var a A
		st := MustParseStruct(&a)
		err := st.StrictDecode(vs)
		require.NoError(t, err)
		assert.Equal(t, A{X: x, Y: y}, a)
	})
	t.Run("subset", func(t *testing.T) {
		vs := Values{xv.Type(): xv}
		var a A
		st := MustParseStruct(&a)
		err := st.StrictDecode(vs)
		require.Error(t, err)
	})
}
func BenchmarkStruct_Decode(b *testing.B) {
	type A struct {
		X int
		Y *float64
	}
	var (
		x = 42
		y = new(float64)
	)
	*y = 3.14
	var (
		xv = reflect.ValueOf(x)
		yv = reflect.ValueOf(y)
	)
	var a A
	vs := Values{xv.Type(): xv, yv.Type(): yv}
	st := MustParseStruct(&a)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st.Decode(vs)
	}
}
