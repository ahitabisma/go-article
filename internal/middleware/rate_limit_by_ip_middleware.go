package middleware

import (
	"go-article/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimitByIP() gin.HandlerFunc {
	limiters := make(map[string]*rate.Limiter)

	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		// Jika limiter untuk IP belum ada, buat baru
		if _, exists := limiters[ip]; !exists {
			limiters[ip] = rate.NewLimiter(1, 10) // 1 request per detik dengan burst 10
		}

		if !limiters[ip].Allow() {
			res := utils.APIResponse("Too many requests. Please try again later.", http.StatusTooManyRequests, "error", nil, nil)
			ctx.AbortWithStatusJSON(429, res)
			return
		}
		ctx.Next()
	}
}
