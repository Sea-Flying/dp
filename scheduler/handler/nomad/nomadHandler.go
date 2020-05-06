package nomad

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"time"
	artifactRepository "voyageone.com/dp/artifactKeeper/model/repository"
	artifactService "voyageone.com/dp/artifactKeeper/service"
	"voyageone.com/dp/infrastructure/git"
	. "voyageone.com/dp/infrastructure/model/global"
	"voyageone.com/dp/infrastructure/model/response"
	nomadService "voyageone.com/dp/infrastructure/nomad/service"
	"voyageone.com/dp/scheduler/model/repository"
	"voyageone.com/dp/scheduler/service"
)

func SubbmitJobImmediately(c *gin.Context) {
	var job repository.DPJob
	var err error
	_ = c.ShouldBindJSON(&job)
	{
		job.Kind = "immediate"
		job.CreatedTime = time.Now()
		job.Status = "executing"
	}
	err = fillDefaultsIntoJob(&job)
	if err != nil {
		goto ErrorReturn
	}
	err = service.CreateOrUpdateJob(job)
	if err != nil {
		goto ErrorReturn
	}
	err = submitJobProcedure(job)
	if err != nil {
		goto ErrorReturn
	} else {
		job.Status = "success"
		_ = service.CreateOrUpdateJob(job)
		DPLogger.Printf("Imediate Nomad Job Create Success \nJob Object: %#v", job)
		response.OkWithMessage("Imediate Nomad Job Create Success ", c)
		return
	}
ErrorReturn:
	{
		job.Status = "failed"
		_ = service.CreateOrUpdateJob(job)
		DPLogger.Printf("Immediate Nomad Job Create Failed \nJob Object: %#v \nError: %#v\n", job, err)
		response.FailWithMessage(fmt.Sprintf("Immediate Nomad Job Create Failed: \n%#v\n", err), c)
		return
	}
}

func CreateOrUpdateDeployer(c *gin.Context) {
	var deployer repository.Deployer
	_ = c.ShouldBindJSON(deployer)
	err := service.CreateOrUpdateDeployer(deployer)
	if err != nil {
		DPLogger.Printf("Create Or Updata Deployer Failed \nEntity Object: %#v \nError: %#v", deployer, err)
		response.FailWithMessage(fmt.Sprintf("Create Or Updata Deployer Failed: \n%#v\n", err), c)
	} else {
		DPLogger.Printf("Create Or Updata Deployer Success \nEntity Object: %#v\n", deployer)
		response.OkWithMessage("Create Or Updata Deployer Success", c)
	}
}

func CreateOrUpdateJob(c *gin.Context) {
	var job repository.DPJob
	_ = c.ShouldBindJSON(job)
	err := service.CreateOrUpdateJob(job)
	if err != nil {
		DPLogger.Printf("Create Or Updata Nomad Job Failed \njob Object: %#v \nError: %#v", job, err)
		response.FailWithMessage(fmt.Sprintf("Create Or Updata Nomad Job Failed: \n%#v\n", err), c)
	} else {
		DPLogger.Printf("Create Or Updata Nomad Job Success \nJob Object: %#v", job)
		response.OkWithMessage("Create Or Updata Nomad Job Success", c)
	}
}

func CreateOrUpdateJobTemplate(c *gin.Context) {
	var nt repository.NomadTemplate
	_ = c.ShouldBindJSON(nt)
	err := service.CreateOrUpdataTemplate(nt)
	if err != nil {
		DPLogger.Printf("Create Or Updata Nomad JobTemplate Failed \njob Object: %#v \nError: %#v", nt, err)
		response.FailWithMessage(fmt.Sprintf("Create Or Updata Nomad JobTemplate Failed: \n%#v\n", err), c)
	} else {
		DPLogger.Printf("Create Or Updata Nomad JobTemplate Success \nJob Object: %#v", nt)
		response.OkWithMessage("Create Or Updata Nomad JobTemplate Success", c)
	}
}

func CheckJobLastDeploymentStatus(c *gin.Context) {
	jobId := c.Param("jobId")
	healthy, err := nomadService.GetJobLastDeploymentHealth(NomadClient, jobId)
	if err != nil {
		DPLogger.Printf("JobId: [%s] , CheckJobLastDeploymentStatus error : %#v", jobId, err)
		response.FailWithMessage(fmt.Sprintf("JobId: [%s] , CheckJobLastDeploymentStatus error : %#v", jobId, err), c)
		return
	}
	if healthy {
		response.OkWithMessage(fmt.Sprintf("JobId: [%s] is healthy now", jobId), c)
	} else {
		response.FailWithMessage(fmt.Sprintf("JobId: [%s] is unhealthy now", jobId), c)
	}
}

func fillDefaultsIntoJob(j *repository.DPJob) (err error) {
	j.Id = gocql.TimeUUID()
	e := artifactRepository.Entity{
		Group:     j.Group,
		ClassName: j.ClassName,
		Version:   j.EntityVersion,
		Profile:   j.Profile,
	}
	err = artifactService.GetEntityByVersionPartitionKey(&e)
	if err != nil {
		return err
	}
	j.EntityGeneratedTime = e.GeneratedTime
	if j.NomadTemplateParams == nil {
		j.NomadTemplateParams = map[string]string{}
	}
	if j.Kind == "immediate" {
		j.ExecutedTime = j.CreatedTime
	}
	if j.NomadTemplateParams["ClassName"] == "" {
		j.NomadTemplateParams["ClassName"] = j.ClassName
	}
	if j.NomadTemplateParams["EntityVersion"] == "" {
		j.NomadTemplateParams["EntityVersion"] = j.EntityVersion
	}
	if j.NomadTemplateParams["EntityUrl"] == "" {
		e := artifactRepository.Entity{
			Group:     j.Group,
			Profile:   j.Profile,
			ClassName: j.ClassName,
			Version:   j.EntityVersion,
		}
		err := artifactService.GetEntityByVersionPartitionKey(&e)
		if err != nil {
			return err
		}
		j.NomadTemplateParams["EntityUrl"] = e.Url
	}
	err = service.MergeTemplateDefaultParamsIntoJob(j)
	return err
}

func submitJobProcedure(job repository.DPJob) (err error) {
	err = git.GitCloneHclTemplate(job, DPConfig.Nomad.NomadJobTplDir)
	if err != nil {
		return
	}
	templateDir := DPConfig.Nomad.NomadJobTplDir + "/" + job.NomadTemplateName + "/" + job.NomadTemplateName + ".tpl"
	jobStr, err := nomadService.RenderJobTemplate(job, templateDir)
	if err != nil {
		return
	}
	nomadJob, err := nomadService.ParseJob(NomadClient, jobStr)
	if err != nil {
		return
	}
	err = nomadService.PlanJob(NomadClient, nomadJob)
	if err != nil {
		return
	}
	err = nomadService.SubmitJob(NomadClient, nomadJob)
	return
}
