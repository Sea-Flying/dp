package repository

import (
	"github.com/scylladb/gocqlx/table"
	"time"
)

var ClassMetadata = table.Metadata{
	Name: "artifact.class",
	Columns: []string{"group", "profile", "name", "kind", "repo_name", "created_time", "git_url", "ci_url",
		"description", "default_nomad_template", "default_template_params"},
	PartKey: []string{"group", "profile", "name"},
	SortKey: nil,
}

type Class struct {
	Group                 string            `json:"group"`
	Profile               string            `json:"profile"`
	Name                  string            `json:"name"`
	Kind                  string            `json:"kind"`
	RepoName              string            `json:"repo_name"`
	CreatedTime           time.Time         `json:"created_time"`
	GitUrl                string            `json:"git_url"`
	CiUrl                 string            `json:"ci_url"`
	Description           string            `json:"description"`
	DefaultNomadTemplate  string            `json:"default_nomad_template"`
	DefaultTemplateParams map[string]string `json:"default_template_params"`
}

var EntityMetadata = table.Metadata{
	Name: "artifact.entity",
	Columns: []string{"group", "profile", "class_name", "generated_time", "version", "repo_name", "class_kind",
		"uploader", "url", "checksum"},
	PartKey: []string{"group", "profile", "class_name"},
	SortKey: []string{"generated_time"},
}

type Entity struct {
	Group         string    `json:"group"`
	ClassName     string    `json:"class_name"`
	Version       string    `json:"version"`
	GeneratedTime time.Time `json:"generated_time"`
	Profile       string    `json:"profile"`
	RepoName      string    `json:"repo_name"`
	ClassKind     string    `json:"class_kind"`
	Uploader      string    `json:"uploader"`
	Url           string    `json:"url"`
	Checksum      string    `json:"checksum"`
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

var DefaultRepoMetadata = table.Metadata{
	Name:    "artifact.default_repo",
	Columns: []string{"group", "profile_repo"},
	PartKey: []string{"group"},
	SortKey: nil,
}

type DefaultRepo struct {
	Group       string            `json:"group"`
	ProfileRepo map[string]string `json:"profile_repo"`
}

var KindMetadata = table.Metadata{
	Name:    "artifact.kind",
	Columns: []string{"group", "name", "profile_default_template", "created_time"},
	PartKey: []string{"group", "name"},
	SortKey: nil,
}

type Kind struct {
	Group                  string            `json:"group"`
	Name                   string            `json:"name"`
	ProfileDefaultTemplate map[string]string `json:"profile_default_template"`
	CreatedTime            time.Time         `json:"created_time"`
}
