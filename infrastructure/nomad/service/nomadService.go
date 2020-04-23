package service

import (
	"fmt"
	"github.com/hashicorp/nomad/api"
	"voyageone.com/dp/infrastructure/entity/customType"
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
	if planResp.JobModifyIndex == 0 {
		return customType.DPError(fmt.Sprintf("Nomad Job Plan Failed: bad JobModifyIndex value: %d\n", planResp.JobModifyIndex))
	}
	if planResp.Diff != nil {
		return customType.DPError(fmt.Sprintf("Nomad Job Plan Failed: got non-nil diff:\n %#v\n", planResp))
	}
	if planResp.Annotations == nil {
		return customType.DPError(fmt.Sprintf("Nomad Job Plan Failed: got nil annotations:\n %#v\n", planResp))
	}
	// Can make this assertion because there are no clients.
	if len(planResp.CreatedEvals) == 0 {
		return customType.DPError(fmt.Sprintf("Nomad Job Plan Failed: got no CreatedEvals:\n %#v\n", planResp))
	}
	return customType.DPError(fmt.Sprintf("Nomad Job Plan Success \n %#v\n", planResp))
}

func SubmitJob(client *api.Client, job *api.Job) error {
	regResp, _, err := client.Jobs().Register(job, nil)
	if err != nil {
		return customType.DPError(fmt.Sprintf("Nomad Job Run: Failed! \n %#v\n %#v\n", regResp, err))
	}
	return nil
}
