package router

import (
	"github.com/gin-gonic/gin"
	artifactHandler "voyageone.com/dp/artifactKeeper/handler"
)

func initArtifactApiGroup(apiRootRouter *gin.RouterGroup) {
	artifactApiGroup := apiRootRouter.Group("artifact")
	{
		artifactApiGroup.POST("repo", artifactHandler.CreataRepo)

		artifactApiGroup.POST("class", artifactHandler.CreateClass)
		artifactApiGroup.GET("class")
		artifactApiGroup.DELETE("class")

		artifactApiGroup.POST("model", artifactHandler.CreateEntity)
		artifactApiGroup.GET("model")
	}
}
