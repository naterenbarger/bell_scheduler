package ratelimiter

import (
    "sync"
    "time"
)

type RateLimiter struct {
    mu       sync.RWMutex
    attempts map[string][]time.Time
    window   time.Duration
    max      int
}

func NewRateLimiter(window time.Duration, max int) *RateLimiter {
    return &RateLimiter{
        attempts: make(map[string][]time.Time),
        window:   window,
        max:      max,
    }
}

func (rl *RateLimiter) IsAllowed(key string) (bool, time.Duration) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    now := time.Now()
    windowStart := now.Add(-rl.window)

    // Clean up old attempts
    attempts := rl.attempts[key]
    valid := attempts[:0]
    for _, t := range attempts {
        if t.After(windowStart) {
            valid = append(valid, t)
        }
    }
    rl.attempts[key] = valid

    // Check if rate limit is exceeded
    if len(valid) >= rl.max {
        // Calculate time until next allowed attempt
        oldestValid := valid[0]
        waitTime := rl.window - now.Sub(oldestValid)
        return false, waitTime
    }

    // Add new attempt
    rl.attempts[key] = append(rl.attempts[key], now)
    return true, 0
}

func (rl *RateLimiter) Clear(key string) {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    delete(rl.attempts, key)
} 