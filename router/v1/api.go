package v1

import (
	"github.com/gin-gonic/gin"
)

func InitApiRouter(Engine *gin.Engine) {
	ApiRouter := Engine.Group("api/v1")
	{
		initArtifactApiGroup(ApiRouter)
		initManagementApiGroup(ApiRouter)
	}
}
