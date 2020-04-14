package init

import (
	"github.com/gocql/gocql"
	"log"
	"strings"
	"voyageone.com/dp/model"
)

func InitCqlSession(cqlSession *gocql.Session, keySpace string, config model.DPConfig) {
	urlsSlice := strings.Split(config.Cassandra.HostsUrls, ",")
	cluster := gocql.NewCluster(urlsSlice...)
	cluster.Keyspace = keySpace
	cluster.Consistency = gocql.Consistency(1)
	cluster.NumConns = 3
	var err error
	cqlSession, err = cluster.CreateSession()
	if err != nil {
		log.Panic("start", err)
		return
	}
}
