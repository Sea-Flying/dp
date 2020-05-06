package repository

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
	"time"
)

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
