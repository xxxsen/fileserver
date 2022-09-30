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

func OnRegistWithConfig(c *RegistConfig) func(router *gin.Engine) {
	return func(router *gin.Engine) {
		OnRegist(router, c)
	}
}

func OnRegist(router *gin.Engine, c *RegistConfig) {
	commonAuth := middlewares.CommonAuth(c.User, c.Pwd)

	//upload
	{
		uploadRouter := router.Group("/upload")
		uploadRouter.POST("/image", commonAuth, naivesvr.WrapHandler(&file.BasicFileUploadRequest{}, codec.CustomCodec(codec.JsonCodec, codec.MultipartCodec), file.ImageUpload))
		uploadRouter.POST("/video", commonAuth, naivesvr.WrapHandler(&file.BasicFileUploadRequest{}, codec.CustomCodec(codec.JsonCodec, codec.MultipartCodec), file.VideoUpload))
		uploadRouter.POST("/file", commonAuth, naivesvr.WrapHandler(&file.BasicFileUploadRequest{}, codec.CustomCodec(codec.JsonCodec, codec.MultipartCodec), file.FileUpload))
		bigFileRouter := uploadRouter.Group("/bigfile")
		bigFileRouter.POST("/begin", commonAuth, naivesvr.WrapHandler(&fileinfo.FileUploadBeginRequest{}, codec.JsonCodec, bigfile.Begin))
		bigFileRouter.POST("/part", commonAuth, naivesvr.WrapHandler(&bigfile.PartUploadRequest{}, codec.CustomCodec(codec.JsonCodec, codec.MultipartCodec), bigfile.Part))
		bigFileRouter.POST("/end", commonAuth, naivesvr.WrapHandler(&fileinfo.FileUploadEndRequest{}, codec.JsonCodec, bigfile.End))

	}
	//download
	{
		router.GET("/file", naivesvr.WrapHandler(&file.BasicFileDownloadRequest{}, codec.CustomCodec(codec.NopCodec, codec.QueryCodec), file.FileDownload)) //input: down_key
	}
	//meta
	{
		router.POST("/file/meta", naivesvr.WrapHandler(&fileinfo.GetFileMetaRequest{}, codec.JsonCodec, file.Meta))
	}
	//s3
	{
		//
		s3Prefix := "/s3/"
		router.GET("/s3/*s3Param", middlewares.S3BucketOpLimitMiddleware(s3Prefix), s3.Download)
		router.PUT("/s3/*s3Param", middlewares.S3BucketOpLimitMiddleware(s3Prefix), s3.Upload)
	}
}
