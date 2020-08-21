package repository

import (
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
	"time"
	"voyageone.com/dp/infrastructure/model/customType"
	. "voyageone.com/dp/infrastructure/model/global"
)

var ClassMetadata = table.Metadata{
	Name: "artifact.class",
	Columns: []string{"name", "kind", "repo_name", "created_time", "git_url", "ci_url",
		"description", "default_nomad_template", "default_template_params", "unit_timeout_seconds"},
	PartKey: []string{"name"},
	SortKey: nil,
}

type Class struct {
	Name                  string            `json:"name"`
	Kind                  string            `json:"kind"`
	RepoName              string            `json:"repo_name"`
	CreatedTime           time.Time         `json:"created_time"`
	GitUrl                string            `json:"git_url"`
	CiUrl                 string            `json:"ci_url"`
	Description           string            `json:"description"`
	DefaultNomadTemplate  string            `json:"default_nomad_template"`
	DefaultTemplateParams map[string]string `json:"default_template_params"`
	UnitTimeoutSeconds    int               `json:"unit_timeout_seconds"`
}

var EntityMetadata = table.Metadata{
	Name: "artifact.entity",
	Columns: []string{"class_name", "generated_time", "version", "repo_name", "class_kind",
		"uploader", "url", "checksum", "valid"},
	PartKey: []string{"class_name"},
	SortKey: []string{"generated_time"},
}

type Entity struct {
	ClassName     string    `json:"class_name"`
	Version       string    `json:"version"`
	GeneratedTime time.Time `json:"generated_time"`
	RepoName      string    `json:"repo_name"`
	ClassKind     string    `json:"class_kind"`
	Uploader      string    `json:"uploader"`
	Url           string    `json:"url"`
	Checksum      string    `json:"checksum"`
	Valid         bool      `json:"valid"`
}

var RepoMetadata = table.Metadata{
	Name:    "artifact.repo",
	Columns: []string{"name", "kind", "web_url", "base_url", "description", "created_time"},
	PartKey: []string{"name"},
	SortKey: nil,
}

type Repo struct {
	Name         string    `json:"name"`
	Kind         string    `json:"kind"`
	ArtifactKind string    `json:"artifact_kind"`
	WebUrl       string    `json:"web_url"`
	BaseUrl      string    `json:"base_url"`
	Description  string    `json:"description"`
	CreatedTime  time.Time `json:"created_time"`
}

var KindMetadata = table.Metadata{
	Name:    "artifact.kind",
	Columns: []string{"name", "default_repo", "default_template", "default_unit_timeout_seconds", "created_time"},
	PartKey: []string{"name"},
	SortKey: nil,
}

type Kind struct {
	Name                      string    `json:"name"`
	DefaultRepo               string    `json:"default_repo"`
	DefaultTemplate           string    `json:"default_template"`
	DefaultUnitTimeoutSeconds int       `json:"default_unit_timeout_seconds"`
	CreatedTime               time.Time `json:"created_time"`
}

func (r *Repo) CreateOrUpdate() error {
	stmt, names := qb.Insert(RepoMetadata.Name).
		Columns(RepoMetadata.Columns...).ToCql()
	q := CqlSession.Query(stmt, names).BindStruct(*r)
	err := q.ExecRelease()
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func (r *Repo) GetByName() error {
	var repoTable = table.New(RepoMetadata)
	stmt, names := repoTable.Get()
	q := CqlSession.Query(stmt, names).BindStruct(*r)
	err := q.GetRelease(r)
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func (c *Class) CreateOrUpdate() error {
	stmt, names := qb.Insert(ClassMetadata.Name).
		Columns(ClassMetadata.Columns...).ToCql()
	q := CqlSession.Query(stmt, names).BindStruct(*c)
	err := q.ExecRelease()
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func (c *Class) GetByPrimaryKey() error {
	var classTable = table.New(ClassMetadata)
	stmt, names := classTable.Get()
	q := CqlSession.Query(stmt, names).BindStruct(*c)
	err := q.GetRelease(c)
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

//todo
func GetManyClassRecently(many int) ([]Class, error) {
	return nil, nil
}

func (e *Entity) CreateOrUpdate() error {
	stmt, names := qb.Insert(EntityMetadata.Name).
		Columns(EntityMetadata.Columns...).ToCql()
	q := CqlSession.Query(stmt, names).BindStruct(*e)
	err := q.ExecRelease()
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func (e *Entity) GetByPrimaryKey() error {
	var entityTable = table.New(EntityMetadata)
	stmt, names := entityTable.Get()
	q := CqlSession.Query(stmt, names).BindStruct(*e)
	err := q.GetRelease(e)
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func (e *Entity) GetByVersionPartitionKey() error {
	if e.ClassName == "" || e.Version == "" {
		return customType.DPError("not sufficient fields provided when execute GetEntityByVersionPartitionKey")
	}
	err := CqlSession.Query(qb.Select("artifact.mv_entity_by_version").
		Columns("url", "generated_time").
		Where(qb.Eq("version"), qb.Eq("class_name")).ToCql()).
		Bind(e.Version, e.ClassName).
		Scan(&e.Url, &e.GeneratedTime)
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

//todo
func ReadManyEntityRecently(many int) {

}

func (k *Kind) CreateOrUpdateKind() error {
	stmt, names := qb.Insert(KindMetadata.Name).
		Columns(KindMetadata.Columns...).ToCql()
	q := CqlSession.Query(stmt, names).BindStruct(*k)
	err := q.ExecRelease()
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func (k *Kind) GetByPrimaryKey() error {
	var classTable = table.New(KindMetadata)
	stmt, names := classTable.Get()
	q := CqlSession.Query(stmt, names).BindStruct(*k)
	err := q.GetRelease(k)
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}
