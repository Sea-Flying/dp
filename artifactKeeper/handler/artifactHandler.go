package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
	"voyageone.com/dp/artifactKeeper/model/repository"
	"voyageone.com/dp/infrastructure/model/customType"
	. "voyageone.com/dp/infrastructure/model/global"
	"voyageone.com/dp/infrastructure/model/response"
)

func CreateRepo(c *gin.Context) {
	var repo repository.Repo
	_ = c.ShouldBindJSON(&repo)
	err := repo.CreateOrUpdate()
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("Artifact Repo Create Failed, %v", err), c)
	} else {
		response.OkWithMessage("Artifact Repo Create Success", c)
	}
}

func CreateClass(c *gin.Context) {
	var class repository.Class
	_ = c.ShouldBindJSON(&class)
	err := fillDefaultsIntoClass(&class)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("artifact class create failed，%v", err), c)
		return
	}
	err = class.CreateOrUpdate()
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("Artifact Class Create Failed，%v", err), c)
	} else {
		response.OkWithMessage("Artifact Class Create Success", c)
	}
}

func CreateEntity(c *gin.Context) {
	var entity repository.Entity
	_ = c.ShouldBindJSON(&entity)
	var class = repository.Class{
		Name: entity.ClassName,
	}
	_ = class.GetByPrimaryKey()
	var repo = repository.Repo{
		Name: class.RepoName,
	}
	err := repo.GetByName()
	if err != nil || repo.BaseUrl == "" {
		DPLogger.Printf("Create or Update Entity Failed: repoName %s is not existed \n", entity.RepoName)
		DPLogger.Println(err)
		response.FailWithMessage(fmt.Sprintf("Artifact Entity Create Failed，%v", err), c)
		return
	} else {
		entity.RepoName = class.RepoName
		entity.ClassKind = class.Kind
	}
	//如果entity的url为空，则生成为默认值
	//TODO 若之后class kind种类多了，或者需要灵活规则，考虑将下面这个if自动生成步骤放在client端或者server端的一个规则配置
	if entity.Url == "" {
		switch {
		case entity.ClassKind == "jar":
			entity.Url = repo.BaseUrl + "/" + entity.ClassName + "/" + entity.ClassName + "-" + entity.Version + ".jar"
		case strings.HasPrefix(entity.ClassKind, "docker"):
			entity.Url = repo.BaseUrl + "/" + entity.ClassName + ":" + entity.Version
		}
		DPLogger.Printf("Entity Url not specify, use defaulf pattern generate it: %s \n", entity.Url)
	}
	entity.GeneratedTime = time.Now()
	err = entity.CreateOrUpdate()
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("Artifact Entity Create Failed，%v", err), c)
	} else {
		response.OkWithMessage("Artifact Entity Create Success", c)
	}
}

func GetClass(c *gin.Context) {
	var class = repository.Class{
		Name: c.Param("className"),
	}
	err := class.GetByPrimaryKey()
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("get artifact class failed: %#v", err), c)
	} else {
		response.OkDetailed(class, "", c)
	}
}

func fillDefaultsIntoClass(c *repository.Class) error {
	var kind = repository.Kind{
		Name: c.Kind,
	}
	err := kind.GetByPrimaryKey()
	if err != nil {
		return customType.DPError(fmt.Sprintf("fill class default template failed, invalid kind %#v", err))
	}
	c.CreatedTime = time.Now()
	//workaround: 如果项目名称以"-task"结尾，则nomad的job模板文件使用以"-task"结尾的名称的模板
	if strings.HasSuffix(c.Name, "-task") {
		c.DefaultNomadTemplate = kind.DefaultTemplate + "-task"
	} else {
		c.DefaultNomadTemplate = kind.DefaultTemplate
	}
	c.RepoName = kind.DefaultRepo
	c.UnitTimeoutSeconds = kind.DefaultUnitTimeoutSeconds
	return nil
}
