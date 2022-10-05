package middlewares

import (
	"fileserver/handler/s3base"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/errs"
)

func S3BucketOpLimitMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		switch ctx.Request.Method {
		case http.MethodGet, http.MethodPut:
		default:
			s3base.WriteError(ctx, http.StatusBadRequest, errs.New(errs.ErrParam, "unsupport method"))
			return
		}
		path := ctx.Request.URL.Path
		path = strings.Trim(path, "/")
		parts := strings.SplitN(path, "/", 2)
		bucket := parts[0]
		s3base.SetS3Bucket(ctx, bucket)
		if len(parts) > 1 {
			obj := parts[1]
			s3base.SetS3Object(ctx, obj)
		}
	}
}
