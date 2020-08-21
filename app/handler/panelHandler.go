package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	watcherService "voyageone.com/dp/app/service/watcher"
	. "voyageone.com/dp/infrastructure/model/global"
)

type appCtlReq struct {
	AppName string `json:"app_name"`
	Action  string `json:"action"`
}

func WsGetAppsStatus(c *gin.Context) {
	WsAppsStatus.HandleRequest(c.Writer, c.Request)
}

func GetAppStatus(c *gin.Context) {
	appId := c.Param("appId")
	appStatus := watcherService.GetAppStatus(appId)
	if appStatus != "" {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusBadRequest)
	}
	c.Writer.Write([]byte(appStatus))
}

func GetAppsStatus(c *gin.Context) {
	ret := watcherService.GetAppsStatus()
	c.JSON(http.StatusOK, ret)
}

func ControlApp(c *gin.Context) {
	var appCtl appCtlReq
	c.ShouldBindJSON(&appCtl)
	var err error
	switch appCtl.Action {
	case "start":
		err = watcherService.StartApp(appCtl.AppName)
		break
	case "stop":
		err = watcherService.StopApp(appCtl.AppName)
		break
	case "restart":
		err = watcherService.RestartApp(appCtl.AppName)
		break
	}
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Writer.Write([]byte("error: " + err.Error()))
	} else {
		c.Status(http.StatusOK)
		c.Writer.Write([]byte(fmt.Sprintf(`actiont "%s" for app "%s" executed`, appCtl.Action, appCtl.AppName)))
	}
}

func AppMarkAsHealthy() {

}

func AppMarkAsUnhealthy() {

}
