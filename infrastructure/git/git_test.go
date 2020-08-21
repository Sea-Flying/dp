package git

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"log"
	"os"
	"testing"
	"time"
	"voyageone.com/dp/infrastructure/model/config"
	"voyageone.com/dp/infrastructure/model/global"
	"voyageone.com/dp/scheduler/model/repository"
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
	global.DPLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile)
	global.DPConfig = *new(config.DPConfig)
	global.DPConfig.Gitlab.Username = "dp"
	global.DPConfig.Gitlab.Token = "bxvHskdgTkk3zUmZKwDF"
}

func TestGitCloneHclTemplate(t *testing.T) {
	var j = repository.DPJob{
		ClassName:           "openvms-restapi-dp",
		EntityGeneratedTime: time.Now(),
		EntityVersion:       "v202003241445",
		CreatedTime:         time.Now(),
		NomadTemplateName:   "vo_local_nomad_jar",
		NomadTemplateParams: map[string]string{
			"CanaryNum":        "0",
			"ClassName":        "openvms-restapi-dp",
			"EntityUrl":        "http://10.0.0.70:4567/voerp-staging/openvms-restapi/openvms-restapi-v202003241445.jar",
			"JvmOpts:-Xms":     "512m",
			"JvmOpts:-Xmx":     "512m",
			"Resources:cpu":    "500",
			"Resources:memory": "1024"},
	}
	_, err := CloneHclTemplate(j, "D:\\tmp\\jobtpldir")
	if err != nil {
		t.Error(err)
	}
}
