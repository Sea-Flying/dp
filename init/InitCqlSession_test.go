package init

import (
	"github.com/gocql/gocql"
	"github.com/ilyakaznacheev/cleanenv"
	"testing"
	"voyageone.com/dp/model"
)

var dpConfig model.DPConfig
var artifactCqlSession gocql.Session

func TestCqlConnet(t *testing.T) {
	_ = cleanenv.ReadConfig("dp.yml", &dpConfig)
	InitCqlSession(&artifactCqlSession, "artifact", dpConfig)
}
