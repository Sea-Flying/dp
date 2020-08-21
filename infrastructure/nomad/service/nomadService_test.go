package service

import (
	"fmt"
	"github.com/hashicorp/nomad/api"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"testing"
	"voyageone.com/dp/infrastructure/model/config"
)

var nomadClient *VoNomadClient

func initNomadClient(nomadConfig config.NomadConfig) (*VoNomadClient, error) {
	clientConfig := api.DefaultConfig()
	clientConfig.Address = nomadConfig.NomadApiUrl
	if nomadConfig.NomadRegion != "" {
		clientConfig.Region = nomadConfig.NomadRegion
	}
	client, err := api.NewClient(clientConfig)
	voNomadClient := NewVoNomadClient(client, log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile))
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
	return voNomadClient, nil
}

func init() {
	var dpConfig config.DPConfig
	_ = cleanenv.ReadConfig("E:/Develop/go/dp/dp.yml", &dpConfig)
	var err error
	nomadClient, err = initNomadClient(dpConfig.Nomad)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetJobLastDeploymentHealth(t *testing.T) {
	health, err := nomadClient.GetJobLastDeploymentHealth("count-api")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(health)
	}
}

func TestGetJobsList(t *testing.T) {
	jobs, err := nomadClient.GetJobsList()
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(jobs)
	}
}

func TestStopJob(t *testing.T) {
	err := nomadClient.StopJob("count-api")
	if err != nil {
		t.Error(err)
	}
}

func TestGetLastDeploymentAllocations(t *testing.T) {
	a, err := nomadClient.GetLastDeploymentAllocations("count-api")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(a)
	}
}

func TestGetJobJson(t *testing.T) {
	j, err := nomadClient.GetJobJson("count-api")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(j)
	}
}

func TestRestartJob(t *testing.T) {
	err := nomadClient.RestartJob("count-api")
	if err != nil {
		t.Error(err)
	}
}
