package proxy

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*client)
)

func getClientLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	//return existing limiter
	if c, exists := clients[ip]; exists {
		c.lastSeen = time.Now()
		return c.limiter
	}
	// new limiter per ip
	c := &client{
		limiter:  rate.NewLimiter(10, 20),
		lastSeen: time.Now(),
	}
	clients[ip] = c
	return c.limiter
}
func init() {
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastSeen) > 3*time.Minute {
					log.Println("Deleting ip address: ", ip)
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
}

func getIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := getClientLimiter(getIP(r))
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
