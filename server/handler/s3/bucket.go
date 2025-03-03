package s3

import (
	"fileserver/server/handler/s3/s3base"

	"github.com/gin-gonic/gin"
)

func GetBucket(c *gin.Context) {
	s3base.SimpleReply(c)
}
