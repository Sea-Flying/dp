package initialize

import (
	"gopkg.in/olahol/melody.v1"
	. "voyageone.com/dp/infrastructure/model/global"
)

func InitWebsocket() {
	WsAppsStatus = melody.New()
	WsAppsStatus.HandleConnect(func(s *melody.Session) {

	})
}
