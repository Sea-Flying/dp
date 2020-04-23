package model

import "github.com/scylladb/gocqlx/table"

var NomadTemplateMetadata = table.Metadata{
	Name:    "schedule.nomad_template",
	Columns: []string{"name", "git_url", "group_tags", "profile_tags", "params"},
	PartKey: []string{"name"},
	SortKey: nil,
}

type NomadTemplate struct {
	Name        string            `json:"name"`
	GitUrl      string            `json:"git_url"`
	GroupTags   []string          `json:"group_tags"`
	ProfileTags []string          `json:"profile_tags"`
	Params      map[string]string `json:"params"`
}
