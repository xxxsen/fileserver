package middlewares

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

var NullReferer = "null"
var errInvalidReferer = fmt.Errorf("invalid referer")

func RefererMiddleware(enable bool, valids []string) gin.HandlerFunc {
	m := map[string]bool{}
	supportNil := false
	for _, item := range valids {
		if strings.EqualFold(item, NullReferer) {
			supportNil = true
			continue
		}
		m[item] = true
	}
	return func(ctx *gin.Context) {
		if !enable {
			return
		}
		logger := logutil.GetLogger(ctx).With(zap.String("method", ctx.Request.Method), zap.String("path", ctx.Request.URL.Path), zap.String("ip", ctx.ClientIP()))
		if len(ctx.Request.Referer()) == 0 {
			if !supportNil {
				ctx.AbortWithError(http.StatusForbidden, errInvalidReferer)
				logger.Error("nil referer")
				return
			}
			return
		}
		uri, err := url.Parse(ctx.Request.Referer())
		if err != nil {
			logger.Error("decode referer fail", zap.Error(err), zap.String("referer", ctx.Request.Referer()))
			ctx.AbortWithError(http.StatusForbidden, errInvalidReferer)
			return
		}
		if _, ok := m[uri.Host]; !ok {
			logger.Error("referer not in white list", zap.String("referer_host", uri.Host))
			ctx.AbortWithError(http.StatusForbidden, errInvalidReferer)
			return
		}
	}
}
