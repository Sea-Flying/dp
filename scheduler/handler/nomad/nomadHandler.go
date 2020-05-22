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
	schedulerRepository "voyageone.com/dp/scheduler/model/repository"
	"voyageone.com/dp/scheduler/service"
)

func SubbmitJobImmediately(c *gin.Context) {
	var job schedulerRepository.DPJob
	var err error
	_ = c.ShouldBindJSON(&job)
	// set some immediate job default
	fillDefaultsIntoImmediateJob(&job)
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
	var deployer schedulerRepository.Deployer
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
	var job schedulerRepository.DPJob
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
	var nt schedulerRepository.NomadTemplate
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

// 填充Immediate Job的特定默认值到Job对象中
func fillDefaultsIntoImmediateJob(j *schedulerRepository.DPJob) {
	j.Kind = "immediate"
	j.CreatedTime = time.Now()
	j.Status = "executing"
}

// 填充默认值到Job对象中
func fillDefaultsIntoJob(j *schedulerRepository.DPJob) (err error) {
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
	var c = artifactRepository.Class{
		Group:   e.Group,
		Name:    e.ClassName,
		Profile: e.Profile,
	}
	err = artifactService.GetClassByPrimaryKey(&c)
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
	if j.NomadTemplateName == "" {
		j.NomadTemplateName = c.DefaultNomadTemplate
	}
	if j.NomadTemplateParams["ClassName"] == "" {
		j.NomadTemplateParams["ClassName"] = j.ClassName
	}
	if j.NomadTemplateParams["EntityVersion"] == "" {
		j.NomadTemplateParams["EntityVersion"] = j.EntityVersion
	}
	if j.NomadTemplateParams["EntityUrl"] == "" {
		j.NomadTemplateParams["EntityUrl"] = e.Url
	}
	if j.DeployerName == "" {
		var dd = schedulerRepository.DefaultDeployer{
			Group: j.Group,
		}
		err = service.GetDefaultDeployerByGroup(&dd)
		if err != nil {
			return err
		}
		j.DeployerName = dd.ProfileDeployer[j.Profile]
	}
	err = service.MergeClassDefaultParamsIntoJob(j, c)
	if err != nil {
		return err
	}
	err = service.MergeTemplateDefaultParamsIntoJob(j)
	return err
}

// Nomad Job的渲染、验证和提交
func submitJobProcedure(job schedulerRepository.DPJob) (err error) {
	//git clone Job HCL template
	err = git.GitCloneHclTemplate(job, DPConfig.Nomad.NomadJobTplDir)
	if err != nil {
		return
	}
	templateDir := DPConfig.Nomad.NomadJobTplDir + "/" + job.NomadTemplateName + "/" + job.NomadTemplateName + ".tpl"
	//根据传入的template参数渲染出实际的Job HCL文本
	jobStr, err := nomadService.RenderJobTemplate(job, templateDir)
	if err != nil {
		return
	}
	//nomad parse Job HCL string to Job JSON Object
	nomadJob, err := nomadService.ParseJob(NomadClient, jobStr)
	if err != nil {
		return
	}
	//nomad plan job
	err = nomadService.PlanJob(NomadClient, nomadJob)
	if err != nil {
		return
	}
	//nomad run job
	err = nomadService.SubmitJob(NomadClient, nomadJob)
	return
}
