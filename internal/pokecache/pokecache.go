package pokecache

import (
	"sync"
	"time"
)

// Cache - A structure that stores key-value pairs with timestamps
// for implementing time-based cache eviction
type Cache struct {
	cache map[string]cacheEntry // Internal map to store cached data with string keys
	mux   *sync.Mutex           // Mutex to ensure thread-safety when accessing the cache
}

// cacheEntry - Represents a single cached item with creation time and actual data
type cacheEntry struct {
	createdAt time.Time // When this entry was added to the cache (for expiration)
	val       []byte    // The actual data being cached as a byte slice
}

// NewCache - Creates and initializes a new cache with automatic cleanup
// interval: how often to check for and remove expired entries
func NewCache(interval time.Duration) Cache {
	// Initialize the cache structure with an empty map and a mutex
	c := Cache{
		cache: make(map[string]cacheEntry), // Create an empty map to store the cache entries
		mux:   &sync.Mutex{},               // Create a new mutex for thread safety
	}

	// Start the cleanup process in a separate goroutine
	// This runs in the background without blocking the main execution
	go c.reapLoop(interval)

	return c // Return the initialized cache
}

// Add - Stores a new key-value pair in the cache
// key: unique identifier for the cached data
// value: the actual data to cache as bytes
func (c *Cache) Add(key string, value []byte) {
	c.mux.Lock()         // Lock the mutex to prevent concurrent access
	defer c.mux.Unlock() // Ensure the mutex is unlocked when the function returns

	// Store the new entry with the current time
	c.cache[key] = cacheEntry{
		createdAt: time.Now().UTC(), // Use UTC time to avoid timezone issues
		val:       value,            // Store the provided bytes
	}
}

// Get - Retrieves a value from the cache by its key
// Returns the cached value and a boolean indicating if the key was found
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()         // Lock the mutex to prevent concurrent access
	defer c.mux.Unlock() // Ensure the mutex is unlocked when the function returns

	val, ok := c.cache[key] // Try to get the entry from the cache
	return val.val, ok      // Return the actual value and whether it was found
}

// reapLoop - A long-running function that periodically calls reap() to clean up expired entries
// interval: how long to wait between cleanup cycles
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval) // Create a ticker that "ticks" at the specified interval
	for range ticker.C {               // Continuously wait for the next tick
		c.reap(time.Now().UTC(), interval) // Call reap with the current time when a tick occurs
	}
}

// reap - Removes cache entries that are older than the specified duration
// now: the current time (used to determine entry age)
// last: duration threshold - entries older than this will be removed
func (c *Cache) reap(now time.Time, last time.Duration) {
	c.mux.Lock()         // Lock the mutex to prevent concurrent access
	defer c.mux.Unlock() // Ensure the mutex is unlocked when the function returns

	// Iterate through all entries in the cache
	for k, v := range c.cache {
		// Check if the entry's creation time is before (now - interval)
		// This means the entry is older than the specified duration
		if v.createdAt.Before(now.Add(-last)) {
			delete(c.cache, k) // Remove the expired entry from the cache
		}
	}
}
