package router

import (
	"github.com/gin-gonic/gin"
	"voyageone.com/dp/app/handler"
)

func initPanelApiGroup(apiRootRouter *gin.RouterGroup) {
	panelApiGroup := apiRootRouter.Group("panel")
	{
		panelApiGroup.GET("ws/appstatus")
		panelApiGroup.GET("apps-status", handler.GetAppsStatus)
		panelApiGroup.GET("app-status/:appId", handler.GetAppStatus)
		panelApiGroup.POST("control", handler.ControlApp)
	}
}
