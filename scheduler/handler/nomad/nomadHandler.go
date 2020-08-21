package nomad

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/imdario/mergo"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	watcherService "voyageone.com/dp/app/service/watcher"
	artifactRepository "voyageone.com/dp/artifactKeeper/model/repository"
	"voyageone.com/dp/infrastructure/git"
	"voyageone.com/dp/infrastructure/model/customType"
	. "voyageone.com/dp/infrastructure/model/global"
	"voyageone.com/dp/infrastructure/model/response"
	templateService "voyageone.com/dp/infrastructure/template/service"
	schedulerRepository "voyageone.com/dp/scheduler/model/repository"
)

func SubmitJobImmediately(c *gin.Context) {
	var job schedulerRepository.DPJob
	var err error
	_ = c.ShouldBindJSON(&job)
	// set some immediate job default
	fillDefaultsIntoImmediateJob(&job)
	err = fillDefaultsIntoJob(&job)
	if err != nil {
		goto ErrorReturn
	}
	err = job.CreateOrUpdate()
	if err != nil {
		goto ErrorReturn
	}
	err = submitJobProcedure(job)
	if err != nil {
		goto ErrorReturn
	} else {
		job.Status = "success"
		_ = job.CreateOrUpdate()
		DPLogger.Printf("Imediate Nomad Job Create Success \nJob Object: %#v", job)
		response.OkWithMessage("Imediate Nomad Job Create Success ", c)
		return
	}
ErrorReturn:
	{
		job.Status = "failed"
		_ = job.CreateOrUpdate()
		DPLogger.Printf("Immediate Nomad Job Create Failed \nJob Object: %#v \nError: %#v\n", job, err)
		response.FailWithMessage(fmt.Sprintf("Immediate Nomad Job Create Failed: \n%#v\n", err), c)
		return
	}
}

func CreateOrUpdateJob(c *gin.Context) {
	var job schedulerRepository.DPJob
	_ = c.ShouldBindJSON(job)
	err := job.CreateOrUpdate()
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
	err := nt.CreateOrUpdate()
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
	if watcherService.AppWatchList[jobId].Current() == "healthy" {
		response.OkWithMessage(fmt.Sprintf("JobId: [%s] is healthy now", jobId), c)
	} else {
		response.FailWithMessage(fmt.Sprintf("JobId: [%s] is unhealthy now", jobId), c)
	}
}

func GetJobLastDeploymentAllocsId(c *gin.Context) {
	jobId := c.Param("jobId")
	allocsId, err := NomadClient.GetLastDeploymentAllocations(jobId)
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Writer.Write([]byte(fmt.Sprintf(`get last deployment allocs id failed, %v`, err)))
		return
	}
	var ret string
	for _, allocId := range allocsId {
		ret += strings.Split(allocId, "-")[0] + "|"
	}
	c.Status(http.StatusOK)
	c.Writer.Write([]byte(ret))
}

// 填充Immediate Job的特定默认值到Job对象中
func fillDefaultsIntoImmediateJob(j *schedulerRepository.DPJob) {
	j.Kind = "immediate"
	j.CreatedTime = time.Now()
	j.Status = "executing"
}

// 填充默认值到Job对象中
func fillDefaultsIntoJob(j *schedulerRepository.DPJob) (err error) {
	e := artifactRepository.Entity{
		ClassName: j.ClassName,
		Version:   j.EntityVersion,
	}
	err = e.GetByVersionPartitionKey()
	if err != nil {
		return err
	}
	var c = artifactRepository.Class{
		Name: e.ClassName,
	}
	err = c.GetByPrimaryKey()
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
	if j.UnitTimeoutSeconds == 0 {
		j.UnitTimeoutSeconds = c.UnitTimeoutSeconds
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
	err = mergeClassDefaultParamsIntoJob(j, c)
	if err != nil {
		return err
	}
	err = mergeTemplateDefaultParamsIntoJob(j)
	return err
}

func mergeClassDefaultParamsIntoJob(j *schedulerRepository.DPJob, c artifactRepository.Class) error {
	err := mergo.Merge(&j.NomadTemplateParams, c.DefaultTemplateParams)
	if err != nil {
		return customType.DPError("merge class defaults to job failed " + err.Error())
	}
	return nil
}

func mergeTemplateDefaultParamsIntoJob(j *schedulerRepository.DPJob) error {
	var nt = schedulerRepository.NomadTemplate{
		Name: j.NomadTemplateName,
	}
	err := nt.GetByName()
	if err != nil {
		return customType.DPError("merge template defaults to job failed " + err.Error())
	}
	err = mergo.Merge(&j.NomadTemplateParams, nt.Params)
	if err != nil {
		return customType.DPError("merge template defaults to job failed " + err.Error())
	}
	return nil
}

// Nomad Job的渲染、验证和提交
func submitJobProcedure(job schedulerRepository.DPJob) (err error) {
	//git clone Job HCL template
	gitRepoDir, err := git.CloneHclTemplate(job, DPConfig.Nomad.NomadJobTplDir)
	if err != nil {
		return
	}
	templateDir := gitRepoDir + "/" + job.NomadTemplateName + ".tpl"
	//根据传入的template参数渲染出实际的Job HCL文本
	jobStr, err := templateService.RenderJobTemplate(job, templateDir)
	if err != nil {
		return
	}
	//nomad parse Job HCL string to Job JSON Object
	nomadJob, err := NomadClient.ParseJob(jobStr)
	if err != nil {
		return
	}
	//nomad plan job
	err = NomadClient.PlanJob(nomadJob)
	if err != nil {
		return
	}
	//nomad run job
	err = NomadClient.SubmitJob(nomadJob)
	//set AppWatcher status and settings
	if _, existed := watcherService.AppWatchList[job.ClassName]; !existed {
		watcherService.AddApp(job.ClassName)
	}
	watcherService.AppWatchList[job.ClassName].UnitTimeoutSeconds = job.UnitTimeoutSeconds
	var parallel int
	if job.NomadTemplateParams["Canary"] == "0" {
		parallel, _ = strconv.Atoi(job.NomadTemplateParams["MaxParallel"])
	} else {
		parallel, _ = strconv.Atoi(job.NomadTemplateParams["Canary"])
	}
	count, _ := strconv.Atoi(job.NomadTemplateParams["Count"])
	factor := int(math.Ceil(float64(count) / float64(parallel)))
	watcherService.AppWatchList[job.ClassName].TimeoutFactor = factor
	err = watcherService.AppWatchList[job.ClassName].Deploy()
	return
}
