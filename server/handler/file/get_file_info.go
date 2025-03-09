package file

import (
	"fmt"
	"net/http"
	"tgfile/proxyutil"
	"tgfile/server/model"
	"tgfile/service"
	"tgfile/utils"

	"github.com/gin-gonic/gin"
)

func GetMetaInfo(c *gin.Context) {
	ctx := c.Request.Context()
	key := c.Param("key")
	fileid, err := utils.DecodeFileId(key)
	if err != nil {
		proxyutil.Fail(c, http.StatusBadRequest, fmt.Errorf("decode down key fail, err:%w", err))
		return
	}
	rsMap, err := service.FileService.BatchGetFileInfo(ctx, []uint64{fileid})
	if err != nil {
		proxyutil.Fail(c, http.StatusInternalServerError, fmt.Errorf("read file info fail, err:%w", err))
		return
	}
	item, ok := rsMap[fileid]
	if !ok {
		proxyutil.Success(c, &model.GetFileInfoResponse{
			Item: &model.GetFileInfoItem{
				Exist: false,
			},
		})
		return
	}
	proxyutil.Success(c, &model.GetFileInfoResponse{
		Item: &model.GetFileInfoItem{
			Key:           key,
			Exist:         true,
			FileName:      item.FileName,
			FileSize:      item.FileSize,
			Ctime:         item.Ctime,
			Mtime:         item.Mtime,
			FilePartCount: item.FilePartCount,
		},
	})
}
