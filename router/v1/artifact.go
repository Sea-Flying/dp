package v1

import "github.com/gin-gonic/gin"

func initArtifactApiGroup(ApiRouter *gin.RouterGroup) {
	ClassRouter := ApiRouter.Group("class")
	{
		ClassRouter.POST("")
		ClassRouter.GET("")
		ClassRouter.PUT("")
		ClassRouter.DELETE("")
	}
	EntityRouter := ApiRouter.Group("entity")
	{
		EntityRouter.POST("")
		EntityRouter.GET("")
	}
}
