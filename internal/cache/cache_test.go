package cache

import (
	"testing"
	"time"
)

func TestMemoryCache_SetAndGet(t *testing.T) {
	cache := NewMemoryCache(time.Hour)
	
	cache.Set("test-key", "test-value", time.Minute)
	
	value, exists := cache.Get("test-key")
	if !exists {
		t.Fatal("Expected key to exist")
	}
	
	if value != "test-value" {
		t.Fatalf("Expected 'test-value', got %v", value)
	}
}

func TestMemoryCache_Expiration(t *testing.T) {
	cache := NewMemoryCache(time.Hour)
	
	cache.Set("test-key", "test-value", time.Millisecond*10)
	
	time.Sleep(time.Millisecond * 20)
	
	_, exists := cache.Get("test-key")
	if exists {
		t.Fatal("Expected key to be expired")
	}
}

func TestMemoryCache_Delete(t *testing.T) {
	cache := NewMemoryCache(time.Hour)
	
	cache.Set("test-key", "test-value", time.Minute)
	cache.Delete("test-key")
	
	_, exists := cache.Get("test-key")
	if exists {
		t.Fatal("Expected key to be deleted")
	}
}

func TestMemoryCache_Clear(t *testing.T) {
	cache := NewMemoryCache(time.Hour)
	
	cache.Set("key1", "value1", time.Minute)
	cache.Set("key2", "value2", time.Minute)
	
	if cache.Size() != 2 {
		t.Fatalf("Expected size 2, got %d", cache.Size())
	}
	
	cache.Clear()
	
	if cache.Size() != 0 {
		t.Fatalf("Expected size 0 after clear, got %d", cache.Size())
	}
}