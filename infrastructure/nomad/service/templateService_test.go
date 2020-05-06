package service

import (
	"fmt"
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
	global.DPLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	global.DPConfig = *new(config.DPConfig)
	global.DPConfig.Gitlab.Username = "dp"
	global.DPConfig.Gitlab.Token = "bxvHskdgTkk3zUmZKwDF"
}

func TestRenderJobTemplate(t *testing.T) {
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
	s, err := RenderJobTemplate(j, "D:/tmp/jobtpldir/vo_local_nomad_jar/vo_local_nomad_jar.tpl")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(s)
	}
}
