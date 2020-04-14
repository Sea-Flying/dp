package global

import (
	"github.com/gocql/gocql"
	"voyageone.com/dp/model"
)

var (
	CqlSession gocql.Session
	DPConfig   model.DPConfig
)
