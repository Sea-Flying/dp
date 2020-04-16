package global

import (
	"github.com/gocql/gocql"
	"log"
	"voyageone.com/dp/infrastructure/entity/config"
)

var (
	CqlSession *gocql.Session
	DPConfig   config.DPConfig
	DPLogger   log.Logger
)
