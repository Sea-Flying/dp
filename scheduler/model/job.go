package model

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
	"time"
)

var DPJobMetadata = table.Metadata{
	Name: "schedule.job",
	Columns: []string{"id", "group", "profile", "class_name", "entity_generated_time", "entity_version", "created_time",
		"deployer_name", "nomad_template_name", "nomad_template_params"},
	PartKey: []string{"group", "profile", "class_name", "entity_generated_time"},
	SortKey: []string{"created_time"},
}

type DPJob struct {
	Id                  gocql.UUID        `json:"id"`
	Group               string            `json:"group"`
	Profile             string            `json:"profile"`
	ClassName           string            `json:"class_name"`
	EntityGeneratedTime time.Time         `json:"entity_generated_time"`
	EntityVersion       string            `json:"entity_version"`
	CreatedTime         time.Time         `json:"created_time"`
	DeployerName        string            `json:"deployer_name"`
	NomadTemplateName   string            `json:"nomad_template_name"`
	NomadTemplateParams map[string]string `json:"nomad_template_params"`
}
