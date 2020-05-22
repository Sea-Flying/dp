package repository

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
	"time"
)

var DeployerMetadata = table.Metadata{
	Name: "schedule.deployer",
	Columns: []string{"name", "kind", "description", "web_url", "nomad_api_url", "nomad_region",
		"nomad_template_git_base_url", "nomad_hcl_git_base_url"},
	PartKey: []string{"name"},
	SortKey: nil,
}

type Deployer struct {
	Name                    string
	Kind                    string
	Description             string
	WebUrl                  string
	NomadApiUrl             string `json:"nomad_api_url"`
	NomadRegion             string `json:"nomad_region"`
	NomadTemplateGitBaseUrl string `json:"nomad_template_git_base_url"`
	NomadHclGitBaseUrl      string `json:"nomad_hcl_git_base_url"`
}

var DefaultDeployerMetadata = table.Metadata{
	Name:    "schedule.default_deployer",
	Columns: []string{"group", "profile_deployer"},
	PartKey: []string{"group"},
	SortKey: nil,
}

type DefaultDeployer struct {
	Group           string            `json:"group"`
	ProfileDeployer map[string]string `json:"profile_deployer"`
}

var DPJobMetadata = table.Metadata{
	Name: "schedule.job",
	Columns: []string{"id", "kind", "group", "profile", "class_name", "entity_generated_time", "entity_version", "created_time",
		"executed_time", "deployer_name", "nomad_template_name", "nomad_template_params", "status"},
	PartKey: []string{"group", "profile", "class_name", "entity_generated_time"},
	SortKey: []string{"created_time"},
}

type DPJob struct {
	Id                  gocql.UUID        `json:"id"`
	Kind                string            `json:"kind"`
	Group               string            `json:"group"`
	Profile             string            `json:"profile"`
	ClassName           string            `json:"class_name"`
	EntityGeneratedTime time.Time         `json:"entity_generated_time"`
	EntityVersion       string            `json:"entity_version"`
	CreatedTime         time.Time         `json:"created_time"`
	ExecutedTime        time.Time         `json:"executed_time"`
	DeployerName        string            `json:"deployer_name"`
	NomadTemplateName   string            `json:"nomad_template_name"`
	NomadTemplateParams map[string]string `json:"nomad_template_params"`
	Status              string            `json:"status"`
}

var NomadTemplateMetadata = table.Metadata{
	Name:    "schedule.nomad_template",
	Columns: []string{"name", "git_url", "group_tags", "profile_tags", "params"},
	PartKey: []string{"name"},
	SortKey: nil,
}

type NomadTemplate struct {
	Name              string            `json:"name"`
	GitUrl            string            `json:"git_url"`
	GroupTags         []string          `json:"group_tags"`
	ProfileTags       []string          `json:"profile_tags"`
	Params            map[string]string `json:"params"`
	ParamsDescription map[string]string `json:"params_description"`
}
