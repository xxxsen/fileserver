package server

import (
	"fileserver/handler/file"
	"fileserver/handler/middleware"
	"fileserver/handler/s3"
	"fileserver/proto/fileserver/fileinfo"
	"fileserver/proxyutil"
	"fmt"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type Server struct {
	addr   string
	c      *config
	engine *gin.Engine
}

func New(bind string, opts ...Option) (*Server, error) {
	c := applyOpts(opts...)
	svr := &Server{addr: bind, c: c}
	if err := svr.init(); err != nil {
		return nil, err
	}
	return svr, nil
}

func (s *Server) init() error {
	s.engine = gin.New()
	s.initBasic(s.engine)
	s.initS3(s.engine)
	return nil
}

func (s *Server) initBasic(router *gin.Engine) {
	authMiddleware := middleware.CommonAuth(s.c.userMap)
	router.POST("/upload/file", authMiddleware, proxyutil.WrapBizFunc(file.FileUpload, &file.BasicFileUploadRequest{}))
	router.GET("/file", proxyutil.WrapBizFunc(file.FileDownload, &file.BasicFileDownloadRequest{}))
	router.POST("/file/meta", proxyutil.WrapBizFunc(file.GetMetaInfo, &fileinfo.GetFileMetaRequest{}))
	for _, bk := range s.c.s3Buckets {
		bucketPath := fmt.Sprintf("/%s", bk)
		routerPath := fmt.Sprintf("%s/*s3Param", bucketPath)
		router.GET(bucketPath, middleware.S3InfoExtractMiddleware(), s3.GetBucket)
		router.GET(routerPath, middleware.S3InfoExtractMiddleware(), s3.Download)
		router.PUT(routerPath, authMiddleware, middleware.S3InfoExtractMiddleware(), s3.Upload)
	}
}

func (s *Server) initS3(router *gin.Engine) {

}

func (s *Server) Run() error {
	return s.engine.Run(s.addr)
}
