package tokenbucket

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

func NewLimiter(r rate.Limit, b int, t time.Duration) gin.HandlerFunc {
	limiters := &sync.Map{}

	return func(c *gin.Context) {
		key := c.FullPath()
		l, _ := limiters.LoadOrStore(key, rate.NewLimiter(r, b))
		ctx, cancel := context.WithTimeout(c, t)
		defer cancel()

		if err := l.(*rate.Limiter).Wait(ctx); err != nil {
			fmt.Println("限流...")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": -1})
		}
		c.Next()
	}
}
