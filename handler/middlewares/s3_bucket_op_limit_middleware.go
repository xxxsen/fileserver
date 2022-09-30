package middlewares

import (
	"fileserver/handler/s3base"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/errs"
)

func S3BucketOpLimitMiddleware(prefix string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		switch ctx.Request.Method {
		case http.MethodGet, http.MethodPut:
		default:
			s3base.WriteError(ctx, http.StatusBadRequest, errs.New(errs.ErrParam, "unsupport method"))
			return
		}
		path := ctx.Request.URL.Path
		if len(prefix) > 0 {
			if !strings.HasPrefix(path, prefix) {
				s3base.WriteError(ctx, http.StatusBadRequest, errs.New(errs.ErrParam, "not contains prefix:%s", prefix))
				return
			}
			path = strings.TrimLeft(path, prefix)
		}
		path = strings.Trim(path, "/")
		parts := strings.SplitN(path, "/", 2)
		if len(parts) != 2 {
			s3base.WriteError(ctx, http.StatusInternalServerError, errs.New(errs.ErrParam, "part not match"))
			return
		}
		bucket := parts[0]
		obj := parts[1]
		if len(bucket) == 0 || len(obj) == 0 {
			s3base.WriteError(ctx, http.StatusInternalServerError, errs.New(errs.ErrParam, "bucket or obj len is 0"))
			return
		}
		s3base.SetS3Bucket(ctx, bucket)
		s3base.SetS3Object(ctx, obj)
	}
}
