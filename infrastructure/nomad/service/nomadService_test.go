package service

import (
	"fmt"
	"github.com/hashicorp/nomad/api"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"testing"
	"voyageone.com/dp/infrastructure/model/config"
	"voyageone.com/dp/infrastructure/model/global"
)

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

func init() {
	_ = cleanenv.ReadConfig("D:/Develop/go/dp/dp.yml", &global.DPConfig)
	var err error
	global.NomadClient, err = initNomadClient(global.DPConfig.Nomad)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetJobLastDeploymentHealth(t *testing.T) {
	health, err := GetJobLastDeploymentHealth(global.NomadClient, "openvms-restapi-dp")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(health)
	}
}
