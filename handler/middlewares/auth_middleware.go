package middlewares

import (
	"encoding/base64"
	"fileserver/handler/s3base"
	"fileserver/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

type AuthFunc func(ctx *gin.Context, u, p string) (bool, error)

func CodeAuth(ctx *gin.Context, u, p string) (bool, error) {
	ak := ctx.GetHeader("x-fs-ak")
	ts := ctx.GetHeader("x-fs-ts")
	code := ctx.GetHeader("x-fs-code")
	if len(ts) == 0 || len(code) == 0 || ak != u {
		return false, nil
	}
	its, _ := strconv.ParseUint(ts, 10, 64)
	now := time.Now().Unix()
	if its+60 < uint64(now) {
		return false, errs.New(errs.ErrParam, "code expire, ts:%s", ts)
	}
	realCode := utils.GetMd5([]byte(fmt.Sprintf("%s:%s:%s", u, p, ts)))
	if code == realCode {
		return true, nil
	}
	return false, nil
}

func BasicAuth(ctx *gin.Context, u, p string) (bool, error) {
	auth := ctx.GetHeader("Authorization")
	if len(auth) == 0 {
		return false, nil
	}
	base := u + ":" + p
	get := "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
	if auth == get {
		return true, nil
	}
	return false, nil
}

func S3V4Auth(ctx *gin.Context, u, p string) (bool, error) {
	if !s3base.IsRequestSignatureV4(ctx.Request) {
		return false, nil
	}
	parsed, _, err := s3base.ParseV4Signature(ctx.Request)
	if err != nil {
		return false, errs.Wrap(errs.ErrParam, "parse v4 signature fail", err)
	}
	if u != parsed.AKey {
		return false, nil
	}
	pass, err := s3base.S3AuthV4(ctx.Request, u, p, parsed)
	if err != nil {
		return false, err
	}
	if !pass {
		return false, nil
	}
	return true, nil
}

func CommonAuth(users map[string]string) gin.HandlerFunc {
	fns := []AuthFunc{
		CodeAuth,
		BasicAuth,
		S3V4Auth,
	}
	return CommonAuthMiddleware(users, fns...)
}

func CommonAuthMiddleware(users map[string]string, fns ...AuthFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for u, p := range users {
			for _, fn := range fns {
				ok, err := fn(ctx, u, p)
				if err != nil {
					logutil.GetLogger(ctx).With(zap.Error(err)).Error("auth error")
					ctx.AbortWithError(http.StatusUnauthorized, errs.Wrap(errs.ErrUnknown, "internal services error", err))
					return
				}
				if ok {
					return
				}
			}
		}
		logutil.GetLogger(ctx).Error("need auth")
		ctx.AbortWithError(http.StatusUnauthorized, errs.New(errs.ErrParam, "need auth"))
	}
}
