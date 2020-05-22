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
		artifactApiGroup.GET("class", artifactHandler.GetClass)
		artifactApiGroup.DELETE("class")

		artifactApiGroup.POST("entity", artifactHandler.CreateEntity)
		artifactApiGroup.GET("entity")
	}
}
