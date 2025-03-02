package middlewares

import (
	"fileserver/handler/middlewares/auth"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

func CreateFullAuthMethods() []auth.IAuth {
	authList := []auth.IAuth{}
	lst := auth.AuthList()
	for _, name := range lst {
		ath := auth.MustCreateByName(name)
		authList = append(authList, ath)
	}
	return authList
}

func CommonAuth(users map[string]string) gin.HandlerFunc {
	return CommonAuthMiddleware(users, CreateFullAuthMethods()...)
}

func CommonAuthMiddleware(users map[string]string, ats ...auth.IAuth) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := logutil.GetLogger(ctx).With(zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.URL.Path), zap.String("ip", ctx.ClientIP()))

		for _, fn := range ats {
			if !fn.IsMatchAuthType(ctx) {
				continue
			}
			ak, err := fn.Auth(ctx, users)
			if err != nil {
				logger.Error("auth error", zap.String("auth", fn.Name()), zap.Error(err))
				ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("internal services error, err:%w", err))
				return
			}
			logger.Debug("user auth succ", zap.String("auth", fn.Name()), zap.String("ak", ak))
			return
		}
		logger.Error("need auth")
		ctx.AbortWithError(http.StatusUnauthorized, errs.New(errs.ErrParam, "need auth"))
	}
}
