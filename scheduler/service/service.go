package service

import (
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"github.com/scylladb/gocqlx/table"
	. "voyageone.com/dp/infrastructure/entity/global"
	"voyageone.com/dp/scheduler/model"
)

func CreateOrUpdateDeployer(d model.Deployer) error {
	stmt, names := qb.Insert(model.DeployerMetadata.Name).
		Columns(model.DeployerMetadata.Columns...).ToCql()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(d)
	return q.ExecRelease()
}

func CreateOrUpdateJob(j model.DPJob) error {
	stmt, names := qb.Insert(model.DPJobMetadata.Name).
		Columns(model.DPJobMetadata.Columns...).ToCql()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(j)
	return q.ExecRelease()
}

func CreateOrUpdataTemplate(nt model.NomadTemplate) error {
	stmt, names := qb.Insert(model.NomadTemplateMetadata.Name).
		Columns(model.NomadTemplateMetadata.Columns...).ToCql()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(nt)
	return q.ExecRelease()
}

func GetDeployerByName(d *model.Deployer) error {
	var deployerTable = table.New(model.DeployerMetadata)
	stmt, names := deployerTable.Get()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(*d)
	err := q.GetRelease(d)
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func GetNomadTemplateByName(nt *model.NomadTemplate) error {
	var nomadTemplateTable = table.New(model.NomadTemplateMetadata)
	stmt, names := nomadTemplateTable.Get()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(*nt)
	return q.GetRelease(nt)
}

func GetJobById(j *model.DPJob) error {
	err := CqlSession.Query(`SELECT id, group, profile, class_name, entity_generated_time, entity_version, created_time, deployer_name, nomad_template_name, nomad_template_params 
		FROM schedule.mv_job_by_id WHERE id = ?`, &j.Id).
		Scan(&j.Id, &j.Group, &j.Profile, &j.ClassName, &j.EntityGeneratedTime, &j.EntityVersion, &j.CreatedTime, &j.DeployerName, &j.NomadTemplateName, &j.NomadTemplateParams)
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func GetJobByEntityVersion(j *model.DPJob, ev string) error {
	return nil
}
