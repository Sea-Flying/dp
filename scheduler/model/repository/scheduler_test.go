package repository

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/imdario/mergo"
	"github.com/scylladb/gocqlx/v2"
	"log"
	"os"
	"testing"
	"time"
	artifactRepository "voyageone.com/dp/artifactKeeper/model/repository"
	"voyageone.com/dp/infrastructure/model/global"
)

func init() {
	cluster := gocql.NewCluster("127.0.0.1:9042")
	cluster.Consistency = gocql.Quorum
	cluster.NumConns = 3
	var err error
	global.CqlSession, err = gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Fatal("CQL Session Create Failed!", err)
	}
	global.DPLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func TestCreateOrUpdataTemplate(t *testing.T) {
	var nt = NomadTemplate{
		Name:   "vo_local_nomad_jar",
		GitUrl: "http://10.0.0.85/nomad-jobs/nomoad-jobs-template/vo_local_nomad_jar.git",
		Params: map[string]string{
			"ClassName":        "",
			"EntityVersion":    "",
			"EntityUrl":        "",
			"CanaryNum":        "0",
			"JvmOpts:-Xmx":     "512m",
			"JvmOpts:-Xms":     "512m",
			"Resources:cpu":    "500",
			"Resources:memory": "1024",
		},
	}
	err := nt.CreateOrUpdate()
	if err != nil {
		t.Error(err)
	}
}

func TestCreateOrUpdateJob(t *testing.T) {
	var j = DPJob{
		Kind:                "immediate",
		ClassName:           "openvms-restapi-dp",
		EntityGeneratedTime: time.Now(),
		EntityVersion:       "v202003241445",
		CreatedTime:         time.Now(),
		NomadTemplateName:   "vo_local_nomad_jar",
		NomadTemplateParams: map[string]string{},
	}
	var e = artifactRepository.Entity{
		ClassName: j.ClassName,
		Version:   j.EntityVersion,
	}
	err := e.GetByVersionPartitionKey()
	if err != nil {
		t.Error(err)
	}
	j.NomadTemplateParams["ClassName"] = j.ClassName
	j.NomadTemplateParams["EntityUrl"] = e.Url
	j.NomadTemplateParams["CanaryNum"] = "0"
	j.NomadTemplateParams["JvmOpts:-Xms"] = "512m"
	j.NomadTemplateParams["JvmOpts:-Xmx"] = "512m"
	j.NomadTemplateParams["Resources:cpu"] = "500"
	j.NomadTemplateParams["Resources:memory"] = "1024"
	err = j.CreateOrUpdate()
	if err != nil {
		t.Error(err)
	}
}

func TestGetNomadTemplateByName(t *testing.T) {
	nt := NomadTemplate{
		Name: "vo_local_nomad_jar",
	}
	err := nt.GetByName()
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(nt)
	}
}

func TestMergeTemplateDefaultParamsIntoJob(t *testing.T) {
	a := map[string]string{
		"a": "1",
		"b": "2",
	}
	b := map[string]string{
		"b": "3",
		"c": "4",
	}
	mergo.Merge(&a, b)
	fmt.Println(a)
}
