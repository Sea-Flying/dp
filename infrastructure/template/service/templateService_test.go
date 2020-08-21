package service

import (
	"fmt"
	"testing"
	"time"
	"voyageone.com/dp/scheduler/model/repository"
)

func TestRenderJobTemplate(t *testing.T) {
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
	s, err := RenderJobTemplate(j, "D:/tmp/jobtpldir/vo_local_nomad_jar/vo_local_nomad_jar.tpl")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(s)
	}
}
