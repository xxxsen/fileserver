package file

import (
	"context"
	"fileserver/dao"
	"fileserver/model"
	"fileserver/proto/fileserver/fileinfo"
	"fileserver/proxyutil"
	"fileserver/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

const (
	maxListMetaSizePerRequest = 20
)

func GetMetaInfo(c *gin.Context, ctx context.Context, request interface{}) {
	req := request.(*fileinfo.GetFileMetaRequest)
	if len(req.DownKey) == 0 || len(req.DownKey) > maxListMetaSizePerRequest {
		proxyutil.Fail(c, http.StatusBadRequest, fmt.Errorf("invalid down key size:%d", len(req.DownKey)))
		return
	}
	downkeys := make([]uint64, 0, len(req.GetDownKey()))
	for _, item := range req.GetDownKey() {
		downkey, err := utils.DecodeFileId(item)
		if err != nil {
			proxyutil.Fail(c, http.StatusBadRequest, fmt.Errorf("decode down key fail, err:%w", err))
			return
		}
		downkeys = append(downkeys, downkey)
	}
	daoRsp, err := dao.FileInfoDao.ListFile(ctx, &model.ListFileRequest{
		Query: &model.ListFileQuery{
			DownKey: downkeys,
		},
		Offset: 0,
		Limit:  uint32(len(req.DownKey)),
	})
	if err != nil {
		proxyutil.Fail(c, http.StatusInternalServerError, fmt.Errorf("read file list fail, err:%w", err))
		return
	}
	metalist := fileinfo2pbmeta(req.GetDownKey(), downkeys, daoRsp.List)
	proxyutil.Success(c, &fileinfo.GetFileMetaResponse{
		List: metalist,
	})
}

func fileinfo2pbmeta(origin []string, decoded []uint64, lst []*model.FileItem) []*fileinfo.FileItem {
	mapper := make(map[uint64]*model.FileItem)
	for _, item := range lst {
		mapper[item.DownKey] = item
	}
	rs := make([]*fileinfo.FileItem, 0, len(lst))
	for index, key := range decoded {
		src := mapper[key]
		dst := &fileinfo.FileItem{
			DownKey: proto.String(origin[index]),
			Exist:   proto.Bool(false),
		}
		if src != nil {
			dst.FileName = proto.String(src.FileName)
			dst.Hash = proto.String(src.Hash)
			dst.FileSize = proto.Uint64(src.FileSize)
			dst.CreateTime = proto.Uint64(src.CreateTime)
			dst.Exist = proto.Bool(true)
		}
		rs = append(rs, dst)
	}
	return rs
}
