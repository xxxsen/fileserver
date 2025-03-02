package proxyutil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	defaultBizErrCode uint32 = 100000
)

type CommonResponse struct {
	Code    uint32
	Message string
	Data    interface{}
}

type iCodeErr interface {
	Code() uint32
}

func makePacket(code uint32, msg string, obj interface{}) *CommonResponse {
	return &CommonResponse{
		Code:    code,
		Message: msg,
		Data:    obj,
	}
}

func Success(c *gin.Context, obj interface{}) {
	c.JSON(http.StatusOK, makePacket(0, "", obj))
}

func Fail(c *gin.Context, code int, err error) {
	bizCode := defaultBizErrCode
	errmsg := err.Error()
	if ie, ok := err.(iCodeErr); ok {
		bizCode = ie.Code()
	}
	c.AbortWithStatusJSON(code, makePacket(bizCode, errmsg, nil))
}
