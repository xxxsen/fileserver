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

const (
	defaultMaxListMetaSizePerRequest = 20
)

func GetMetaInfo(c *gin.Context, ctx context.Context, request interface{}) {
	req := request.(*model.GetFileInfoRequest)
	if len(req.DownKey) == 0 || len(req.DownKey) > defaultMaxListMetaSizePerRequest {
		proxyutil.Fail(c, http.StatusBadRequest, fmt.Errorf("invalid down key size:%d", len(req.DownKey)))
		return
	}
	fileids := make([]uint64, 0, len(req.DownKey))
	for _, item := range req.DownKey {
		fileid, err := utils.DecodeFileId(item)
		if err != nil {
			proxyutil.Fail(c, http.StatusBadRequest, fmt.Errorf("decode down key fail, err:%w", err))
			return
		}
		fileids = append(fileids, fileid)
	}
	rs, err := service.FileService.BatchGetFileInfo(ctx, fileids)
	if err != nil {
		proxyutil.Fail(c, http.StatusInternalServerError, fmt.Errorf("read file info fail, err:%w", err))
		return
	}
	rsp := &model.GetFileInfoResponse{List: make([]*model.GetFileInfoItem, 0, len(req.DownKey))}
	for idx, fileid := range fileids {
		item, ok := rs[fileid]
		if !ok {
			continue
		}
		rsp.List = append(rsp.List, &model.GetFileInfoItem{
			DownKey:       req.DownKey[idx],
			FileName:      item.FileName,
			FileSize:      item.FileSize,
			Ctime:         item.Ctime,
			Mtime:         item.Mtime,
			FilePartCount: item.FilePartCount,
		})
	}
	proxyutil.Success(c, rsp)
}
