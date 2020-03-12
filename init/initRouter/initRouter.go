package initRouter

import "github.com/gin-gonic/gin"

func InitRouter() (ApiGroup *gin.RouterGroup) {
	var Router = gin.Default()
	ApiGroup = Router.Group("")
	return
}
