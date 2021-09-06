package object

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeSlice(t *testing.T) {
	type A struct {
		X int
		Y *float64
	}
	t.Run("ok", func(t *testing.T) {
		st := reflect.TypeOf([]A(nil))
		_, err := MakeSlice(st, 2, 2)
		require.NoError(t, err)
	})
	t.Run("err1", func(t *testing.T) {
		st := reflect.TypeOf([]int(nil))
		_, err := MakeSlice(st, 2, 2)
		require.Error(t, err)
	})
	t.Run("err2", func(t *testing.T) {
		st := reflect.TypeOf(int(0))
		_, err := MakeSlice(st, 2, 2)
		require.Error(t, err)
	})
}

func TestSlice_Decode(t *testing.T) {
	type A struct {
		X int
		Y *float64
	}
	a0 := MustParseStruct(A{X: 13})
	a1 := MustParseStruct(A{X: 42})
	a2 := MustParseStruct(A{X: 18})
	objects := []Values{a0.Values(), a1.Values(), a2.Values()}
	sl := MustMakeSlice(reflect.TypeOf([]A(nil)), 2, 2)
	sl.Decode(objects...)
	assert.Equal(t, []A{{X: 13}, {X: 42}}, sl.Value.Interface())
}
