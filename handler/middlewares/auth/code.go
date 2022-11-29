package auth

import (
	"fileserver/utils"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/errs"
)

const (
	CodeAuthName = "code"
)

func init() {
	Regist(CodeAuthName, func() IAuth {
		return &codeAuth{}
	})
}

type codeAuth struct {
}

func (c *codeAuth) Name() string {
	return CodeAuthName
}

func (c *codeAuth) IsMatchAuthType(ctx *gin.Context) bool {
	ak := ctx.GetHeader("x-fs-ak")
	code := ctx.GetHeader("x-fs-code")
	ts := ctx.GetHeader("x-fs-ts")
	return len(ak) != 0 && len(code) != 0 && len(ts) != 0
}

func (c *codeAuth) Auth(ctx *gin.Context, users map[string]string) (string, error) {
	ak := ctx.GetHeader("x-fs-ak")
	sk, ok := users[ak]
	if !ok {
		return "", errs.New(errs.ErrParam, "user:%s not found", ak)
	}
	ts := ctx.GetHeader("x-fs-ts")
	code := ctx.GetHeader("x-fs-code")
	if len(ts) == 0 || len(code) == 0 {
		return "", errs.New(errs.ErrParam, "invalid ts/code, ts:%s, code:%s", ts, code)
	}
	its, _ := strconv.ParseUint(ts, 10, 64)
	now := time.Now().Unix()
	if its < uint64(now) {
		return "", errs.New(errs.ErrParam, "code expire, ts:%s", ts)
	}
	realCode := utils.GetMd5([]byte(fmt.Sprintf("%s:%s:%s", ak, sk, ts)))
	if code != realCode {
		return "", errs.New(errs.ErrParam, "code not match, code carry:%s, calc:%s", code, realCode)
	}
	return ak, nil
}
