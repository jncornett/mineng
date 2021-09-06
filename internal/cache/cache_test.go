package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	t.Run("multiple Do calls only fires once before a call to Reset", func(t *testing.T) {
		var c Cache
		var n int
		c.Do(func() { n++ })
		c.Do(func() { n++ })
		c.Do(func() { n++ })
		assert.Equal(t, 1, n)
	})
	t.Run("calling Reset allows Do to be called again", func(t *testing.T) {
		var c Cache
		var n int
		c.Do(func() { n++ })
		assert.Equal(t, 1, n)
		c.Reset()
		c.Do(func() { n++ })
		assert.Equal(t, 2, n)
	})
	t.Run("calling Reset before calling Do for the first time has no effect", func(t *testing.T) {
		var c Cache
		var n int
		c.Reset()
		c.Do(func() { n++ })
		assert.Equal(t, 1, n)
		c.Reset()
		c.Do(func() { n++ })
		assert.Equal(t, 2, n)
	})
}
