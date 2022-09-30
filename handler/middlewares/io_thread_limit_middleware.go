package middlewares

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/semaphore"
)

type IOTheadContext struct {
	cnt int64
	wg  *semaphore.Weighted
}

func NewIOThreadContext(mt int64) *IOTheadContext {
	return &IOTheadContext{wg: semaphore.NewWeighted(mt), cnt: mt}
}

func IOThreadLimitMiddleware(ctx *IOTheadContext) gin.HandlerFunc {
	return func(g *gin.Context) {
		if ctx.cnt == 0 {
			return
		}
		ctx.wg.Acquire(g, 1)
		defer ctx.wg.Release(1)
		g.Next()
	}
}
