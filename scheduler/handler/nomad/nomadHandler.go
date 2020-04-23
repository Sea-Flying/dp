package nomad

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/hashicorp/nomad/api"
	. "voyageone.com/dp/infrastructure/entity/global"
	"voyageone.com/dp/infrastructure/entity/response"
	"voyageone.com/dp/infrastructure/git"
	nomadService "voyageone.com/dp/infrastructure/nomad/service"
	"voyageone.com/dp/scheduler/model"
	"voyageone.com/dp/scheduler/service"
)

func SubbmitJobImmediately(c *gin.Context) {
	var job model.DPJob
	var err error
	var templateDir, jobStr string
	var nomadJob *api.Job
	_ = c.ShouldBindJSON(&job)
	job.Id = gocql.TimeUUID()
	err = service.CreateOrUpdateJob(job)
	if err != nil {
		goto ErrorReturn
	}
	err = git.GitCloneHclTemplate(job, DPConfig.Nomad.NomadJobTplDir)
	if err != nil {
		goto ErrorReturn
	}
	templateDir = DPConfig.Nomad.NomadJobTplDir + job.NomadTemplateName
	jobStr, err = nomadService.RenderJobTemplate(job, templateDir)
	if err != nil {
		goto ErrorReturn
	}
	nomadJob, err = nomadService.ParseJob(NomadClient, jobStr)
	if err != nil {
		goto ErrorReturn
	}
	err = nomadService.PlanJob(NomadClient, nomadJob)
	if err != nil {
		goto ErrorReturn
	}
	err = nomadService.SubmitJob(NomadClient, nomadJob)
	if err != nil {
		goto ErrorReturn
	} else {
		DPLogger.Printf("Imediate Nomad Job Create Success! \nJob Object: %#v", job)
		response.OkWithMessage("Imediate Nomad Job Create Success!", c)
		return
	}
ErrorReturn:
	{
		DPLogger.Printf("Immediate Nomad Job Create Failed! \nJob Object: %#v \nError: %#v", job, err)
		response.FailWithMessage(fmt.Sprintf("Immediate Nomad Job Create Failed! \n%#v", err), c)
		return
	}
}

func CreateDeployer(c *gin.Context) {
	var deployer model.Deployer
	_ = c.ShouldBindJSON(deployer)

}

func CreateOrUpdateJob(c *gin.Context) {

}

func CreateOrUpdateJobTemplate(c *gin.Context) {

}
