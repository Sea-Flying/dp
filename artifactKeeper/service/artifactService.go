package service

import (
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"github.com/scylladb/gocqlx/table"
	"voyageone.com/dp/artifactKeeper/model/repository"
	. "voyageone.com/dp/infrastructure/entity/global"
)

func CreateOrUpdateRepo(r repository.Repo) error {
	stmt, names := qb.Insert(repository.RepoMetadata.Name).
		Columns(repository.RepoMetadata.Columns...).ToCql()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(r)
	err := q.ExecRelease()
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func CreateOrUpdateClass(c repository.Class) error {
	stmt, names := qb.Insert(repository.ClassMetadata.Name).
		Columns(repository.ClassMetadata.Columns...).ToCql()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(c)
	err := q.ExecRelease()
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func CreateOrUpdateEntity(e repository.Entity) error {
	stmt, names := qb.Insert(repository.EntityMetadata.Name).
		Columns(repository.EntityMetadata.Columns...).ToCql()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(e)
	err := q.ExecRelease()
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func GetRepoByName(r *repository.Repo) error {
	var repoTable = table.New(repository.RepoMetadata)
	stmt, names := repoTable.Get()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(*r)
	err := q.GetRelease(r)
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func GetClassByPrimaryKey(c *repository.Class) error {
	var classTable = table.New(repository.ClassMetadata)
	stmt, names := classTable.Get()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(*c)
	err := q.GetRelease(c)
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

//todo
func GetManyClassRecently(many int) ([]repository.Class, error) {
	return nil, nil
}

func GetEntityByPrimaryKey(e *repository.Entity) error {
	var entityTable = table.New(repository.EntityMetadata)
	stmt, names := entityTable.Get()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(*e)
	err := q.GetRelease(e)
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func GetEntityByVersionPartitionKey(e *repository.Entity) error {
	err := CqlSession.Query(`SELECT url, generated_time FROM artifact.mv_entity_by_version 
          WHERE version=? and group=? and profile=? and class_name=?`, e.Version, e.Group, e.Profile, e.ClassName).
		Scan(&e.Url, &e.GeneratedTime)
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

//todo
func ReadManyEntityRecently(many int) {

}
