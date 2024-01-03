package handler

import (
	"fileserver/handler/file"
	"fileserver/handler/file/bigfile"
	"fileserver/handler/middlewares"
	"fileserver/handler/s3"
	"fileserver/proto/fileserver/fileinfo"
	"fmt"

	"github.com/xxxsen/common/cgi"
	"github.com/xxxsen/common/cgi/codec"

	"github.com/gin-gonic/gin"
)

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

	refererMiddleware := middlewares.RefererMiddleware(c.enableRefererCheck, c.referers)

	router.Use(middlewares.TimeCostMiddleware())
	router.Use(refererMiddleware)
	//upload
	{
		uploadRouter := router.Group("/upload", authMiddleware, uploadLimitMiddleware)
		uploadRouter.POST("/image", cgi.WrapHandler(&file.BasicFileUploadRequest{}, codec.CustomCodec(codec.JsonCodec, codec.MultipartCodec), file.ImageUpload))
		uploadRouter.POST("/video", cgi.WrapHandler(&file.BasicFileUploadRequest{}, codec.CustomCodec(codec.JsonCodec, codec.MultipartCodec), file.VideoUpload))
		uploadRouter.POST("/file", cgi.WrapHandler(&file.BasicFileUploadRequest{}, codec.CustomCodec(codec.JsonCodec, codec.MultipartCodec), file.FileUpload))
		bigFileRouter := uploadRouter.Group("/bigfile")
		bigFileRouter.POST("/begin", cgi.WrapHandler(&fileinfo.FileUploadBeginRequest{}, codec.JsonCodec, bigfile.Begin))
		bigFileRouter.POST("/part", cgi.WrapHandler(&bigfile.PartUploadRequest{}, codec.CustomCodec(codec.JsonCodec, codec.MultipartCodec), bigfile.Part))
		bigFileRouter.POST("/end", cgi.WrapHandler(&fileinfo.FileUploadEndRequest{}, codec.JsonCodec, bigfile.End))

	}
	//download
	{
		router.GET("/file", downloadLimitMiddleware, cgi.WrapHandler(&file.BasicFileDownloadRequest{}, codec.CustomCodec(codec.NopCodec, codec.QueryCodec), file.FileDownload)) //input: down_key
	}
	//meta
	{
		router.POST("/file/meta", cgi.WrapHandler(&fileinfo.GetFileMetaRequest{}, codec.JsonCodec, file.Meta))
	}
	registS3(router, c, authMiddleware, downloadLimitMiddleware, uploadLimitMiddleware)
}

func registS3(router *gin.Engine, c *config, authMiddleware, downloadLimitMiddleware, uploadLimitMiddleware gin.HandlerFunc) {
	if !c.enableFakeS3 {
		return
	}
	for _, bk := range c.fakeS3Buckets {
		bucketPath := fmt.Sprintf("/%s", bk)
		routerPath := fmt.Sprintf("%s/*s3Param", bucketPath)
		router.GET(bucketPath, middlewares.S3BucketOpLimitMiddleware(), s3.GetBucket)
		router.GET(routerPath, downloadLimitMiddleware, middlewares.S3BucketOpLimitMiddleware(), s3.Download)
		router.PUT(routerPath, authMiddleware, uploadLimitMiddleware, middlewares.S3BucketOpLimitMiddleware(), s3.Upload)
	}
}
