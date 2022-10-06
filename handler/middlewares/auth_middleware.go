package middlewares

import (
	"encoding/base64"
	"fileserver/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

type IAuth interface {
	Name() string
	Auth(ctx *gin.Context, users map[string]string) (string, bool, error)
}

var authList = []IAuth{
	&codeAuth{},
	&basicAuth{},
}

type codeAuth struct {
}

func (c *codeAuth) Name() string {
	return "code_auth"
}

func (c *codeAuth) Auth(ctx *gin.Context, users map[string]string) (string, bool, error) {
	ak := ctx.GetHeader("x-fs-ak")
	sk, ok := users[ak]
	if !ok {
		return "", false, nil
	}
	ts := ctx.GetHeader("x-fs-ts")
	code := ctx.GetHeader("x-fs-code")
	if len(ts) == 0 || len(code) == 0 {
		return "", false, nil
	}
	its, _ := strconv.ParseUint(ts, 10, 64)
	now := time.Now().Unix()
	if its < uint64(now) {
		return "", false, errs.New(errs.ErrParam, "code expire, ts:%s", ts)
	}
	realCode := utils.GetMd5([]byte(fmt.Sprintf("%s:%s:%s", ak, sk, ts)))
	if code == realCode {
		return ak, true, nil
	}
	return "", false, nil
}

type basicAuth struct {
}

func (c *basicAuth) Name() string {
	return "basic_auth"
}

func (b *basicAuth) Auth(ctx *gin.Context, users map[string]string) (string, bool, error) {
	auth := ctx.GetHeader("Authorization")
	if len(auth) == 0 {
		return "", false, nil
	}
	authData := strings.SplitN(auth, " ", 2)
	if len(authData) != 2 {
		return "", false, errs.New(errs.ErrParam, "invalid auth data:%s", auth)
	}
	if authData[0] != "Basic" {
		return "", false, errs.New(errs.ErrParam, "invalid auth prefix, data:%s", auth)
	}
	bdata, err := base64.StdEncoding.DecodeString(authData[1])
	if err != nil {
		return "", false, errs.Wrap(errs.ErrParam, "decode auth data fail", err)
	}
	data := string(bdata)
	userdata := strings.SplitN(data, ":", 2)
	if len(userdata) != 2 {
		return "", false, errs.New(errs.ErrParam, "invalid user pwd data:%s", data)
	}
	sk, ok := users[userdata[0]]
	if !ok || sk != userdata[1] {
		return "", false, nil
	}
	return userdata[0], true, nil
}

// func S3V4Auth(ctx *gin.Context, u, p string) (bool, error) {
// 	if !s3base.IsRequestSignatureV4(ctx.Request) {
// 		return false, nil
// 	}
// 	parsed, _, err := s3base.ParseV4Signature(ctx.Request)
// 	if err != nil {
// 		return false, errs.Wrap(errs.ErrParam, "parse v4 signature fail", err)
// 	}
// 	if u != parsed.AKey {
// 		return false, nil
// 	}
// 	pass, err := s3base.S3AuthV4(ctx.Request, u, p, parsed)
// 	if err != nil {
// 		return false, err
// 	}
// 	if !pass {
// 		return false, nil
// 	}
// 	return true, nil
// }

func CommonAuth(users map[string]string) gin.HandlerFunc {
	return CommonAuthMiddleware(users, authList...)
}

func CommonAuthMiddleware(users map[string]string, ats ...IAuth) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := logutil.GetLogger(ctx).With(zap.String("method", ctx.Request.Method), zap.String("path", ctx.Request.URL.Path), zap.String("ip", ctx.ClientIP()))
		for _, fn := range ats {
			ak, ok, err := fn.Auth(ctx, users)
			if err != nil {
				logger.With(zap.String("auth", fn.Name()), zap.Error(err)).Error("auth error")
				ctx.AbortWithError(http.StatusUnauthorized, errs.Wrap(errs.ErrUnknown, "internal services error", err))
				return
			}
			if ok {
				logger.With(zap.String("auth", fn.Name()), zap.String("ak", ak)).Debug("user auth succ")
				return
			}
		}
		logger.Error("need auth")
		ctx.AbortWithError(http.StatusUnauthorized, errs.New(errs.ErrParam, "need auth"))
	}
}
