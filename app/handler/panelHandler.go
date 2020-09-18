package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/scylladb/gocqlx/v2/qb"
	"net/http"
	panelRepository "voyageone.com/dp/app/model/repository"
	watcherService "voyageone.com/dp/app/service/watcher"
	. "voyageone.com/dp/infrastructure/model/global"
)

type appCtlReq struct {
	AppName string `json:"app_name"`
	Action  string `json:"action"`
}

type appHistoryReq struct {
	TimeOrder bool `json:"time_order"`
	PageSize  int  `json:"page_size"`
	PageNum   int  `json:"page_num"`
}

type appHistoryResp struct {
	Total        int                             `json:"total"`
	AppHistories []panelRepository.StatusHistory `json:"app_histories"`
	Error        string                          `json:"error"`
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

func GetAppStatusHistories(c *gin.Context) {
	appId := c.Param("appId")
	var historyReq appHistoryReq
	_ = c.ShouldBindJSON(&historyReq)
	var resp = appHistoryResp{}
	var err error
	resp.Total, resp.AppHistories, err =
		panelRepository.GetByAppNameOrderByTime(appId, qb.Order(historyReq.TimeOrder), historyReq.PageSize, historyReq.PageNum)
	var httpCode int
	if err != nil {
		resp.Error = err.Error()
		httpCode = http.StatusBadRequest
	} else {
		httpCode = http.StatusOK
	}
	c.JSON(httpCode, resp)
}

func GetAppsStatusHistories(c *gin.Context) {

}

func AppMarkAsHealthy() {

}

func AppMarkAsUnhealthy() {

}
