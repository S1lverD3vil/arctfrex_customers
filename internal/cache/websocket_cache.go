package cache

import (
	"github.com/patrickmn/go-cache"
)

// Global cache instance accessible throughout the package
var ClientCache *cache.Cache

func init() {
	// Initialize the cache without expiration and no cleanup
	ClientCache = cache.New(cache.NoExpiration, cache.NoExpiration)
}
