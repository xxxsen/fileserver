package handler

import (
	"fileserver/handler/file"
	"fileserver/handler/file/bigfile"
	"fileserver/handler/middlewares"
	"fileserver/handler/s3"
	"fileserver/proto/fileserver/fileinfo"

	"github.com/xxxsen/common/naivesvr"
	"github.com/xxxsen/common/naivesvr/codec"

	"github.com/gin-gonic/gin"
)

type RegistConfig struct {
	User string
	Pwd  string
}

func OnRegistWithConfig(opts ...Option) func(router *gin.Engine) {
	return func(router *gin.Engine) {
		OnRegist(router, opts...)
	}
}

func OnRegist(router *gin.Engine, opts ...Option) {
	c := &config{
		users: make(map[string]string),
	}
	for _, opt := range opts {
		opt(c)
	}

	authMiddleware := middlewares.CommonAuth(c.users)

	uploadLimitMiddleware := middlewares.IOThreadLimitMiddleware(middlewares.NewIOThreadContext(int64(c.maxUploadThread)))
	downloadLimitMiddleware := middlewares.IOThreadLimitMiddleware(middlewares.NewIOThreadContext(int64(c.maxDownloadThread)))

	//upload
	{
		uploadRouter := router.Group("/upload", authMiddleware, uploadLimitMiddleware)
		uploadRouter.POST("/image", naivesvr.WrapHandler(&file.BasicFileUploadRequest{}, codec.CustomCodec(codec.JsonCodec, codec.MultipartCodec), file.ImageUpload))
		uploadRouter.POST("/video", naivesvr.WrapHandler(&file.BasicFileUploadRequest{}, codec.CustomCodec(codec.JsonCodec, codec.MultipartCodec), file.VideoUpload))
		uploadRouter.POST("/file", naivesvr.WrapHandler(&file.BasicFileUploadRequest{}, codec.CustomCodec(codec.JsonCodec, codec.MultipartCodec), file.FileUpload))
		bigFileRouter := uploadRouter.Group("/bigfile")
		bigFileRouter.POST("/begin", naivesvr.WrapHandler(&fileinfo.FileUploadBeginRequest{}, codec.JsonCodec, bigfile.Begin))
		bigFileRouter.POST("/part", naivesvr.WrapHandler(&bigfile.PartUploadRequest{}, codec.CustomCodec(codec.JsonCodec, codec.MultipartCodec), bigfile.Part))
		bigFileRouter.POST("/end", naivesvr.WrapHandler(&fileinfo.FileUploadEndRequest{}, codec.JsonCodec, bigfile.End))

	}
	//download
	{
		router.GET("/file", downloadLimitMiddleware, naivesvr.WrapHandler(&file.BasicFileDownloadRequest{}, codec.CustomCodec(codec.NopCodec, codec.QueryCodec), file.FileDownload)) //input: down_key
	}
	//meta
	{
		router.POST("/file/meta", naivesvr.WrapHandler(&fileinfo.GetFileMetaRequest{}, codec.JsonCodec, file.Meta))
	}
	//s3
	{
		//
		s3Prefix := "/s3/"
		router.GET("/s3/*s3Param", downloadLimitMiddleware, middlewares.S3BucketOpLimitMiddleware(s3Prefix), s3.Download)
		router.PUT("/s3/*s3Param", authMiddleware, uploadLimitMiddleware, middlewares.S3BucketOpLimitMiddleware(s3Prefix), s3.Upload)
	}
}
