package service

import (
	"github.com/imdario/mergo"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"github.com/scylladb/gocqlx/table"
	"voyageone.com/dp/infrastructure/model/customType"
	. "voyageone.com/dp/infrastructure/model/global"
	"voyageone.com/dp/scheduler/model/repository"
)

func CreateOrUpdateDeployer(d repository.Deployer) error {
	stmt, names := qb.Insert(repository.DeployerMetadata.Name).
		Columns(repository.DeployerMetadata.Columns...).ToCql()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(d)
	return q.ExecRelease()
}

func CreateOrUpdateJob(j repository.DPJob) error {
	stmt, names := qb.Insert(repository.DPJobMetadata.Name).
		Columns(repository.DPJobMetadata.Columns...).ToCql()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(j)
	return q.ExecRelease()
}

func CreateOrUpdataTemplate(nt repository.NomadTemplate) error {
	stmt, names := qb.Insert(repository.NomadTemplateMetadata.Name).
		Columns(repository.NomadTemplateMetadata.Columns...).ToCql()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(nt)
	return q.ExecRelease()
}

func GetDeployerByName(d *repository.Deployer) error {
	var deployerTable = table.New(repository.DeployerMetadata)
	stmt, names := deployerTable.Get()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(*d)
	err := q.GetRelease(d)
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func GetNomadTemplateByName(nt *repository.NomadTemplate) error {
	var nomadTemplateTable = table.New(repository.NomadTemplateMetadata)
	stmt, names := nomadTemplateTable.Get()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(*nt)
	return q.GetRelease(nt)
}

func GetJobById(j *repository.DPJob) error {
	err := CqlSession.Query(`SELECT id, group, profile, class_name, entity_generated_time, entity_version, created_time, deployer_name, nomad_template_name, nomad_template_params 
		FROM schedule.mv_job_by_id WHERE id = ?`, &j.Id).
		Scan(&j.Id, &j.Group, &j.Profile, &j.ClassName, &j.EntityGeneratedTime, &j.EntityVersion, &j.CreatedTime, &j.DeployerName, &j.NomadTemplateName, &j.NomadTemplateParams)
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func GetJobByEntityVersion(j *repository.DPJob, ev string) error {
	return nil
}

func MergeTemplateDefaultParamsIntoJob(j *repository.DPJob) error {
	var nt = repository.NomadTemplate{
		Name: j.NomadTemplateName,
	}
	err := GetNomadTemplateByName(&nt)
	if err != nil {
		return customType.DPError("merge template defaults to job failed " + err.Error())
	}
	err = mergo.Merge(&j.NomadTemplateParams, nt.Params)
	if err != nil {
		return customType.DPError("merge template defaults to job failed " + err.Error())
	}
	return nil
}
