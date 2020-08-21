package service

import (
	"github.com/hashicorp/consul/api"
	"log"
)

type VoConsulClient struct {
	*api.Client
	Logger *log.Logger
}

func NewVoConsulClient(client *api.Client, logger *log.Logger) *VoConsulClient {
	return &VoConsulClient{
		Client: client,
		Logger: logger,
	}
}

func (client *VoConsulClient) GetServiceHealthCheckPassNum(appId string) (int, error) {
	healthEndpoint := client.Health()
	instances, _, err := healthEndpoint.Service(appId, "green", true, nil)
	return len(instances), err
}
