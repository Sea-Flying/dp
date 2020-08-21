package initialize

import (
	"github.com/gocql/gocql"
	"github.com/ilyakaznacheev/cleanenv"
	"testing"
	"voyageone.com/dp/infrastructure/model/config"
	"voyageone.com/dp/infrastructure/model/global"
)

var dpConfig config.DPConfig
var artifactCqlSession gocql.Session

func TestCqlConnet(t *testing.T) {
	_ = cleanenv.ReadConfig("dp.yml", &dpConfig)
	initCqlSession(global.DPConfig.Cassandra)
}
