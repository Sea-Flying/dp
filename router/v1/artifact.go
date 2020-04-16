package router

import (
	"github.com/gin-gonic/gin"
	"voyageone.com/dp/artifactKeeper/handler"
)

func initArtifactApiGroup(ApiRouter *gin.RouterGroup) {
	artifactApiGroup := ApiRouter.Group("artifact")
	ClassRouter := artifactApiGroup.Group("class")
	{
		ClassRouter.POST("", handler.CreateClass)
		ClassRouter.GET("")
		ClassRouter.DELETE("")
	}
	EntityRouter := artifactApiGroup.Group("entity")
	{
		EntityRouter.POST("")
		EntityRouter.GET("")
	}
}
