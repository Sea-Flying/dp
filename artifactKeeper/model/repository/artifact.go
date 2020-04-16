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
	Group       string
	Name        string
	Profile     string
	Kind        string
	RepoName    string    `json:"repo_name"`
	CreatedTime time.Time `json:"created_time"`
	GitUrl      string    `json:"git_url"`
	CiUrl       string    `json:"ci_url"`
	Description string
}

var EntityMetadata = table.Metadata{
	Name:    "entity",
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
	Name:    "repo",
	Columns: []string{"name", "kind", "url", "baseUrl", "description", "created_time"},
	PartKey: []string{"name"},
	SortKey: []string{"created_time"},
}

type Repo struct {
	Name         string
	Kind         string
	ArtifactKind string `json:"artifact_kind"`
	Url          string
	BaseUrl      string
	Descprition  string
	CreatedTime  string `json:"created_time"`
}
