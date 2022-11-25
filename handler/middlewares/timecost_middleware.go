package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

func TimeCostMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		logger := logutil.GetLogger(ctx).With(zap.Time("request_start", start), zap.String("path", ctx.Request.URL.Path))
		logger.Debug("recv request")
		defer func() {
			cost := time.Since(start)
			logger.Debug("handle request finish", zap.Int64("cost(ms)", int64(cost/time.Millisecond)))
		}()
		ctx.Next()
	}
}
