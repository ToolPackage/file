package server

import (
	"fmt"
	"github.com/ToolPackage/fse/server/api"
	"github.com/ToolPackage/fse/server/config"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strings"
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
	s.engine.Use(Cors())
	s.engine.GET("/files", api.GetFilesList)
	s.engine.POST("/files", api.PostFile)
	s.engine.GET("/files/:fileId", api.GetFile)
	s.engine.DELETE("/files/:fileId", api.DeleteFile)
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			//下面的都是乱添加的-_-~
			// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "0.0.0.0")
			c.Header("Access-Control-Allow-Headers", headerStr)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			// c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			// c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}
