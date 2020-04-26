package server

import (
	"github.com/ToolPackage/fse/server/api"
	"github.com/ToolPackage/fse/server/config"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Server struct {
	engine *gin.Engine
}

func New() *Server {
	server := &Server{
		engine: gin.Default(),
	}

	switch config.Conf.RunMode {
	case config.ModeDev:
		gin.SetMode(gin.DebugMode)
	case config.ModeProd:
		gin.SetMode(gin.ReleaseMode)
	default:
		panic(errors.New("未知的mode:" + config.Conf.RunMode))
	}

	return server
}

func (s *Server) Start() {
	s.initRouter()

	addr := config.Conf.Host + ":" + config.Conf.Port
	err := s.engine.Run(addr)
	if err != nil {
		panic(errors.WithStack(err))
	}
}

func (s *Server) initRouter() {
	s.engine.GET("/files", api.GetFilesList)
	s.engine.POST("/files", api.PostFile)
	s.engine.GET("/files/:fileId", api.GetFile)
	s.engine.DELETE("/files/:fileId", api.DeleteFile)
}
