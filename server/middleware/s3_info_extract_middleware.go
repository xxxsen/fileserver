package middleware

import (
	"fileserver/proxyutil"
	"fileserver/server/handler/s3/s3base"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ExtractS3InfoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet, http.MethodPut:
		default:
			s3base.WriteError(c, http.StatusBadRequest, fmt.Errorf("unsupport method"))
			return
		}
		path := c.Request.URL.Path
		path = strings.Trim(path, "/")
		parts := strings.SplitN(path, "/", 2)
		bucket := parts[0]
		s3info := &proxyutil.S3Info{
			Bucket: bucket,
		}
		if len(parts) > 1 {
			obj := parts[1]
			s3info.FileID = obj
		}
		c.Request = c.Request.WithContext(proxyutil.SetS3Info(c.Request.Context(), s3info))
	}
}
