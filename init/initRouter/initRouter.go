package initRouter

import "github.com/gin-gonic/gin"

func InitRouter() *gin.Engine {
	var Router = gin.Default()
	ApiGroup := Router.Group("")
}
