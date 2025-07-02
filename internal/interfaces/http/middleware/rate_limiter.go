package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	redisClient *redis.Client
	limit       int
	period      int // in seconds
}

func NewRateLimiter(redisClient *redis.Client, limit, period int) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		limit:       limit,
		period:      period,
	}
}

func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("rate_limit:%s", c.ClientIP())
		ctx := context.Background()

		// Get the current count
		count, err := rl.redisClient.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rate_limit_error", "message": "Error checking rate limit"})
			c.Abort()
			return
		}

		if count >= rl.limit {
			c.Header("X-RateLimit-Limit", strconv.Itoa(rl.limit))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Unix()+int64(rl.period), 10))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate_limit_exceeded",
				"message": fmt.Sprintf("Rate limit of %d requests per %d seconds exceeded", rl.limit, rl.period),
			})
			c.Abort()
			return
		}

		// Increment the counter
		if count == 0 {
			err = rl.redisClient.SetNX(ctx, key, 1, time.Duration(rl.period)*time.Second).Err()
		} else {
			err = rl.redisClient.Incr(ctx, key).Err()
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rate_limit_error", "message": "Error updating rate limit"})
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(rl.limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(rl.limit-count-1))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Unix()+int64(rl.period), 10))

		c.Next()
	}
}

func (rl *RateLimiter) LoginLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginReq struct {
			Email string `json:"email"`
		}
		c.ShouldBindJSON(&loginReq)
		key := fmt.Sprintf("rate_limit:login:%s", loginReq.Email)
		ctx := context.Background()

		count, err := rl.redisClient.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rate_limit_error", "message": "Error checking rate limit"})
			c.Abort()
			return
		}

		if count >= rl.limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate_limit_exceeded",
				"message": "Too many login attempts. Please try again later.",
			})
			c.Abort()
			return
		}

		if count == 0 {
			err = rl.redisClient.SetNX(ctx, key, 1, time.Duration(rl.period)*time.Second).Err()
		} else {
			err = rl.redisClient.Incr(ctx, key).Err()
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rate_limit_error", "message": "Error updating rate limit"})
			c.Abort()
			return
		}

		c.Next()
	}
}
