package file

import (
	"fileserver/dao"
	"fileserver/model"
	"fileserver/proto/fileserver/fileinfo"
	"fileserver/utils"
	"net/http"

	"github.com/xxxsen/common/errs"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

const (
	maxListMetaSizePerRequest = 20
)

func Meta(ctx *gin.Context, request interface{}) (int, interface{}, error) {
	req := request.(*fileinfo.GetFileMetaRequest)
	if len(req.DownKey) == 0 || len(req.DownKey) > maxListMetaSizePerRequest {
		return http.StatusOK, nil, errs.New(errs.ErrParam, "invalid down key size:%d", len(req.DownKey))
	}
	downkeys := make([]uint64, 0, len(req.GetDownKey()))
	for _, item := range req.GetDownKey() {
		downkey, err := utils.DecodeFileId(item)
		if err != nil {
			return http.StatusOK, nil, errs.Wrap(errs.ErrParam, "decode down key fail", err)
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
		return http.StatusOK, errs.Wrap(errs.ErrDatabase, "read file list fail", err), nil
	}
	metalist := fileinfo2pbmeta(req.GetDownKey(), downkeys, daoRsp.List)
	return http.StatusOK, &fileinfo.GetFileMetaResponse{
		List: metalist,
	}, nil
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
