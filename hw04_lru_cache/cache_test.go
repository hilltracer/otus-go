package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		// Pushing due to the size of the queue

		_ = c.Set("1", 101)
		_ = c.Set("2", 102)
		_ = c.Set("3", 103)

		_, ok := c.Get("1")
		require.True(t, ok)

		_, ok = c.Get("2")
		require.True(t, ok)

		_, ok = c.Get("3")
		require.True(t, ok)

		_ = c.Set("4", 104)

		_, ok = c.Get("1")
		require.False(t, ok)

		_, ok = c.Get("2")
		require.True(t, ok)

		_, ok = c.Get("3")
		require.True(t, ok)

		_, ok = c.Get("4")
		require.True(t, ok)

		// Clearing the cache
		c.Clear()

		_, ok = c.Get("2")
		require.False(t, ok)

		_, ok = c.Get("3")
		require.False(t, ok)

		_, ok = c.Get("4")
		require.False(t, ok)

		// Pushing due to the last access
		_ = c.Set("1", 101)
		_ = c.Set("2", 102)
		_ = c.Set("3", 103)

		_, ok = c.Get("3")
		require.True(t, ok)

		_, ok = c.Get("2")
		require.True(t, ok)

		_, ok = c.Get("1")
		require.True(t, ok)

		// "3" will be pushed out of the cache because
		// it was accessed the least recently
		_ = c.Set("4", 104)

		_, ok = c.Get("1")
		require.True(t, ok)

		_, ok = c.Get("2")
		require.True(t, ok)

		_, ok = c.Get("3")
		require.False(t, ok)

		_, ok = c.Get("4")
		require.True(t, ok)
	})

	t.Run("overwrite with same value", func(t *testing.T) {
		c := NewCache(2)

		_ = c.Set("x", 123)
		wasInCache := c.Set("x", 123)

		require.True(t, wasInCache)

		val, ok := c.Get("x")
		require.True(t, ok)
		require.Equal(t, 123, val)
	})

	t.Run("zero capacity", func(t *testing.T) {
		c := NewCache(0)

		_ = c.Set("k", 1)

		_, ok := c.Get("k")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
