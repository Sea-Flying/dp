package git

import (
	"github.com/gocql/gocql"
	"log"
	"os"
	"testing"
	"time"
	"voyageone.com/dp/infrastructure/model/config"
	"voyageone.com/dp/infrastructure/model/global"
	"voyageone.com/dp/scheduler/model/repository"
)

func init() {
	cluster := gocql.NewCluster("127.0.0.1:9043")
	cluster.Consistency = gocql.Quorum
	cluster.NumConns = 3
	var err error
	global.CqlSession, err = cluster.CreateSession()
	if err != nil {
		log.Fatal("CQL Session Create Failed!", err)
	}
	global.DPLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	global.DPConfig = *new(config.DPConfig)
	global.DPConfig.Gitlab.Username = "dp"
	global.DPConfig.Gitlab.Token = "bxvHskdgTkk3zUmZKwDF"
}

func TestGitCloneHclTemplate(t *testing.T) {
	jid, _ := gocql.ParseUUID("9f8126c7-845d-11ea-b535-8cec4baae5bd")
	var j = repository.DPJob{
		Id:                  jid,
		Group:               "voerp",
		Profile:             "staging",
		ClassName:           "openvms-restapi-dp",
		EntityGeneratedTime: time.Now(),
		EntityVersion:       "v202003241445",
		CreatedTime:         time.Now(),
		DeployerName:        "vo_local_nomad",
		NomadTemplateName:   "vo_local_nomad_jar",
		NomadTemplateParams: map[string]string{
			"CanaryNum":        "0",
			"ClassName":        "openvms-restapi-dp",
			"EntityUrl":        "http://10.0.0.70:4567/voerp-staging/openvms-restapi/openvms-restapi-v202003241445.jar",
			"Group":            "voerp",
			"JvmOpts:-Xms":     "512m",
			"JvmOpts:-Xmx":     "512m",
			"Profile":          "staging",
			"Resources:cpu":    "500",
			"Resources:memory": "1024"},
	}
	err := GitCloneHclTemplate(j, "D:\\tmp\\jobtpldir")
	if err != nil {
		t.Error(err)
	}
}
