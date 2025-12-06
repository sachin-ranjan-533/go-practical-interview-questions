package main

import (
	"fmt"
	"rate-limiter/rate_limiter"
	"time"
)

func main() {
	rateLimiter := rate_limiter.NewRateLimiter(5, 3*time.Second)

	for i := 0; i < 100; i++ {
		if rateLimiter.Allow() {
			fmt.Println("Request", i, "allowed at", time.Now().Format("15:04:05"))
		} else {
			fmt.Println("Request", i, "DENIED at", time.Now().Format("15:04:05"))
		}

		time.Sleep(500 * time.Millisecond)
	}
}
