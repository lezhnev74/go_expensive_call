package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
	"golang.org/x/time/rate"
	"lezhnev74/go_expensive_call/internal"
	"log"
	"net/http"
	"time"
)

func ExpensiveCall() (string, error) {
	time.Sleep(time.Second)
	log.Printf("Expensive call invoked.")
	result := fmt.Sprintf("Done at %s", time.Now())
	return result, nil
}

func main() {
	r := gin.Default()
	internal.InitCache()
	g := new(singleflight.Group)
	l := rate.NewLimiter(rate.Every(time.Second), 1)

	r.GET("/calculate", func(c *gin.Context) {
		// Level-3: Throttling
		allowed := l.Allow()
		if !allowed {
			c.JSON(http.StatusTooManyRequests, nil) // <-- backoff immediately
			return
		}

		cacheKey := "expensive"
		ttl := 10 * time.Second
		decorated := func() (any, error) {
			// Level-1: Caching
			return internal.Cache(cacheKey, ttl, ExpensiveCall)
		}

		// Level-2: Queueing
		result, err, _ := g.Do(cacheKey, decorated) // <-- Here it queues all the callers
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, result)
	})

	err := r.Run("0.0.0.0:8088")
	if err != nil {
		panic(err)
	}
}
