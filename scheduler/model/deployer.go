package model

import "github.com/scylladb/gocqlx/table"

var DeployerMetadata = table.Metadata{
	Name: "schedule.deployer",
	Columns: []string{"name", "kind", "description", "web_url", "nomad_api_url", "nomad_region",
		"nomad_template_git_base_url", "nomad_hcl_git_base_url"},
	PartKey: []string{"name"},
	SortKey: []string{},
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
