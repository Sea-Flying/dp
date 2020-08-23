package initialize

import (
	"github.com/hashicorp/consul/api"
	"log"
	"net"
	"net/http"
	"time"
	consulService "voyageone.com/dp/infrastructure/consul/service"
	"voyageone.com/dp/infrastructure/model/config"
	"voyageone.com/dp/infrastructure/model/global"
)

func InitConsulClient() {
	var err error
	global.ConsulClient, err = initConsulClient(global.DPConfig.Consul)
	if err != nil {
		log.Fatal(err)
	}
}

func initConsulClient(consulConfig config.ConsulConfig) (*consulService.VoConsulClient, error) {
	clientConfig := api.DefaultConfig()
	clientConfig.Address = consulConfig.ConsulApiUrl
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
	if consulConfig.ConsulDataCenter != "" {
		clientConfig.Datacenter = consulConfig.ConsulDataCenter
	}
	client, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}
	voConsulClient := consulService.NewVoConsulClient(client, global.DPLogger)
	return voConsulClient, nil
}
