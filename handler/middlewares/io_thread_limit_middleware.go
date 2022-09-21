package middlewares

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/semaphore"
)

type IOTheadContext struct {
	wg *semaphore.Weighted
}

func NewIOThreadContext(mt int64) *IOTheadContext {
	return &IOTheadContext{semaphore.NewWeighted(mt)}
}

func IOThreadLimitMiddleware(ctx *IOTheadContext) gin.HandlerFunc {
	return func(g *gin.Context) {
		ctx.wg.Acquire(g, 1)
		defer ctx.wg.Release(1)
		g.Next()
	}
}
