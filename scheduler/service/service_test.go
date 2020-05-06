package service

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/imdario/mergo"
	"log"
	"os"
	"testing"
	"time"
	"voyageone.com/dp/artifactKeeper/model/repository"
	artifactService "voyageone.com/dp/artifactKeeper/service"
	"voyageone.com/dp/infrastructure/model/global"
	repository2 "voyageone.com/dp/scheduler/model/repository"
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
}

func TestCreateOrUpdateDeployer(t *testing.T) {
	var deployer = repository2.Deployer{
		Name:                    "vo_local_nomad",
		Kind:                    "nomad",
		Description:             "VO local staging Nomad cluster",
		WebUrl:                  "http://10.0.0.152:4646/ui",
		NomadApiUrl:             "http://10.0.0.152:4646",
		NomadRegion:             "vo_local",
		NomadTemplateGitBaseUrl: "http://10.0.0.85/nomad-jobs/nomoad-jobs-template",
		NomadHclGitBaseUrl:      "http://10.0.0.85/nomad-jobs",
	}
	err := CreateOrUpdateDeployer(deployer)
	if err != nil {
		t.Error(err)
	}
}

func TestGetDeployerByName(t *testing.T) {
	var d = repository2.Deployer{
		Name: "vo_local_nomad",
	}
	err := GetDeployerByName(&d)
	if err != nil {
		t.Error(err)
	} else {
		global.DPLogger.Println(d)
	}
}

func TestCreateOrUpdataTemplate(t *testing.T) {
	var nt = repository2.NomadTemplate{
		Name:        "vo_local_nomad_jar",
		GitUrl:      "http://10.0.0.85/nomad-jobs/nomoad-jobs-template/vo_local_nomad_jar.git",
		GroupTags:   nil,
		ProfileTags: nil,
		Params: map[string]string{
			"Group":            "",
			"Profile":          "",
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
	err := CreateOrUpdataTemplate(nt)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateOrUpdateJob(t *testing.T) {
	var j = repository2.DPJob{
		Id:                  gocql.TimeUUID(),
		Kind:                "immediate",
		Group:               "voerp",
		Profile:             "staging",
		ClassName:           "openvms-restapi-dp",
		EntityGeneratedTime: time.Now(),
		EntityVersion:       "v202003241445",
		CreatedTime:         time.Now(),
		DeployerName:        "vo_local_nomad",
		NomadTemplateName:   "vo_local_nomad_jar",
		NomadTemplateParams: map[string]string{},
	}
	var e repository.Entity = repository.Entity{
		Group:     j.Group,
		ClassName: j.ClassName,
		Version:   j.EntityVersion,
		Profile:   j.Profile,
	}
	err := artifactService.GetEntityByVersionPartitionKey(&e)
	if err != nil {
		t.Error(err)
	}
	j.NomadTemplateParams["Group"] = j.Group
	j.NomadTemplateParams["Profile"] = j.Profile
	j.NomadTemplateParams["ClassName"] = j.ClassName
	j.NomadTemplateParams["EntityUrl"] = e.Url
	j.NomadTemplateParams["CanaryNum"] = "0"
	j.NomadTemplateParams["JvmOpts:-Xms"] = "512m"
	j.NomadTemplateParams["JvmOpts:-Xmx"] = "512m"
	j.NomadTemplateParams["Resources:cpu"] = "500"
	j.NomadTemplateParams["Resources:memory"] = "1024"
	err = CreateOrUpdateJob(j)
	if err != nil {
		t.Error(err)
	}
}

func TestGetNomadTemplateByName(t *testing.T) {
	nt := repository2.NomadTemplate{
		Name: "vo_local_nomad_jar",
	}
	err := GetNomadTemplateByName(&nt)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(nt)
	}
}

func TestGetJobById(t *testing.T) {
	jid, _ := gocql.ParseUUID("9f8126c7-845d-11ea-b535-8cec4baae5bd")
	var j = repository2.DPJob{
		Id: jid,
	}
	err := GetJobById(&j)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(j)
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
