package file

import (
	"context"
	"fileserver/constants"
	"fileserver/dao"
	"fileserver/model"
	"fileserver/proto/fileserver/fileinfo"
	"fileserver/utils"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/cache"
	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/logutil"
	"github.com/xxxsen/common/naivesvr"
	"github.com/xxxsen/common/naivesvr/codec"
	"github.com/xxxsen/common/s3"
	"github.com/yitter/idgenerator-go/idgen"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var fileCache, _ = cache.NewLocalCache(20000)

var ImageUpload = CommonFilePostUpload(NewFileUploader(uint32(model.FileTypeImage), ImageExtChecker))
var VideoUpload = CommonFilePostUpload(NewFileUploader(uint32(model.FileTypeVideo), VideoExtChecker))
var FileUpload = CommonFilePostUpload(NewFileUploader(uint32(model.FileTypeFile), nil))
var FileDownload = CommonFileDownload(NewFileDownloader(fileCache))

type TypeCheckFunc func(meta *FileUploadMeta) error

func ExtNameChecker(exts ...string) TypeCheckFunc {
	valid := map[string]interface{}{}
	for _, ext := range exts {
		valid[strings.ToLower(ext)] = true
	}
	return func(meta *FileUploadMeta) error {
		ext := strings.ToLower(filepath.Ext(meta.FileName))
		if _, ok := valid[ext]; ok {
			return nil
		}
		return errs.New(errs.ErrParam, "not support ext:%s", ext)
	}
}

var ImageExtChecker = ExtNameChecker(".jpg", ".png")
var VideoExtChecker = ExtNameChecker(".mp4")

type FileUploadMeta struct {
	Reader   io.ReadSeekCloser
	FileName string
	DownKey  string
	FileSize int64
	MD5      string
}

type ISmallFileUploader interface {
	BeforeUpload(ctx *gin.Context, request interface{}) (*FileUploadMeta, bool, error)
	OnUpload(ctx *gin.Context, meta *FileUploadMeta) error
	AfterUpload(ctx *gin.Context, realUpload bool, meta *FileUploadMeta) (interface{}, error)
}

type S3SmallFileUploader struct {
}

type BasicFileUploadRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
	MD5  string                `form:"md5" binding:"required"`
}

func (f *S3SmallFileUploader) BeforeUpload(ctx *gin.Context, request interface{}) (*FileUploadMeta, bool, error) {
	req := request.(*BasicFileUploadRequest)
	header := req.File
	file, err := header.Open()
	if err != nil {
		return nil, false, errs.Wrap(errs.ErrParam, "open file fail", err)
	}
	md5 := req.MD5
	if header.Size > constants.MaxPostUploadSize {
		file.Close()
		return nil, false, errs.New(errs.ErrParam, "file size out of limit")
	}
	return &FileUploadMeta{
		Reader:   file,
		FileName: header.Filename,
		MD5:      md5,
		DownKey:  utils.EncodeFileId(uint64(idgen.NextId())),
		FileSize: header.Size,
	}, true, nil
}

func (f *S3SmallFileUploader) OnUpload(ctx *gin.Context, meta *FileUploadMeta) error {
	if err := s3.Client.Upload(ctx, meta.DownKey, meta.Reader, meta.FileSize, meta.MD5); err != nil {
		return errs.Wrap(errs.ErrS3, "upload to s3 fail", err)
	}
	return nil
}

func (f *S3SmallFileUploader) AfterUpload(ctx *gin.Context, realUpload bool, meta *FileUploadMeta) (interface{}, error) {
	return nil, nil
}

type FileDownloadMeta struct {
	DownKey     string
	FileName    string
	FileSize    int64
	ContentType string
}

type IFileDownloader interface {
	BeforeDownload(ctx *gin.Context, request interface{}) (*FileDownloadMeta, error)
	OnDownload(ctx *gin.Context, meta *FileDownloadMeta) error
	AfterDownload(ctx *gin.Context, meta *FileDownloadMeta, err error)
}

type S3FileDownloader struct {
}

func (f *S3FileDownloader) BeforeDownload(ctx *gin.Context, request interface{}) (*FileDownloadMeta, error) {
	return nil, fmt.Errorf("need impl")
}

func (f *S3FileDownloader) OnDownload(ctx *gin.Context, meta *FileDownloadMeta) error {
	reader, err := s3.Client.Download(ctx, meta.DownKey)
	if err != nil {
		return errs.Wrap(errs.ErrS3, "create download stream fail", err)
	}
	defer reader.Close()
	contentType := meta.ContentType
	if len(contentType) == 0 {
		contentType = mime.TypeByExtension(filepath.Ext(meta.FileName))
	}
	writer := ctx.Writer
	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", strconv.Quote(meta.FileName)))
	writer.Header().Set("Content-Length", fmt.Sprintf("%d", meta.FileSize))
	writer.Header().Set("Content-Type", contentType)
	sz, err := io.Copy(ctx.Writer, reader)
	if err != nil {
		return errs.Wrap(errs.ErrIO, "copy stream fail", err)
	}
	if sz != int64(meta.FileSize) {
		return errs.New(errs.ErrIO, "io size not match, need %d, write:%d", meta.FileSize, sz)
	}
	return nil
}

func (f *S3FileDownloader) AfterDownload(ctx *gin.Context, meta *FileDownloadMeta, err error) {
	if err == nil {
		return
	}
	logutil.GetLogger(ctx).With(zap.String("path", ctx.Request.URL.Path), zap.Error(err)).Error("file download fail")
}

func CommonFilePostUpload(uploader ISmallFileUploader) naivesvr.ProcessFunc {
	return func(ctx *gin.Context, req interface{}) (int, errs.IError, interface{}) {
		caller := func() (interface{}, errs.IError) {
			meta, needUpload, err := uploader.BeforeUpload(ctx, req)
			defer func() {
				if meta != nil && meta.Reader != nil {
					meta.Reader.Close()
				}
			}()
			if err != nil {
				return nil, errs.Wrap(errs.ErrServiceInternal, "before post upload fail", err)
			}
			if needUpload {
				if err := uploader.OnUpload(ctx, meta); err != nil {
					return nil, errs.Wrap(errs.ErrStorage, "on upload fail", err)
				}
			}
			rsp, err := uploader.AfterUpload(ctx, needUpload, meta)
			if err != nil {
				return nil, errs.Wrap(errs.ErrServiceInternal, "after upload fail", err)
			}
			return rsp, nil
		}
		rsp, err := caller()
		return http.StatusOK, err, rsp
	}
}

func CommonFileDownload(downloader IFileDownloader) naivesvr.ProcessFunc {
	return func(ctx *gin.Context, request interface{}) (int, errs.IError, interface{}) {
		caller := func() error {
			meta, err := downloader.BeforeDownload(ctx, request)
			if err == nil {
				err = downloader.OnDownload(ctx, meta)
			}
			downloader.AfterDownload(ctx, meta, err)
			if err != nil {
				return err
			}
			return nil
		}
		if err := caller(); err != nil {
			e := errs.FromError(err)
			logutil.GetLogger(ctx).With(
				zap.String("path", ctx.Request.URL.Path), zap.Error(e),
			).Error("call file download fail")
			codec.JsonCodec.Encode(ctx, http.StatusOK, e, nil)
		}
		return http.StatusOK, nil, nil
	}
}

func cacheGetFileMeta(ctx context.Context, c cache.ICache, key interface{},
	cb func() (interface{}, bool, error)) (interface{}, bool, error) {

	ival, exist, _ := c.Get(ctx, key)
	if exist {
		return ival, true, nil
	}
	val, exist, err := cb()
	if err != nil {
		return nil, false, err
	}
	if exist {
		c.Set(ctx, key, val, 10*time.Minute)
	}
	return val, exist, nil
}

type FileUploader struct {
	S3SmallFileUploader
	typ  uint32
	ckfn TypeCheckFunc
}

func NewFileUploader(typ uint32, ckfn TypeCheckFunc) *FileUploader {
	if ckfn == nil {
		ckfn = func(meta *FileUploadMeta) error { return nil }
	}
	return &FileUploader{
		typ:  typ,
		ckfn: ckfn,
	}
}

func (uploader *FileUploader) BeforeUpload(ctx *gin.Context, request interface{}) (*FileUploadMeta, bool, error) {
	meta, _, err := uploader.S3SmallFileUploader.BeforeUpload(ctx, request)
	if err != nil {
		return nil, false, err
	}
	meta.DownKey = fmt.Sprintf("%d_%s", uploader.typ, meta.DownKey)
	if err := uploader.ckfn(meta); err != nil {
		return nil, false, errs.Wrap(errs.ErrParam, "meta check not pass", err)
	}
	return meta, true, nil
}

func (uploader *FileUploader) AfterUpload(ctx *gin.Context, realUpload bool, meta *FileUploadMeta) (interface{}, error) {
	if _, err := dao.FileInfoDao.CreateFile(ctx, &model.CreateFileRequest{
		Item: &model.FileItem{
			FileName:   meta.FileName,
			Hash:       meta.MD5,
			FileSize:   uint64(meta.FileSize),
			CreateTime: uint64(time.Now().UnixMilli()),
			DownKey:    meta.DownKey,
		},
	}); err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "insert image to db fail", err)
	}
	return &fileinfo.FileUploadResponse{
		DownKey: proto.String(meta.DownKey),
	}, nil
}

type BasicFileDownloadRequest struct {
	DownKey string `form:"down_key" binding:"required"`
}

type FileDownloader struct {
	S3FileDownloader
	c cache.ICache
}

func NewFileDownloader(c cache.ICache) *FileDownloader {
	return &FileDownloader{c: c}
}

func (d *FileDownloader) BeforeDownload(ctx *gin.Context, request interface{}) (*FileDownloadMeta, error) {
	req := request.(*BasicFileDownloadRequest)
	downKey := req.DownKey
	ifileinfo, exist, err := cacheGetFileMeta(ctx, d.c, downKey, func() (interface{}, bool, error) {
		daoRsp, exist, err := dao.FileInfoDao.GetFile(ctx, &model.GetFileRequest{
			DownKey: downKey,
		})
		if err != nil {
			return nil, false, err
		}
		if !exist {
			return nil, false, nil
		}
		return daoRsp.Item, true, nil
	})
	if err != nil {
		return nil, errs.Wrap(errs.ErrStorage, "cache get file meta fail", err)
	}
	if !exist {
		return nil, errs.New(errs.ErrNotFound, "not found file meta")
	}
	fileinfo := ifileinfo.(*model.FileItem)
	return &FileDownloadMeta{
		DownKey:  downKey,
		FileName: fileinfo.FileName,
		FileSize: int64(fileinfo.FileSize),
	}, nil
}
