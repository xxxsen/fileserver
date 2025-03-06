package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xxxsen/common/trace"
)

const (
	defaultRequestIdKey = "x-request-id"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqid := c.Request.Header.Get(defaultRequestIdKey)
		if len(reqid) == 0 {
			reqid = uuid.NewString()
		}
		ctx := c.Request.Context()
		ctx = trace.WithTraceId(ctx, reqid)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
