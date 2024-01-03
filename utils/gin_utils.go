package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/cgi"
	"github.com/xxxsen/common/cgi/codec"
)

func GinToProcessFunc(fn gin.HandlerFunc) cgi.ProcessFunc {
	return func(c *gin.Context, req interface{}) (int, interface{}, error) {
		fn(c)
		return http.StatusOK, nil, nil
	}
}

func WrapGinFunc(fn gin.HandlerFunc) gin.HandlerFunc {
	return cgi.WrapHandler(nil, codec.NopCodec, GinToProcessFunc(fn))
}
