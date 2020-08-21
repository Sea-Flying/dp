package service

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"testing"
	"voyageone.com/dp/infrastructure/model/config"
)

var consulClient *VoConsulClient

func initConsulClient(consulConfig config.ConsulConfig) (*VoConsulClient, error) {
	clientConfig := api.DefaultConfig()
	clientConfig.Address = consulConfig.ConsulApiUrl
	if consulConfig.ConsulDataCenter != "" {
		clientConfig.Datacenter = consulConfig.ConsulDataCenter
	}
	client, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}
	voConsulClient := NewVoConsulClient(client, log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile))
	return voConsulClient, nil
}

func init() {
	var dpConfig config.DPConfig
	_ = cleanenv.ReadConfig("E:/Develop/go/dp/dp.yml", &dpConfig)
	var err error
	consulClient, err = initConsulClient(dpConfig.Consul)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetServiceHealthCheckPassNum(t *testing.T) {
	count, err := consulClient.GetServiceHealthCheckPassNum("openvms-restapi")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(count)
	}
}
