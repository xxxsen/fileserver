package proxyutil

import (
	"context"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

type HandleFunc func(c *gin.Context, ctx context.Context, request interface{})

func innerHandle(fn HandleFunc, ptr interface{}, onDecode func(c *gin.Context, input interface{}) error) gin.HandlerFunc {
	typ := reflect.TypeOf(ptr)
	return func(c *gin.Context) {
		val := reflect.New(typ.Elem())
		iface := val.Interface()
		ctx := c.Request.Context()
		if err := onDecode(c, iface); err != nil {
			logutil.GetLogger(ctx).Error("bind object failed", zap.Error(err), zap.String("uri", c.Request.RequestURI))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		fn(c, ctx, iface)
	}

}

func WrapBizFunc(fn HandleFunc, ptr interface{}) gin.HandlerFunc {
	return innerHandle(fn, ptr, func(c *gin.Context, input interface{}) error {
		return c.ShouldBind(input)
	})
}
