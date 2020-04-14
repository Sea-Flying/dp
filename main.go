package main

import (
	"net/http"
	"time"
	"voyageone.com/dp/global"
	"voyageone.com/dp/init"
)

func main() {
	init.InitConfig(&global.DPConfig)
	init.InitCqlSession(&global.CqlSession, global.DPConfig)
	dpRouter := init.InitRouter()
	dpServer := &http.Server{
		Addr:           global.DPConfig.Base.HttpHost + ":" + global.DPConfig.Base.HttpPort,
		Handler:        dpRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	dpServer.ListenAndServe()
}
