package service

import (
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"voyageone.com/dp/global"
	"voyageone.com/dp/model/repository/artifact"
)

func CreateArtifactClass(j artifact.Class) error {
	stmt, names := qb.Insert("artifact.class").Columns("").ToCql()
	q := gocqlx.Query(global.CqlSession.Query(stmt), names).BindStruct(j)
	return q.ExecRelease()
}
