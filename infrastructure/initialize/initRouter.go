package initialize

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"voyageone.com/dp/router/v1"
)

func InitRouter() (Router *gin.Engine) {
	Router = gin.Default()
	Router.Use(cors.Default())
	router.InitApiRouter(Router)
	return
}
