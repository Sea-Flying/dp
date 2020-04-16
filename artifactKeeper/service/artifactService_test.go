package service

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"testing"
	"time"
	"voyageone.com/dp/artifactKeeper/model/repository"
	"voyageone.com/dp/infrastructure/entity/global"
)

func TestCreateArtifactClass(t *testing.T) {
	var c = repository.Class{
		Name:        "voerp-restapi",
		Group:       "voerp",
		Profile:     "prod",
		Kind:        "jar",
		RepoName:    "hk_oss",
		CreatedTime: time.Now(),
		GitUrl:      "https://github.com",
		CiUrl:       "",
		Description: "",
	}

	cluster := gocql.NewCluster("127.0.0.1:9043")
	cluster.Consistency = gocql.Quorum
	cluster.NumConns = 3
	var err error
	global.CqlSession, err = cluster.CreateSession()
	if err != nil {
		log.Panic("CQL Session Create Failed!", err)
		return
	}
	err = CreateOrUpdateClass(c)
	t.Log(err)
}

func TestReadClassByPrimaryKey(t *testing.T) {
	var c = repository.Class{
		Group:   "voerp",
		Profile: "staging",
		Name:    "voerp-product",
	}
	cluster := gocql.NewCluster("127.0.0.1:9043")
	cluster.Consistency = gocql.Quorum
	cluster.NumConns = 3
	var err error
	global.CqlSession, err = cluster.CreateSession()
	if err != nil {
		log.Panic("CQL Session Create Failed!", err)
		return
	}
	err = ReadClassByPrimaryKey(&c)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(c)
	}
}

func TestCreateOrUpdateEntity(t *testing.T) {
	var e = repository.Entity{
		Group:         "voerp",
		Profile:       "dev",
		ClassName:     "voerp-product",
		ClassKind:     "jar",
		GeneratedTime: time.Now(),
		RepoName:      "local_http",
		Uploader:      "jenkins",
		Version:       "v202004161644",
	}
	cluster := gocql.NewCluster("127.0.0.1:9043")
	cluster.Consistency = gocql.Quorum
	cluster.NumConns = 3
	var err error
	global.CqlSession, err = cluster.CreateSession()
	if err != nil {
		log.Panic("CQL Session Create Failed!", err)
		return
	}
	err = CreateOrUpdateEntity(e)
	if err != nil {
		t.Error(err)
	}
}

func TestReadEntityByVersionPartitionKey(t *testing.T) {
	var e = repository.Entity{
		Group:     "voerp",
		Profile:   "dev",
		ClassName: "voerp-product",
		Version:   "v202004161644",
	}
	cluster := gocql.NewCluster("127.0.0.1:9043")
	cluster.Consistency = gocql.Quorum
	cluster.NumConns = 3
	var err error
	global.CqlSession, err = cluster.CreateSession()
	if err != nil {
		log.Panic("CQL Session Create Failed!", err)
		return
	}
	err = ReadEntityByVersionPartitionKey(&e)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(e)
	}
}
