package service

import (
	"fmt"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"github.com/scylladb/gocqlx/table"
	"voyageone.com/dp/artifactKeeper/model/repository"
	"voyageone.com/dp/infrastructure/entity/global"
)

func CreateOrUpdateRepo(r repository.Repo) error {
	return nil
}

func CreateOrUpdateClass(c repository.Class) error {
	stmt, names := qb.Insert("artifact.class").
		Columns(repository.ClassMetadata.Columns...).
		ToCql()
	q := gocqlx.Query(global.CqlSession.Query(stmt), names).BindStruct(c)
	return q.ExecRelease()
}

func CreateOrUpdateEntity(e repository.Entity) error {
	stmt, names := qb.Insert("artifact.entity").
		Columns(repository.EntityMetadata.Columns...).
		ToCql()
	q := gocqlx.Query(global.CqlSession.Query(stmt), names).BindStruct(e)
	return q.ExecRelease()
}

func ReadClassByPrimaryKey(c *repository.Class) error {
	var classTable = table.New(repository.ClassMetadata)
	stmt, names := classTable.Get()
	q := gocqlx.Query(global.CqlSession.Query(stmt), names).BindStruct(*c)
	return q.GetRelease(c)
}

func ReadManyClassRecently(many int) ([]repository.Class, error) {
	return nil, nil
}

func ReadEntityByPrimaryKey(e *repository.Entity) error {
	var entityTable = table.New(repository.EntityMetadata)
	stmt, names := entityTable.Get()
	q := gocqlx.Query(global.CqlSession.Query(stmt), names).BindStruct(*e)
	return q.GetRelease(e)
}

func ReadEntityByVersionPartitionKey(e *repository.Entity) (err error) {
	err = global.CqlSession.Query(`SELECT url, generated_time FROM artifact.mv_entity_by_version 
          WHERE version=? and group=? and profile=? and class_name=?`, e.Version, e.Group, e.Profile, e.ClassName).
		Scan(&e.Url, &e.GeneratedTime)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func ReadManyEntityRecently(many int) {

}
