package internal

import (
	"context"
	internalCacheVendor "github.com/Code-Hex/go-generics-cache"
	"time"
)

var internalCache *internalCacheVendor.Cache[string, string]

func InitCache() {
	internalCache = internalCacheVendor.NewContext[string, string](context.Background())
}

// cache decorates expensive data fetch call with a caching layer.
func Cache(key string, ttl time.Duration, expensiveCall func() (string, error)) (string, error) {
	var (
		val string
		err error
		ok  bool
	)

	val, ok = internalCache.Get(key)
	if ok {
		return val, nil // fast track: return from the cache
	}

	val, err = expensiveCall() // slow track: jump into the expensive call

	if err == nil {
		internalCache.Set(key, val, internalCacheVendor.WithExpiration(ttl))
	}

	return val, err
}
