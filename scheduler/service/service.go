package service

import (
	"github.com/imdario/mergo"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"github.com/scylladb/gocqlx/table"
	artifactRepository "voyageone.com/dp/artifactKeeper/model/repository"
	"voyageone.com/dp/infrastructure/model/customType"
	. "voyageone.com/dp/infrastructure/model/global"
	schedulerRepository "voyageone.com/dp/scheduler/model/repository"
)

func CreateOrUpdateDeployer(d schedulerRepository.Deployer) error {
	stmt, names := qb.Insert(schedulerRepository.DeployerMetadata.Name).
		Columns(schedulerRepository.DeployerMetadata.Columns...).ToCql()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(d)
	return q.ExecRelease()
}

func CreateOrUpdateJob(j schedulerRepository.DPJob) error {
	stmt, names := qb.Insert(schedulerRepository.DPJobMetadata.Name).
		Columns(schedulerRepository.DPJobMetadata.Columns...).ToCql()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(j)
	return q.ExecRelease()
}

func CreateOrUpdataTemplate(nt schedulerRepository.NomadTemplate) error {
	stmt, names := qb.Insert(schedulerRepository.NomadTemplateMetadata.Name).
		Columns(schedulerRepository.NomadTemplateMetadata.Columns...).ToCql()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(nt)
	return q.ExecRelease()
}

func GetDeployerByName(d *schedulerRepository.Deployer) error {
	var deployerTable = table.New(schedulerRepository.DeployerMetadata)
	stmt, names := deployerTable.Get()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(*d)
	err := q.GetRelease(d)
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func GetNomadTemplateByName(nt *schedulerRepository.NomadTemplate) error {
	var nomadTemplateTable = table.New(schedulerRepository.NomadTemplateMetadata)
	stmt, names := nomadTemplateTable.Get()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(*nt)
	return q.GetRelease(nt)
}

func GetJobById(j *schedulerRepository.DPJob) error {
	err := CqlSession.Query(`SELECT id, group, profile, class_name, entity_generated_time, entity_version, created_time, deployer_name, nomad_template_name, nomad_template_params 
		FROM schedule.mv_job_by_id WHERE id = ?`, &j.Id).
		Scan(&j.Id, &j.Group, &j.Profile, &j.ClassName, &j.EntityGeneratedTime, &j.EntityVersion, &j.CreatedTime, &j.DeployerName, &j.NomadTemplateName, &j.NomadTemplateParams)
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func GetJobByEntityVersion(j *schedulerRepository.DPJob, ev string) error {
	return nil
}

func GetDefaultDeployerByGroup(dd *schedulerRepository.DefaultDeployer) error {
	var defaultDeployerTable = table.New(schedulerRepository.DefaultDeployerMetadata)
	stmt, names := defaultDeployerTable.Get()
	q := gocqlx.Query(CqlSession.Query(stmt), names).BindStruct(*dd)
	return q.GetRelease(dd)
}

func MergeClassDefaultParamsIntoJob(j *schedulerRepository.DPJob, c artifactRepository.Class) error {
	err := mergo.Merge(&j.NomadTemplateParams, c.DefaultTemplateParams)
	if err != nil {
		return customType.DPError("merge class defaults to job failed " + err.Error())
	}
	return nil
}

func MergeTemplateDefaultParamsIntoJob(j *schedulerRepository.DPJob) error {
	var nt = schedulerRepository.NomadTemplate{
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
