package service

import (
	"fmt"
	"github.com/hashicorp/nomad/api"
	"voyageone.com/dp/infrastructure/model/customType"
	"voyageone.com/dp/infrastructure/model/global"
)

func ParseJob(client *api.Client, jobHCL string) (*api.Job, error) {
	job, err := client.Jobs().ParseHCL(jobHCL, true)
	if err != nil {
		return &api.Job{}, err
	}
	return job, nil
}

func PlanJob(client *api.Client, job *api.Job) error {
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
		global.DPLogger.Printf("nomad Job plan success with warning: %s\n", planResp.Warnings)
	}
	return nil
}

func SubmitJob(client *api.Client, job *api.Job) error {
	regResp, _, err := client.Jobs().Register(job, nil)
	if err != nil {
		return customType.DPError(fmt.Sprintf("Nomad Job Run: Failed! \n %#v\n %#v\n", regResp, err))
	}
	return nil
}

func GetJobLastDeploymentHealth(client *api.Client, jobId string) (healthy bool, err error) {
	var nomadJobsEndpoint = client.Jobs()
	scaleStatusResp, _, err := nomadJobsEndpoint.ScaleStatus(jobId, nil)
	if err != nil {
		healthy = false
		return
	}
	if scaleStatusResp.TaskGroups["group"].Placed == scaleStatusResp.TaskGroups["group"].Desired &&
		scaleStatusResp.TaskGroups["group"].Healthy == scaleStatusResp.TaskGroups["group"].Placed {
		return true, nil
	} else {
		return false, nil
	}
}
