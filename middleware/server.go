//Package middleware 中间件仓库
package middleware

import "github.com/gin-gonic/gin"
import "code.htres.cn/casicloud/adc-genius/middleware/v1"

// Server 中间件API服务,使用gin实现
type Server struct {
	Route *gin.Engine
}

// Config 服务器相关配置
type Config struct {
}

//NewServer 构造Server
func NewServer(config *Config) (*Server, error) {
	server := &Server{}
	// build route
	router := gin.New()
	router.Use(gin.Recovery())
	// todo: 配置gin的日志格式
	router.Use(gin.Logger())
	server.Route = router

	return server, server.setUp()
}

//Start 启动Web服务
func (s *Server) Start(addr ...string) error {
	return s.Route.Run(addr...)
}

func (s *Server) setUp() error {
	// api healh related
	v1.ApplyRoutes(s.Route)
	s.Route.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return nil
}
