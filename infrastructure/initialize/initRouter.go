package initialize

import (
	"github.com/gin-gonic/gin"
	"voyageone.com/dp/router/v1"
)

func InitRouter() (Router *gin.Engine) {
	Router = gin.Default()
	router.InitApiRouter(Router)
	return
}
