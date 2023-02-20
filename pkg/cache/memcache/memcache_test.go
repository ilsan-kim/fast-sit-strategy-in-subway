package memcache

import (
	"testing"
	"time"
)

func TestMemoryCache_Get(t *testing.T) {
	cache := NewMemCache()

	if value, exists := cache.Get("non-existent-key"); value != nil || exists != false {
		t.Errorf("Expected Get to return (nil, false), but got (%v, %v)", value, exists)
	}

	cache.Set("test-key", "test-value", time.Minute)
	if value, exists := cache.Get("test-key"); value != "test-value" || exists != true {
		t.Errorf("Expected Get to return (\"test-value\", true), but got (%v, %v)", value, exists)
	}

	cache.Set("expired-key", "expired-value", time.Nanosecond)
	time.Sleep(time.Nanosecond) // wait for value to expire
	if value, exists := cache.Get("expired-key"); value != nil || exists != false {
		t.Errorf("Expected Get to return (nil, false) for expired key, but got (%v, %v)", value, exists)
	}
}

func TestMemoryCache_Set(t *testing.T) {
	cache := NewMemCache()

	cache.Set("test-key", "test-value", time.Minute)
	if value, exists := cache.Get("test-key"); value != "test-value" || exists != true {
		t.Errorf("Expected Set to set the value and Get to return it, but got (%v, %v)", value, exists)
	}

	cache.Set("test-key", "new-value", time.Minute)
	if value, exists := cache.Get("test-key"); value != "new-value" || exists != true {
		t.Errorf("Expected Set to overwrite the value and Get to return the new value, but got (%v, %v)", value, exists)
	}
}
