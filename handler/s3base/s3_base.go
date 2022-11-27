package s3base

import (
	"encoding/xml"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/logutil"
	"github.com/xxxsen/common/trace"
	"go.uber.org/zap"
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

func SetS3Bucket(ctx *gin.Context, bk string) {
	ctx.Set(keyS3Bucket, bk)
}

func SetS3Object(ctx *gin.Context, obj string) {
	ctx.Set(keyS3Object, obj)
}

type S3ErrorMessage struct {
	XMLName    xml.Name `xml:"Error"`
	Code       string   `xml:"Code"`
	Message    string   `xml:"Message"`
	Key        string   `xml:"Key"`
	BucketName string   `xml:"BucketName"`
	Resouce    string   `xml:"Resource"`
	RequestId  string   `xml:"RequestId"`
	HostId     string   `xml:"HostId"`
}

func ResponseWithError(ctx *gin.Context, code int, e *S3ErrorMessage) {
	ctx.XML(code, e)
}

func SimpleReply(ctx *gin.Context) {
	data := []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>" +
		"<LocationConstraint xmlns=\"http://s3.amazonaws.com/doc/2006-03-01/\"></LocationConstraint>")
	_, err := ctx.Writer.Write(data)
	if err != nil {
		logutil.GetLogger(ctx).Error("write msg fail", zap.Error(err))
		return
	}
}

func WriteError(ctx *gin.Context, statuscode int, xerr error) {
	bucket, _ := GetS3Bucket(ctx)
	obj, _ := GetS3Object(ctx)
	logutil.GetLogger(ctx).Error("write err to client",
		zap.Int("status_code", statuscode),
		zap.String("bucket", bucket),
		zap.String("obj", obj),
		zap.Error(xerr))
	traceid, _ := trace.GetTraceId(ctx)
	err := errs.FromError(xerr)
	e := &S3ErrorMessage{
		Code:       fmt.Sprintf("%d", err.Code()),
		Message:    err.Message(),
		Key:        obj,
		BucketName: bucket,
		Resouce:    fmt.Sprintf("%s/%s", bucket, obj),
		RequestId:  traceid,
		HostId:     traceid,
	}
	ResponseWithError(ctx, statuscode, e)
}
