package initialize

import (
	"fmt"
	"gopkg.in/olahol/melody.v1"
	"time"
	watcherService "voyageone.com/dp/app/service/watcher"
	. "voyageone.com/dp/infrastructure/model/global"
)

func InitWebsocket() {
	WsAppsStatus = melody.New()
	WsAppsStatus.HandleConnect(func(s *melody.Session) {
		s.Write([]byte(generateWsRet()))
	})

	WsAppsStatus.HandleMessage(func(s *melody.Session, msg []byte) {
		WsAppsStatus.Broadcast([]byte(generateWsRet()))
	})

	wsTimer := time.NewTicker(5 * time.Second)
	go func() {
		for {
			<-wsTimer.C
			WsAppsStatus.Broadcast([]byte(generateWsRet()))
		}
	}()

}

func generateWsRet() string {
	var ret = ``
	ret += `[`
	for appId, appEntity := range watcherService.AppWatchList {
		ret += fmt.Sprintf(`{"appId":"%s","status":"%s"}`, appId, appEntity.Current())
		ret += `,`
	}
	ret = ret[0 : len(ret)-1]
	ret += `]`
	return ret
}
