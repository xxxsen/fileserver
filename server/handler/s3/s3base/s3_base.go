package s3base

import (
	"encoding/xml"

	"github.com/gin-gonic/gin"
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

func SimpleReply(ctx *gin.Context) {
	data := []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>" +
		"<LocationConstraint xmlns=\"http://s3.amazonaws.com/doc/2006-03-01/\"></LocationConstraint>")
	_, err := ctx.Writer.Write(data)
	if err != nil {
		logutil.GetLogger(ctx).Error("write msg fail", zap.Error(err))
		return
	}
}

func WriteError(c *gin.Context, statuscode int, err error) {
	ctx := c.Request.Context()
	logutil.GetLogger(ctx).Error("write err to client",
		zap.Error(err),
		zap.Int("status_code", statuscode))
	traceid, _ := trace.GetTraceId(ctx)
	e := &S3ErrorMessage{
		Code:      "500",
		Message:   err.Error(),
		Resouce:   c.Request.URL.Path,
		RequestId: traceid,
		HostId:    traceid,
	}
	ResponseWithError(c, statuscode, e)
}
