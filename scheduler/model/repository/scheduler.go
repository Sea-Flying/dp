package repository

import (
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
	"time"
	. "voyageone.com/dp/infrastructure/model/global"
)

var DPJobMetadata = table.Metadata{
	Name: "schedule.job",
	Columns: []string{"kind", "class_name", "entity_generated_time", "entity_version", "created_time",
		"executed_time", "nomad_template_name", "nomad_template_params", "unit_timeout_seconds", "status"},
	PartKey: []string{"class_name"},
	SortKey: []string{"created_time"},
}

type DPJob struct {
	ClassName           string            `json:"class_name"`
	CreatedTime         time.Time         `json:"created_time"`
	EntityGeneratedTime time.Time         `json:"entity_generated_time"`
	EntityVersion       string            `json:"entity_version"`
	Kind                string            `json:"kind"`
	ExecutedTime        time.Time         `json:"executed_time"`
	NomadTemplateName   string            `json:"nomad_template_name"`
	NomadTemplateParams map[string]string `json:"nomad_template_params"`
	UnitTimeoutSeconds  int               `json:"unit_timeout_seconds"`
	Status              string            `json:"status"`
}

var NomadTemplateMetadata = table.Metadata{
	Name:    "schedule.nomad_template",
	Columns: []string{"name", "git_url", "params", "params_description", "tags"},
	PartKey: []string{"name"},
	SortKey: nil,
}

type NomadTemplate struct {
	Name              string            `json:"name"`
	GitUrl            string            `json:"git_url"`
	Params            map[string]string `json:"params"`
	ParamsDescription map[string]string `json:"params_description"`
	Tags              []string          `json:"tags"`
}

func (j *DPJob) CreateOrUpdate() error {
	stmt, names := qb.Insert(DPJobMetadata.Name).
		Columns(DPJobMetadata.Columns...).ToCql()
	q := CqlSession.Query(stmt, names).BindStruct(*j)
	return q.ExecRelease()
}

func (nt *NomadTemplate) CreateOrUpdate() error {
	stmt, names := qb.Insert(NomadTemplateMetadata.Name).
		Columns(NomadTemplateMetadata.Columns...).ToCql()
	q := CqlSession.Query(stmt, names).BindStruct(*nt)
	return q.ExecRelease()
}

func (nt *NomadTemplate) GetByName() error {
	var nomadTemplateTable = table.New(NomadTemplateMetadata)
	stmt, names := nomadTemplateTable.Get()
	q := CqlSession.Query(stmt, names).BindStruct(*nt)
	return q.GetRelease(nt)
}

//TODO
func (j *DPJob) GetJobByEntityVersion(ev string) error {
	return nil
}
