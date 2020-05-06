package main

import (
	"flag"
	"net/http"
	"time"
	"voyageone.com/dp/infrastructure/initialize"
	. "voyageone.com/dp/infrastructure/model/global"
)

var dpConfigPath = flag.String("config", "D:/Develop/go/dp/dp.yml", "the YAML config file ")

func main() {
	flag.Parse()
	initialize.InitConfig(*dpConfigPath)
	initialize.InitCqlSession()
	initialize.InitNomadClient()
	initialize.InitLogger()
	dpRouter := initialize.InitRouter()
	dpServer := &http.Server{
		Addr:           DPConfig.Base.HttpHost + ":" + DPConfig.Base.HttpPort,
		Handler:        dpRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := dpServer.ListenAndServe()
	if err != nil {
		DPLogger.Fatal(err)
	}
}
