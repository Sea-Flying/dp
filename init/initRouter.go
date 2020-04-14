package init

import (
	"github.com/gin-gonic/gin"
	v1 "voyageone.com/dp/router/v1"
)

func InitRouter() (Router *gin.Engine) {
	Router = gin.Default()
	v1.InitApiRouter(Router)
	return
}
