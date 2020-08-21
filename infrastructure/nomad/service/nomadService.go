package service

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/nomad/api"
	"log"
	"time"
	"voyageone.com/dp/infrastructure/model/customType"
)

type VoNomadClient struct {
	*api.Client
	Logger *log.Logger
}

func NewVoNomadClient(client *api.Client, logger *log.Logger) *VoNomadClient {
	return &VoNomadClient{
		Client: client,
		Logger: logger,
	}
}

func (client *VoNomadClient) ParseJob(jobHCL string) (*api.Job, error) {
	job, err := client.Jobs().ParseHCL(jobHCL, true)
	if err != nil {
		return &api.Job{}, err
	}
	return job, nil
}

func (client *VoNomadClient) PlanJob(job *api.Job) error {
	planResp, _, err := client.Jobs().Plan(job, true, nil)
	if err != nil {
		return customType.DPError("Nomad Job Plan Failed " + err.Error())
	}
	if len(planResp.FailedTGAllocs) > 0 {
		return customType.DPError(fmt.Sprintf("Nomad Job Plan Failed: got failed Allocs :\n %#v\n", planResp))
	}
	if planResp.Annotations == nil {
		return customType.DPError(fmt.Sprintf("Nomad Job Plan Failed: got nil annotations:\n %#v\n", planResp))
	}
	if planResp.Warnings != "" {
		client.Logger.Printf("nomad Job plan success with warnig: %s\n", planResp.Warnings)
	}
	return nil
}

func (client *VoNomadClient) SubmitJob(job *api.Job) error {
	regResp, _, err := client.Jobs().Register(job, nil)
	if err != nil {
		return customType.DPError(fmt.Sprintf("Nomad Job Run: Failed! \n %#v\n %#v\n", regResp, err))
	}
	return nil
}

func (client *VoNomadClient) GetJob(jobId string) (job *api.Job, err error) {
	var nomadJobsEndpoint = client.Jobs()
	job, _, err = nomadJobsEndpoint.Info(jobId, nil)
	return
}

func (client *VoNomadClient) GetJobJson(jobId string) (jobJson string, err error) {
	var nomadJobsEndpoint = client.Jobs()
	job, _, err := nomadJobsEndpoint.Info(jobId, nil)
	if err != nil {
		return "", err
	}
	jobJsonBytes, _ := json.Marshal(job)
	jobJson = string(jobJsonBytes)
	return
}

func (client *VoNomadClient) GetJobsList() (jobs []string, err error) {
	var nomadJobsEndpoint = client.Jobs()
	jobs = make([]string, 0, 100)
	JobListStub, _, err := nomadJobsEndpoint.List(nil)
	if err != nil {
		return nil, err
	}
	for _, job := range JobListStub {
		jobs = append(jobs, job.ID)
	}
	return jobs, err
}

func (client *VoNomadClient) GetJobLastDeploymentStatus(jobId string) string {
	var nomadJobsEndpoint = client.Jobs()
	latestDeployment, _, err := nomadJobsEndpoint.LatestDeployment(jobId, nil)
	if err != nil {
		return ""
	}
	if latestDeployment != nil {
		return latestDeployment.Status
	} else {
		return ""
	}
}

func (client *VoNomadClient) GetJobLastDeploymentHealth(jobId string) (healthy bool, err error) {
	var nomadJobsEndpoint = client.Jobs()
	latestDeployment, _, err := nomadJobsEndpoint.LatestDeployment(jobId, nil)
	if err != nil {
		healthy = false
		return
	}
	job, _, err := nomadJobsEndpoint.Info(jobId, nil)
	if err != nil {
		healthy = false
		return
	}
	if *(job.Status) == "running" && latestDeployment.Status == "successful" {
		return true, nil
	} else {
		return false, nil
	}
}

func (client *VoNomadClient) FailJobLastDeployment(jobId string) (err error) {
	var nomadJobsEndpoint = client.Jobs()
	latestDeployment, _, err := nomadJobsEndpoint.LatestDeployment(jobId, nil)
	if err != nil {
		return
	}
	var nomadDeploymentsEndpoint = client.Deployments()
	nomadDeploymentsEndpoint.Fail(latestDeployment.ID, nil)
	if err != nil {
		return
	}
	return nil
}

func (client *VoNomadClient) GetLastDeploymentAllocations(jobId string) (allocationsId []string, err error) {
	var nomadJobsEndpoint = client.Jobs()
	latestDeployment, _, err := nomadJobsEndpoint.LatestDeployment(jobId, nil)
	if err != nil {
		return
	}
	var nomadDeploymentsEndpoint = client.Deployments()
	allocations, _, err := nomadDeploymentsEndpoint.Allocations(latestDeployment.ID, nil)
	for _, a := range allocations {
		allocationsId = append(allocationsId, a.ID)
	}
	return
}

func (client *VoNomadClient) StopJob(jobId string) (err error) {
	var nomadJobsEndpoint = client.Jobs()
	_, _, err = nomadJobsEndpoint.Deregister(jobId, false, nil)
	return
}

func (client *VoNomadClient) RestartJob(jobId string) (err error) {
	j, err := client.GetJobJson(jobId)
	if err != nil {
		return
	}
	var J api.Job
	var canary = 1
	var autoPromote = true
	var autoRevert = true
	json.Unmarshal([]byte(j), &J)
	J.TaskGroups[0].Update.Canary = &canary
	J.TaskGroups[0].Update.AutoPromote = &autoPromote
	J.TaskGroups[0].Update.AutoRevert = &autoRevert
	if J.Meta == nil {
		metas := map[string]string{"manual_restart": time.Now().String()}
		J.Meta = metas
	} else {
		J.Meta["manual_restart"] = time.Now().String()
	}
	err = client.SubmitJob(&J)
	return
}
