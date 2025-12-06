package rate_limiter

import (
	"time"
)

type RateLimiter struct {
	tokens chan struct{}
}

func NewRateLimiter(burstSize int, refillInterval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, burstSize),
	}

	// fill initial burst
	for i := 0; i < burstSize; i++ {
		rl.tokens <- struct{}{}
	}

	// refill tokens without blocking
	go func() {
		ticker := time.NewTicker(refillInterval)
		for range ticker.C {
			select {
			case rl.tokens <- struct{}{}:
			default:
				// channel full → don't block
			}
		}
	}()

	return rl
}

func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}
