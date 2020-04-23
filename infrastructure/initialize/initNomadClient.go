package initialize

import (
	"github.com/hashicorp/nomad/api"
	"log"
	"os"
	"voyageone.com/dp/infrastructure/entity/config"
	"voyageone.com/dp/infrastructure/entity/global"
)

func InitNomadClient() {
	var err error
	global.NomadClient, err = initNomadClient(global.DPConfig.Nomad)
	if err != nil {
		log.Fatal(err)
	}
}

func initNomadClient(nomadConfig config.NomadConfig) (*api.Client, error) {
	config := api.DefaultConfig()
	config.Address = nomadConfig.NomadApiUrl
	if nomadConfig.NomadRegion != "" {
		config.Region = nomadConfig.NomadRegion
	}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(nomadConfig.NomadJobTplDir, os.ModeDir|os.ModePerm)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(nomadConfig.NomadJobHclDir, os.ModeDir|os.ModePerm)
	if err != nil {
		return nil, err
	}
	return client, nil
}
