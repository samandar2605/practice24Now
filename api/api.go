package api

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/practice2311/api/v1"
	"github.com/practice2311/config"
	"github.com/practice2311/storage"

	_ "github.com/practice2311/api/docs"       // for swagger
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type RouterOptions struct {
	Cfg      *config.Config
	Storage  storage.StorageI
	InMemory storage.InMemoryStorageI
}

// @title           Swagger for blog api
// @version         1.0
// @description     This is a blog service api.
// @host      		localhost:8000
// @BasePath  		/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @Security ApiKeyAuth
func New(opt RouterOptions) *gin.Engine {
	router := gin.Default()

	handlerV1 := v1.New(&v1.HandlerV1Options{
		Cfg:      opt.Cfg,
		Storage:  opt.Storage,
		InMemory: opt.InMemory,
	})
	router.Static("/media", "./media")
	apiV1 := router.Group("/v1")

	// Register
	apiV1.POST("/users/auth", handlerV1.Register)
	apiV1.POST("/auth/verify", handlerV1.Verify)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
