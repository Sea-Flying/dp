package initialize

import (
	"github.com/hashicorp/consul/api"
	"log"
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
