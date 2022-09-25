package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	keyS3Bucket = "x-s3-bucket"
	keyS3Object = "x-s3-object"
)

func GetS3Bucket(ctx *gin.Context) (string, bool) {
	val, ok := ctx.Get(keyS3Bucket)
	if !ok {
		return "", false
	}
	return val.(string), true
}

func GetS3Object(ctx *gin.Context) (string, bool) {
	val, ok := ctx.Get(keyS3Object)
	if !ok {
		return "", false
	}
	return val.(string), true
}

func setS3Bucket(ctx *gin.Context, bk string) {
	ctx.Set(keyS3Bucket, bk)
}

func setS3Object(ctx *gin.Context, obj string) {
	ctx.Set(keyS3Object, obj)
}

func S3BucketOpLimitMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		switch ctx.Request.Method {
		case http.MethodGet, http.MethodPut:
		default:
			http.Error(ctx.Writer, "unsupport method", http.StatusInternalServerError)
			return
		}
		path := ctx.Request.URL.Path
		path = strings.TrimLeft(path, "/s3")
		path = strings.Trim(path, "/")
		parts := strings.SplitN(path, "/", 2)
		if len(parts) != 2 {
			http.Error(ctx.Writer, "part not match", http.StatusInternalServerError)
			return
		}
		bucket := parts[0]
		obj := parts[1]
		if len(bucket) == 0 || len(obj) == 0 {
			http.Error(ctx.Writer, "bucket or obj len is 0", http.StatusInternalServerError)
			return
		}
		setS3Bucket(ctx, bucket)
		setS3Object(ctx, obj)
	}
}
