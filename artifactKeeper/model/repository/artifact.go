package repository

import (
	"github.com/scylladb/gocqlx/table"
	"time"
)

var ClassMetadata = table.Metadata{
	Name:    "artifact.class",
	Columns: []string{"group", "profile", "name", "kind", "repo_name", "created_time", "git_url", "ci_url"},
	PartKey: []string{"group", "profile", "name"},
	SortKey: []string{},
}

type Class struct {
	Group       string    `json:"group"`
	Name        string    `json:"name"`
	Profile     string    `json:"profile"`
	Kind        string    `json:"kind"`
	RepoName    string    `json:"repo_name"`
	CreatedTime time.Time `json:"created_time"`
	GitUrl      string    `json:"git_url"`
	CiUrl       string    `json:"ci_url"`
	Description string    `json:"description"`
}

var EntityMetadata = table.Metadata{
	Name:    "artifact.model",
	Columns: []string{"group", "profile", "class_name", "generated_time", "version", "repo_name", "class_kind", "uploader", "url", "checksum"},
	PartKey: []string{"group", "profile", "class_name"},
	SortKey: []string{"generated_time"},
}

type Entity struct {
	Group         string
	ClassName     string
	Version       string
	GeneratedTime time.Time `json:"generated_time"`
	Profile       string
	RepoName      string `json:"repo_name"`
	ClassKind     string `json:"class_kind"`
	Uploader      string
	Url           string
	Checksum      string
}

var RepoMetadata = table.Metadata{
	Name:    "artifact.repo",
	Columns: []string{"name", "kind", "web_url", "base_url", "description", "created_time"},
	PartKey: []string{"name"},
	SortKey: []string{},
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
