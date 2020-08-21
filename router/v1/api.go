package router

import (
	"github.com/gin-gonic/gin"
)

func InitApiRouter(Engine *gin.Engine) {
	ApiRouter := Engine.Group("api/v1")
	{
		initArtifactApiGroup(ApiRouter)
		InitScheduleApiGroup(ApiRouter)
		initPanelApiGroup(ApiRouter)
	}
}
