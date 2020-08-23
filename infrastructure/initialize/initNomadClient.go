package initialize

import (
	"github.com/hashicorp/nomad/api"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	"voyageone.com/dp/infrastructure/model/config"
	"voyageone.com/dp/infrastructure/model/global"
	nomadService "voyageone.com/dp/infrastructure/nomad/service"
)

func InitNomadClient() {
	var err error
	global.NomadClient, err = initNomadClient(global.DPConfig.Nomad)
	if err != nil {
		log.Fatal(err)
	}
}

func initNomadClient(nomadConfig config.NomadConfig) (*nomadService.VoNomadClient, error) {
	clientConfig := api.DefaultConfig()
	clientConfig.Address = nomadConfig.NomadApiUrl
	clientConfig.HttpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	if nomadConfig.NomadRegion != "" {
		clientConfig.Region = nomadConfig.NomadRegion
	}
	client, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}
	voNomadClient := nomadService.NewVoNomadClient(client, global.DPLogger)
	err = os.MkdirAll(nomadConfig.NomadJobTplDir, os.ModeDir|os.ModePerm)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(nomadConfig.NomadJobHclDir, os.ModeDir|os.ModePerm)
	if err != nil {
		return nil, err
	}
	return voNomadClient, nil
}
