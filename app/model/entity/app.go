package entity

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/nomad/api"
	"github.com/looplab/fsm"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"
	panelRepository "voyageone.com/dp/app/model/repository"
	artifactRepository "voyageone.com/dp/artifactKeeper/model/repository"
	"voyageone.com/dp/infrastructure/model/customType"
	. "voyageone.com/dp/infrastructure/model/global"
)

type App struct {
	AppId                string
	WatcherTicker        *time.Ticker
	Deleted              chan struct{}
	TimeoutCounter       *time.Timer
	TimeoutCounterStatus string
	RefreshMutex         sync.Mutex
	UnitTimeoutSeconds   int
	TimeoutFactor        int
	*fsm.FSM
}

func NewApp(appId string) *App {
	unitTimeout, timeoutFactor := initAppTimeoutFiled(appId)
	app := App{
		AppId:                appId,
		Deleted:              make(chan struct{}),
		TimeoutCounter:       time.NewTimer(math.MaxInt64),
		TimeoutCounterStatus: "idle",
		UnitTimeoutSeconds:   unitTimeout,
		TimeoutFactor:        timeoutFactor,
		FSM: fsm.NewFSM(
			"unhealthy",
			fsm.Events{
				{Name: "deploy", Src: []string{"stopped", "healthy", "unhealthy"}, Dst: "deploying"},
				{Name: "start", Src: []string{"stopped"}, Dst: "starting"},
				{Name: "stop", Src: []string{"deploying", "starting", "restarting", "healthy", "unhealthy"}, Dst: "stopping"},
				{Name: "restart", Src: []string{"healthy", "unhealthy"}, Dst: "restarting"},
				{Name: "bootTimely", Src: []string{"deploying", "starting", "restarting"}, Dst: "healthy"},
				{Name: "bootTimeout", Src: []string{"deploying", "starting", "restarting"}, Dst: "unhealthy"},
				//确认已停止，可以被手动触发，或自动状态轮询触发，正常情况下，手动停止的状态流转应该为 x -> stopping -> stopped
				{Name: "confirmStopped", Src: []string{"healthy", "unhealthy", "deploying", "starting", "restarting", "stopping"}, Dst: "stopped"},
				{Name: "markAsUnhealthy", Src: []string{"stopped", "healthy"}, Dst: "unhealthy"},
				{Name: "markAsHealthy", Src: []string{"stopped", "unhealthy"}, Dst: "healthy"},
			},
			fsm.Callbacks{},
		),
	}
	app.TimeoutCounter.Stop()
	return &app
}

func (a *App) StartWatcher() {
	//5s轮询一次，加上一个2s内的随机时间分散请求
	a.WatcherTicker = time.NewTicker(time.Duration(5000+rand.Int63n(2000)) * time.Millisecond)
	go func() {
		for {
			select {
			case <-a.WatcherTicker.C:
				a.RefreshAppsStatus()
			case <-a.Deleted:
				close(a.Deleted)
				return
			}
		}
	}()
}

func (a *App) StopWatcher() {
	a.Deleted <- struct{}{}
	a.WatcherTicker.Stop()
}

func (a *App) RefreshAppsStatus() {
	var err error
	jobStatus, err := getJobStatusFromNomadAndConsul(a.AppId)
	if err != nil {
		if strings.Contains(err.Error(), "job not found") {
			a.SetState("stopped")
		} else {
			DPLogger.Printf("error when get job [%s] status from Nomad\n, %v", a.AppId, err)
		}
		return
	}
	a.RefreshMutex.Lock()
	defer a.RefreshMutex.Unlock()
	switch a.Current() {
	case "stopped":
		if jobStatus == "scheduling" {
			err = a.MarkAsUnhealthy()
		} else if jobStatus == "healthy" {
			err = a.MarkAsHealthy()
		}
		break
	case "stopping", "healthy":
		if jobStatus == "stopped" {
			err = a.ConfirmStopped()
		}
		break
	case "starting", "deploying", "restarting":
		if jobStatus == "healthy" {
			err = a.BootTimely()
		} else if jobStatus == "scheduling" && a.TimeoutCounterStatus == "done" {
			err = a.BootTimeout()
		} else if jobStatus == "stopped" {
			err = a.ConfirmStopped()
		}
		break
	case "unhealthy":
		if jobStatus == "stopped" {
			err = a.ConfirmStopped()
		} else if jobStatus == "healthy" {
			err = a.MarkAsHealthy()
		}
		break
	}
	if err != nil {
		DPLogger.Printf(`error when RefreshAppsStatus("%s"), current status: %s`, a.AppId, a.Current())
	}
}

func (a *App) Stop() (err error) {
	if a.Can("stop") {
		//获取当前job定义Json, 落库保存
		jobJson, _ := NomadClient.GetJobJson(a.AppId)
		c := panelRepository.StoppedAppCache{
			AppName:      a.AppId,
			NomadJobJson: jobJson,
		}
		err = c.CreateOrUpdate()
		if err != nil {
			return
		}
		//Nomad stop job
		err = NomadClient.StopJob(a.AppId)
		if err != nil {
			return
		}
		//app实体状态流转
		err = a.Event("stop")
		if err != nil {
			return
		}
		a.TimeoutCounter.Stop()
		a.TimeoutCounterStatus = "idle"
		//状态变化落库
		err = panelRepository.AddStatusHistory(a.AppId, a.Current())
		return
	} else {
		return customType.DPError(fmt.Sprintf(`can not execute action "stop", app "%s" current status is [%s]`, a.AppId, a.Current()))
	}
}

func (a *App) Start() (err error) {
	if a.Can("start") {
		//从库中获取上次停止时的job定义
		c := panelRepository.StoppedAppCache{
			AppName: a.AppId,
		}
		err = c.GetByPrimaryKey()
		if err != nil {
			return
		}
		var job api.Job
		json.Unmarshal([]byte(c.NomadJobJson), &job)
		//Nomad run job
		err = NomadClient.SubmitJob(&job)
		if err != nil {
			return
		}
		//app实体状态流转
		err = a.Event("start")
		if err != nil {
			return
		}
		a.TimeoutCounter.Reset(time.Duration(a.UnitTimeoutSeconds) * time.Duration(a.TimeoutFactor) * time.Second)
		a.TimeoutCounterStatus = "counting"
		go func() {
			<-a.TimeoutCounter.C
			a.TimeoutCounterStatus = "done"
		}()
		//状态变化入库
		err = panelRepository.AddStatusHistory(a.AppId, a.Current())
		return
	} else {
		return customType.DPError(fmt.Sprintf(`can not execute action "start", app "%s" current status is [%s]`, a.AppId, a.Current()))
	}
}

func (a *App) Restart() (err error) {
	if a.Can("restart") {
		err = NomadClient.RestartJob(a.AppId)
		if err != nil {
			return err
		}
		err = a.Event("restart")
		if err != nil {
			return err
		}
		//重启app，超时计时器比发布时长30s，用于前导的停止app耗时
		a.TimeoutCounter.Reset(time.Duration(a.UnitTimeoutSeconds+30) * time.Duration(a.TimeoutFactor) * time.Second)
		a.TimeoutCounterStatus = "counting"
		go func() {
			<-a.TimeoutCounter.C
			a.TimeoutCounterStatus = "done"
		}()
		err = panelRepository.AddStatusHistory(a.AppId, a.Current())
		return
	} else {
		return customType.DPError(fmt.Sprintf(`can not execute action "restart", app "%s" current status is [%s]`, a.AppId, a.Current()))
	}
}

func (a *App) Deploy() (err error) {
	if a.Can("deploy") {
		err = a.Event("deploy")
		if err != nil {
			return err
		}
		a.TimeoutCounter.Reset(time.Duration(a.UnitTimeoutSeconds) * time.Duration(a.TimeoutFactor) * time.Second)
		a.TimeoutCounterStatus = "counting"
		go func() {
			<-a.TimeoutCounter.C
			a.TimeoutCounterStatus = "done"
		}()
		err = panelRepository.AddStatusHistory(a.AppId, a.Current())
		return
	} else {
		return customType.DPError(fmt.Sprintf(`can not execute action "deploy", app "%s" current status is [%s]`, a.AppId, a.Current()))
	}
}

func (a *App) BootTimely() (err error) {
	if a.Can("bootTimely") {
		err = a.Event("bootTimely")
		if err != nil {
			return err
		}
		err = panelRepository.AddStatusHistory(a.AppId, a.Current())
		return
	} else {
		return customType.DPError(fmt.Sprintf(`can not execute action "deploy", app "%s" current status is [%s]`, a.AppId, a.Current()))
	}
}

func (a *App) BootTimeout() (err error) {
	if a.Can("bootTimeout") {
		err = a.Event("bootTimeout")
		if err != nil {
			return err
		}
		err = panelRepository.AddStatusHistory(a.AppId, a.Current())
		return
	} else {
		return customType.DPError(fmt.Sprintf(`can not execute action "bootTimeout", app "%s" current status is [%s]`, a.AppId, a.Current()))
	}
}

func (a *App) ConfirmStopped() (err error) {
	if a.Can("confirmStopped") {
		err = a.Event("confirmStopped")
		if err != nil {
			return err
		}
		err = panelRepository.AddStatusHistory(a.AppId, a.Current())
		return
	} else {
		return customType.DPError(fmt.Sprintf(`can not execute action "confirmStopped", app "%s" current status is [%s]`, a.AppId, a.Current()))
	}
}

func (a *App) MarkAsHealthy() (err error) {
	if a.Can("markAsHealthy") {
		err = a.Event("markAsHealthy")
		if err != nil {
			return err
		}
		err = panelRepository.AddStatusHistory(a.AppId, a.Current())
		return
	} else {
		return customType.DPError(fmt.Sprintf(`can not execute action "markAsHealthy", app "%s" current status is [%s]`, a.AppId, a.Current()))
	}
}

func (a *App) MarkAsUnhealthy() (err error) {
	if a.Can("markAsUnhealthy") {
		err = a.Event("markAsUnhealthy")
		if err != nil {
			return err
		}
		err = panelRepository.AddStatusHistory(a.AppId, a.Current())
		return
	} else {
		return customType.DPError(fmt.Sprintf(`can not execute action "markAsUnhealthy", app "%s" current status is [%s]`, a.AppId, a.Current()))
	}
}

func getJobStatusFromNomadAndConsul(jobId string) (status string, err error) {
	nomadJob, err := NomadClient.GetJob(jobId)
	if err != nil {
		return "unknown", err
	}
	consulServiceHealthyInstanceCount, err := ConsulClient.GetServiceHealthCheckPassNum(jobId)
	if err != nil {
		return "unknown", err
	}
	lastDeploymentStatus := NomadClient.GetJobLastDeploymentStatus(jobId)

	if *(nomadJob.Status) == "dead" || *(nomadJob.Stop) {
		return "stopped", nil
	} else if *(nomadJob.Status) == "running" {
		if *(nomadJob.TaskGroups[0].Count) == consulServiceHealthyInstanceCount {
			//最近的一个Deployment正在运行，且服务当前的运行状态时健康的，则说明该服务正在通过dp(灰度)发布/重启中
			if lastDeploymentStatus == "running" || lastDeploymentStatus == "paused" {
				return "scheduling", nil
			} else {
				return "healthy", nil
			}
		} else {
			return "scheduling", nil
		}
	} else {
		return "scheduling", nil
	}
}

func initAppTimeoutFiled(appId string) (unitTimeout int, factor int) {
	var err error
	artifactClass := artifactRepository.Class{
		Name: appId,
	}
	err = artifactClass.GetByPrimaryKey()
	if err != nil {
		DPLogger.Printf("get artifact class error when initAppTimeoutFiled, error: %v\n", err)
		//默认超时120s
		unitTimeout = 120
	} else {
		unitTimeout = artifactClass.UnitTimeoutSeconds
	}
	nomadJob, err := NomadClient.GetJob(appId)
	if err != nil {
		DPLogger.Printf("get nomad job error when initAppTimeoutFiled, error: %v\n", err)
		//默认实例并行后因子1
		factor = 1
	} else {
		var parallel int
		if *(nomadJob.Update.Canary) == 0 {
			parallel = *(nomadJob.Update.MaxParallel)
		} else {
			parallel = *(nomadJob.Update.Canary)
		}
		count := *(nomadJob.TaskGroups[0].Count)
		factor = int(math.Ceil(float64(count) / float64(parallel)))
	}
	return
}
