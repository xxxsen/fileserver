package server

import (
	"fileserver/proxyutil"
	"fileserver/server/handler/file"
	"fileserver/server/handler/s3"
	"fileserver/server/middleware"
	"fileserver/server/model"
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
	s.initMiddleware(s.engine)
	s.initAPI(s.engine)
	return nil
}

func (s *Server) initMiddleware(router *gin.Engine) {
	mds := []gin.HandlerFunc{
		middleware.PanicRecoverMiddleware(),
		middleware.TraceMiddleware(),
		middleware.LogRequestMiddleware(),
	}
	router.Use(mds...)
}

func (s *Server) initAPI(router *gin.Engine) {
	authMiddleware := middleware.CommonAuth(s.c.userMap)
	fileRouter := router.Group("/file")
	fileRouter.POST("/upload", authMiddleware, proxyutil.WrapBizFunc(file.FileUpload, &model.UploadFileRequest{}))
	fileRouter.GET("/download/:key", file.FileDownload)
	fileRouter.GET("/meta/:key", file.GetMetaInfo)
	for _, bk := range s.c.s3Buckets {
		bucketPath := fmt.Sprintf("/:%s", bk)
		routerPath := fmt.Sprintf("%s/*object", bucketPath)
		router.GET(bucketPath, middleware.ExtractS3InfoMiddleware(), s3.GetBucket)
		router.GET(routerPath, middleware.ExtractS3InfoMiddleware(), s3.DownloadObject)
		router.PUT(routerPath, authMiddleware, middleware.ExtractS3InfoMiddleware(), s3.UploadObject)
	}
}
func (s *Server) Run() error {
	return s.engine.Run(s.addr)
}
