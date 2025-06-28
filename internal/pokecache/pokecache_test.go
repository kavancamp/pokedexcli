package pokecache

import (
	"testing"
	"time"
)

func TestCacheAddAndGet(t *testing.T) {
	cache := NewCache(10 * time.Second)

	key := "https://pokeapi.co/test"
	val := []byte("pikachu!")

	cache.Add(key, val)
	got, ok := cache.Get(key)

	if !ok {
		t.Errorf("Expected key %q to be in cache", key)
	}
	if string(got) != string(val) {
		t.Errorf("Expected value %q, got %q", val, got)
	}
}
func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}