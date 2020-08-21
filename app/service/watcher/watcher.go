package watcher

import (
	"regexp"
	"time"
	appPanelEntity "voyageone.com/dp/app/model/entity"
	panelRepository "voyageone.com/dp/app/model/repository"
	. "voyageone.com/dp/infrastructure/model/global"
	"voyageone.com/dp/infrastructure/utils"
)

var (
	AppWatchRegexpInclude *regexp.Regexp
	AppWatchRegexpExclude *regexp.Regexp
	AppWatchList          map[string]*appPanelEntity.App
)

func InitAppWatcher() {
	AppWatchRegexpInclude = regexp.MustCompile(DPConfig.AppPanel.WatcherAppsRegexpInclude)
	AppWatchRegexpExclude = regexp.MustCompile(DPConfig.AppPanel.WatcherAppsRegexpExclude)
	AppWatchList = make(map[string]*appPanelEntity.App)
	RefreshAppWatchList()
	ticker1m := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			<-ticker1m.C
			RefreshAppWatchList()
		}
	}()
}

//周期性执行，根据nomad的job list 和 stopped_app_cache表的 name 的并集，然后经过正则表达式规则过滤，获取到需要监控的应用清单
func RefreshAppWatchList() {
	jobsList, err := NomadClient.GetJobsList()
	if err != nil {
		DPLogger.Println("error when get jobs list from Nomad")
		//TODO mark nomad server unresponsive
		return
	}
	stoppedApps, err := panelRepository.GetStoppedAppsList()
	if err != nil {
		DPLogger.Println("error when get stopped apps list from Cassandra")
	}
	watchList := utils.MergeStringSlice(jobsList, stoppedApps)
	//检查最新的watchlist是否有未加入轮询列表中的app
	for _, appName := range watchList {
		if _, existed := AppWatchList[appName]; !existed {
			if AppWatchRegexpInclude.MatchString(appName) && !AppWatchRegexpExclude.MatchString(appName) {
				AddApp(appName)
			}
		}
	}
	//检查是否有不需要监控的App, 移出轮询列表
	for appName := range AppWatchList {
		if !utils.IsInSlice(watchList, appName) {
			DelApp(appName)
		}
	}
}

func AddApp(appId string) {
	AppWatchList[appId] = appPanelEntity.NewApp(appId)
	AppWatchList[appId].StartWatcher()
}

func DelApp(appId string) {
	AppWatchList[appId].StopWatcher()
	delete(AppWatchList, appId)
}

func GetAppStatus(appId string) string {
	return AppWatchList[appId].Current()
}

func GetAppsStatus() (ret map[string]string) {
	ret = make(map[string]string)
	for appName, app := range AppWatchList {
		ret[appName] = app.Current()
	}
	return
}

func StopApp(appId string) error {
	return AppWatchList[appId].Stop()
}

func StartApp(appId string) error {
	return AppWatchList[appId].Start()
}

func RestartApp(appId string) error {
	return AppWatchList[appId].Restart()
}
