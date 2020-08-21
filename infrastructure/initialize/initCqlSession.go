package initialize

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"log"
	"strings"
	"time"
	"voyageone.com/dp/infrastructure/model/config"
	"voyageone.com/dp/infrastructure/model/global"
)

func InitCqlSession() {
	var err error
	global.CqlSession, err = initCqlSession(global.DPConfig.Cassandra)
	if err != nil {
		log.Fatal(err)
	}
}

func initCqlSession(cassandraConfig config.CassandraConfig) (gocqlx.Session, error) {
	urlsSlice := strings.Split(cassandraConfig.HostsUrls, ",")
	cluster := gocql.NewCluster(urlsSlice...)
	cluster.Consistency = gocql.Quorum
	cluster.NumConns = 3
	cluster.ConnectTimeout = 1 * time.Second
	cluster.Timeout = 2 * time.Second
	return gocqlx.WrapSession(cluster.CreateSession())
}
