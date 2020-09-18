package router

import (
	"github.com/gin-gonic/gin"
	"voyageone.com/dp/app/handler"
)

func initPanelApiGroup(apiRootRouter *gin.RouterGroup) {
	panelApiGroup := apiRootRouter.Group("panel")
	{
		panelApiGroup.GET("ws/appstatus", handler.WsGetAppsStatus)
		panelApiGroup.GET("apps-status", handler.GetAppsStatus)
		panelApiGroup.GET("app-status/:appId", handler.GetAppStatus)
		panelApiGroup.GET("apps-histories", handler.GetAppsStatusHistories)
		panelApiGroup.GET("app-histories/:appId", handler.GetAppStatusHistories)
		panelApiGroup.POST("control", handler.ControlApp)
	}
}
