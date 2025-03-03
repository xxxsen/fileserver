package s3

import (
	"bytes"
	"fileserver/core"
	"fileserver/dao"
	"fileserver/handler/common"
	"fileserver/handler/s3base"
	"fileserver/model"
	"fileserver/utils"
	"fmt"
	"io"
	"net/http"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/idgen"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

func bigFileUpload(ctx *gin.Context) (uint64, string, error) {
	length := ctx.Request.ContentLength
	fs := core.GetFsCore()
	//生成uploadid
	beginRsp, err := fs.BeginFileUpload(ctx, &core.BeginFileUploadRequest{
		FileSize: length,
	})
	if err != nil {
		return 0, "", fmt.Errorf("begin upload fail, err:%w", err)
	}
	file := ctx.Request.Body
	uploadid := beginRsp.UploadID

	//分批上传
	blkcnt := utils.CalcFileBlockCount(uint64(length), uint64(fs.BlockSize()))
	for i := 0; i < blkcnt; i++ {
		partid := i + 1
		r := io.LimitReader(file, fs.BlockSize())
		raw, err := io.ReadAll(r)
		if err != nil {
			return 0, "", fmt.Errorf("read io data fail, err:%w", err)
		}
		md5v := utils.GetMd5(raw)

		_, err = fs.PartFileUpload(ctx, &core.PartFileUploadRequest{
			ReadSeeker: bytes.NewReader(raw),
			UploadId:   uploadid,
			PartId:     uint64(partid),
			Size:       int64(len(raw)),
			MD5:        md5v,
		})
		if err != nil {
			return 0, "", fmt.Errorf("part upload fail, err:%w", err)
		}
	}
	obj, _ := s3base.GetS3Object(ctx)
	name := path.Base(obj)
	//完成上传
	rsp, err := fs.FinishFileUpload(ctx, &core.FinishFileUploadRequest{
		UploadId: uploadid,
		FileName: name,
	})
	if err != nil {
		return 0, "", fmt.Errorf("finish file upload fail, err:%w", err)
	}
	//写入db
	fileid := idgen.NextId()
	_, err = dao.FileInfoDao.CreateFile(ctx, &model.CreateFileRequest{
		Item: &model.FileItem{
			FileName:   name,
			Hash:       rsp.CheckSum,
			FileSize:   uint64(rsp.FileSize),
			CreateTime: uint64(time.Now().UnixMilli()),
			FileKey:    rsp.Key,
			Extra:      rsp.Extra,
			DownKey:    fileid,
			StType:     fs.StType(),
		},
	})
	if err != nil {
		return 0, "", fmt.Errorf("write file to db fail, err:%w", err)
	}
	return fileid, rsp.CheckSum, nil
}

func smallFileUpload(ctx *gin.Context) (uint64, string, error) {
	md5Base64 := ctx.Request.Header.Get("Content-MD5")
	length := ctx.Request.ContentLength

	raw, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return 0, "", fmt.Errorf("read body fail, err:%w", err)
	}
	checksum, err := utils.Base64Md52HexMd5(md5Base64)
	if err != nil {
		return 0, "", fmt.Errorf("invalid checksum, err:%w", err)
	}
	if len(checksum) == 0 {
		checksum = utils.GetMd5(raw)
	}

	obj, _ := s3base.GetS3Object(ctx)
	name := path.Base(obj)
	uploadRequest := common.CommonUploadContext{
		Reader: bytes.NewReader(raw),
		Size:   length,
		Name:   name,
		Md5Sum: checksum,
	}
	fileid, err := common.Upload(ctx, &uploadRequest)
	if err != nil {
		return 0, "", fmt.Errorf("upload fail, err:%w", err)
	}
	return fileid, checksum, nil
}

func Upload(ctx *gin.Context) {
	caller := smallFileUpload
	length := ctx.Request.ContentLength
	if length > core.GetFsCore().BlockSize() {
		caller = bigFileUpload
	}
	fileid, checksum, err := caller(ctx)
	if err != nil {
		s3base.WriteError(ctx, http.StatusInternalServerError, fmt.Errorf("do file upload fail, err:%w", err))
		return
	}
	bucket, _ := s3base.GetS3Bucket(ctx)
	obj, _ := s3base.GetS3Object(ctx)
	filename := fmt.Sprintf("%s/%s", bucket, obj)
	if _, err := dao.MappingInfoDao.CreateMappingInfo(ctx, &model.CreateMappingInfoRequest{
		Item: &model.MappingInfoItem{
			FileName: filename,
			FileId:   fileid,
		},
	}); err != nil {
		s3base.WriteError(ctx, http.StatusInternalServerError, fmt.Errorf("create mapping fail, err:%w", err))
		return
	}
	ctx.Writer.Header().Set("ETag", `"`+checksum+`"`)
	ctx.Writer.WriteHeader(http.StatusOK)
	logutil.GetLogger(ctx).Info("upload file finish", zap.Int64("size", length), zap.String("bucket", bucket), zap.String("obj", obj))
}
