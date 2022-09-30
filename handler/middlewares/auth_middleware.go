package middlewares

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/errs"
)

type AuthFunc func(ctx *gin.Context, u, p string) (bool, error)

func CodeAuth(ctx *gin.Context, u, p string) (bool, error) {
	code := ctx.GetHeader("x-fs-code")
	if code == p {
		return true, nil
	}
	return false, nil
}

func BasicAuth(ctx *gin.Context, u, p string) (bool, error) {
	auth := ctx.GetHeader("Authorization")
	base := u + ":" + p
	get := "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
	if auth == get {
		return true, nil
	}
	return false, nil
}

func S3V2Auth(ctx *gin.Context, u, p string) (bool, error) {
	//TODO: finish it
	return false, nil
}

func CommonAuth(users map[string]string) gin.HandlerFunc {
	fns := []AuthFunc{
		CodeAuth,
		BasicAuth,
		S3V2Auth,
	}
	return CommonAuthMiddleware(users, fns...)
}

func CommonAuthMiddleware(users map[string]string, fns ...AuthFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for u, p := range users {
			for _, fn := range fns {
				ok, err := fn(ctx, u, p)
				if err != nil {
					ctx.AbortWithError(http.StatusUnauthorized, errs.Wrap(errs.ErrUnknown, "internal services error", err))
					return
				}
				if ok {
					return
				}
			}
		}
		ctx.AbortWithError(http.StatusUnauthorized, errs.New(errs.ErrParam, "need auth"))
	}
}
