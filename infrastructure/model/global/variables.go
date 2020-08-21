package global

import (
	"github.com/scylladb/gocqlx/v2"
	"gopkg.in/olahol/melody.v1"
	"log"
	"sync"
	consulService "voyageone.com/dp/infrastructure/consul/service"
	"voyageone.com/dp/infrastructure/model/config"
	nomadService "voyageone.com/dp/infrastructure/nomad/service"
)

var (
	CqlSession   gocqlx.Session
	DPConfig     config.DPConfig
	DPLogger     *log.Logger
	NomadClient  *nomadService.VoNomadClient
	ConsulClient *consulService.VoConsulClient
	WsAppsStatus *melody.Melody
	GitMutex     = make(map[string]*sync.Mutex)
)
