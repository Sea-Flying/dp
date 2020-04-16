package initialize

import (
	"github.com/gocql/gocql"
	"log"
	"strings"
	"voyageone.com/dp/infrastructure/entity/config"
	"voyageone.com/dp/infrastructure/entity/global"
)

func InitCqlSession(config config.DPConfig) {
	urlsSlice := strings.Split(config.Cassandra.HostsUrls, ",")
	cluster := gocql.NewCluster(urlsSlice...)
	cluster.Consistency = gocql.Quorum
	cluster.NumConns = 3
	var err error
	global.CqlSession, err = cluster.CreateSession()
	if err != nil {
		log.Panic("CQL Session Create Failed!", err)
		return
	}
}
