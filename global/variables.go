package global

import (
	"github.com/gocql/gocql"
	"voyageone.com/dp/model"
)

var (
	ArtifactCqlSession   gocql.Session
	ManagementCqlSession gocql.Session
	DPConfig             model.DPConfig
)
