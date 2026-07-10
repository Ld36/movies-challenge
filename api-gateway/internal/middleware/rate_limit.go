package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients = make(map[string]*clientLimiter)
	mu      sync.Mutex
)

func RateLimit(requestsPerSecond rate.Limit, burst int) gin.HandlerFunc {
	go cleanupClients()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()

		if _, exists := clients[ip]; !exists {
			clients[ip] = &clientLimiter{
				limiter: rate.NewLimiter(requestsPerSecond, burst),
			}
		}

		clients[ip].lastSeen = time.Now()

		limiter := clients[ip].limiter

		mu.Unlock()

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func cleanupClients() {
	for {
		time.Sleep(time.Minute)

		mu.Lock()

		for ip, client := range clients {
			if time.Since(client.lastSeen) > 5*time.Minute {
				delete(clients, ip)
			}
		}

		mu.Unlock()
	}
}
