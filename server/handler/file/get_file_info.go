package file

import (
	"context"
	"fileserver/proxyutil"
	"fileserver/server/model"
	"fileserver/service"
	"fileserver/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMetaInfo(c *gin.Context, ctx context.Context, request interface{}) {
	req := request.(*model.GetFileInfoRequest)
	if len(req.Key) == 0 {
		proxyutil.Fail(c, http.StatusBadRequest, fmt.Errorf("empty key"))
		return
	}
	fileid, err := utils.DecodeFileId(req.Key)
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
			Key:           req.Key,
			Exist:         true,
			FileName:      item.FileName,
			FileSize:      item.FileSize,
			Ctime:         item.Ctime,
			Mtime:         item.Mtime,
			FilePartCount: item.FilePartCount,
		},
	})
}
