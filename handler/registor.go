package handler

import (
	"fileserver/handler/file"
	"fileserver/handler/file/bigfile"
	"fileserver/proto/fileserver/fileinfo"

	"github.com/xxxsen/common/naivesvr"
	"github.com/xxxsen/common/naivesvr/codec"

	"github.com/gin-gonic/gin"
)

func OnRegist(router *gin.Engine) {
	//upload
	{
		uploadRouter := router.Group("/upload")
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
		router.GET("/file", naivesvr.WrapHandler(&file.BasicFileDownloadRequest{}, codec.CustomCodec(codec.NopCodec, codec.QueryCodec), file.FileDownload)) //input: down_key
	}
	//meta
	{
		router.POST("/file/meta", naivesvr.WrapHandler(&fileinfo.GetFileMetaRequest{}, codec.JsonCodec, file.Meta))
	}
}
