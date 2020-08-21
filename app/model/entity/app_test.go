package entity

import (
	"github.com/gocql/gocql"
	consulApi "github.com/hashicorp/consul/api"
	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/scylladb/gocqlx/v2"
	"log"
	"os"
	"testing"
	"time"
	consulService "voyageone.com/dp/infrastructure/consul/service"
	"voyageone.com/dp/infrastructure/model/config"
	"voyageone.com/dp/infrastructure/model/global"
	nomadService "voyageone.com/dp/infrastructure/nomad/service"
)

func initNomadClient(nomadConfig config.NomadConfig) (*nomadService.VoNomadClient, error) {
	clientConfig := nomadApi.DefaultConfig()
	clientConfig.Address = nomadConfig.NomadApiUrl
	if nomadConfig.NomadRegion != "" {
		clientConfig.Region = nomadConfig.NomadRegion
	}
	client, err := nomadApi.NewClient(clientConfig)
	voNomadClient := nomadService.NewVoNomadClient(client, log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile))
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

func initConsulClient(consulConfig config.ConsulConfig) (*consulService.VoConsulClient, error) {
	clientConfig := consulApi.DefaultConfig()
	clientConfig.Address = consulConfig.ConsulApiUrl
	if consulConfig.ConsulDataCenter != "" {
		clientConfig.Datacenter = consulConfig.ConsulDataCenter
	}
	client, err := consulApi.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}
	voConsulClient := consulService.NewVoConsulClient(client, global.DPLogger)
	return voConsulClient, nil
}

func init() {
	var dpConfig config.DPConfig
	_ = cleanenv.ReadConfig("E:/Develop/go/dp/dp.yml", &dpConfig)
	var err error
	global.NomadClient, err = initNomadClient(dpConfig.Nomad)
	if err != nil {
		log.Fatal(err)
	}
	global.ConsulClient, err = initConsulClient(dpConfig.Consul)
	if err != nil {
		log.Fatal(err)
	}
	cluster := gocql.NewCluster("127.0.0.1:9042")
	cluster.Consistency = gocql.Quorum
	cluster.NumConns = 3
	global.CqlSession, err = gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Fatal("CQL Session Create Failed!", err)
	}
	global.DPLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile)
}

var a = NewApp("openvms-restapi-dp")

func TestApp_Stop(t *testing.T) {
	a.SetState("healthy")
	a.StartWatcher()
	a.Stop()
	ticker1s := time.NewTicker(1 * time.Second)
	counter30s := time.NewTimer(30 * time.Second)
	for {
		select {
		case <-ticker1s.C:
			if a.Current() == "stopped" {
				return
			}
		case <-counter30s.C:
			t.Error("Stop timeout")
		}
	}
}

func TestApp_Start(t *testing.T) {
	a.UnitTimeoutSeconds = 180
	a.TimeoutFactor = 1
	a.SetState("healthy")
	a.StartWatcher()
	a.Start()

	ticker1s := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker1s.C:
			if a.Current() == "healthy" {
				return
			}
		case <-a.TimeoutCounter.C:
			a.TimeoutCounterStatus = "done"
			t.Error("Stop timeout")
			return
		}
	}
}

func TestApp_Restart(t *testing.T) {
	a.UnitTimeoutSeconds = 180
	a.TimeoutFactor = 1
	a.SetState("healthy")
	a.StartWatcher()
	a.Restart()

	ticker1s := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker1s.C:
			if a.Current() == "healthy" {
				return
			}
		case <-a.TimeoutCounter.C:
			a.TimeoutCounterStatus = "done"
			t.Error("Stop timeout")
			return
		}
	}
}

func TestTimerBlocking(t *testing.T) {
	timer := time.NewTimer(0)
	<-timer.C
	//timer.Stop()
	//timer.Reset(1*time.Second)
	go func() {
		<-timer.C
		println("Ding~")
	}()
	time.Sleep(3 * time.Second)
}
