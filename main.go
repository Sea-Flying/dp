package main

import (
	"net/http"
	"time"
	. "voyageone.com/dp/infrastructure/entity/global"
	"voyageone.com/dp/infrastructure/initialize"
)

func main() {
	initialize.InitConfig()
	initialize.InitCqlSession()
	initialize.InitNomadClient()
	dpRouter := initialize.InitRouter()
	dpServer := &http.Server{
		Addr:           DPConfig.Base.HttpHost + ":" + DPConfig.Base.HttpPort,
		Handler:        dpRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	dpServer.ListenAndServe()
}
