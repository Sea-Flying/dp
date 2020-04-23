package global

import (
	"github.com/gocql/gocql"
	"github.com/hashicorp/nomad/api"
	"log"
	"voyageone.com/dp/infrastructure/entity/config"
)

var (
	CqlSession  *gocql.Session
	DPConfig    config.DPConfig
	DPLogger    *log.Logger
	NomadClient *api.Client
)
