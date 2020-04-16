package main

import (
	"log"
	"net/http"
	"time"
	"voyageone.com/dp/infrastructure/entity/global"
	"voyageone.com/dp/infrastructure/initialize"
	"voyageone.com/dp/infrastructure/utils"
)

func main() {
	log.Println(utils.GetExecPath())
	initialize.InitConfig(&global.DPConfig)
	initialize.InitCqlSession(global.DPConfig)
	dpRouter := initialize.InitRouter()
	dpServer := &http.Server{
		Addr:           global.DPConfig.Base.HttpHost + ":" + global.DPConfig.Base.HttpPort,
		Handler:        dpRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	dpServer.ListenAndServe()
}
