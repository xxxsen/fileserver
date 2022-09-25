package s3

import (
	"encoding/xml"
	"fileserver/handler/middlewares"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/logutil"
	"github.com/xxxsen/common/trace"
	"go.uber.org/zap"
)

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

func WriteError(ctx *gin.Context, statuscode int, err errs.IError) {
	bucket, _ := middlewares.GetS3Bucket(ctx)
	obj, _ := middlewares.GetS3Object(ctx)
	logutil.GetLogger(ctx).With(
		zap.Int("status_code", statuscode),
		zap.String("bucket", bucket),
		zap.String("obj", obj),
		zap.Error(err),
	).Error("write err to client")
	traceid, _ := trace.GetTraceId(ctx)
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
