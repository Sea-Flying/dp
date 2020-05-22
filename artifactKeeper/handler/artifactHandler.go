package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"voyageone.com/dp/artifactKeeper/model/repository"
	"voyageone.com/dp/artifactKeeper/service"
	"voyageone.com/dp/infrastructure/model/customType"
	. "voyageone.com/dp/infrastructure/model/global"
	"voyageone.com/dp/infrastructure/model/response"
)

func CreataRepo(c *gin.Context) {
	var repo repository.Repo
	_ = c.ShouldBindJSON(&repo)
	err := service.CreateOrUpdateRepo(repo)
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
	err = service.CreateOrUpdateClass(class)
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
		Group:   entity.Group,
		Name:    entity.ClassName,
		Profile: entity.Profile,
	}
	_ = service.GetClassByPrimaryKey(&class)
	var repo = repository.Repo{
		Name: class.RepoName,
	}
	err := service.GetRepoByName(&repo)
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
	if entity.Url == "" {
		switch entity.ClassKind {
		case "jar":
			entity.Url = repo.BaseUrl + "/" + entity.Group + "-" + entity.Profile + "/" + entity.ClassName + "/" + entity.ClassName + "-" + entity.Version + ".jar"
		case "docker":
			entity.Url = repo.BaseUrl + "/" + entity.Group + "-" + entity.Profile + "/" + entity.ClassName + ":" + entity.Version
		}
		DPLogger.Printf("Entity Url not specify, use defaulf pattern generate it: %s \n", entity.Url)
	}
	entity.GeneratedTime = time.Now()
	err = service.CreateOrUpdateEntity(entity)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("Artifact Entity Create Failed，%v", err), c)
	} else {
		response.OkWithMessage("Arifact Entity Create Success", c)
	}
}

func GetClass(c *gin.Context) {
	var class repository.Class
	_ = c.ShouldBindJSON(&class)
	err := service.GetClassByPrimaryKey(&class)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("get artifact class failed: %#v", err), c)
	} else {
		response.OkDetailed(class, "", c)
	}
}

func fillDefaultsIntoClass(c *repository.Class) error {
	var kind = repository.Kind{
		Group: c.Group,
		Name:  c.Kind,
	}
	err := service.GetKindByPrimaryKey(&kind)
	if err != nil {
		return customType.DPError(fmt.Sprintf("fill class default template failed, invalid kind %#v", err))
	}
	c.CreatedTime = time.Now()
	c.DefaultNomadTemplate = kind.ProfileDefaultTemplate[c.Profile]
	var defaultRepo = repository.DefaultRepo{
		Group: c.Group,
	}
	if c.RepoName == "" {
		err = service.GetDefaultRepoByGroup(&defaultRepo)
		if err != nil {
			return customType.DPError(fmt.Sprintf("fill class default template failed when get default repo %#v", err))
		}
		c.RepoName = defaultRepo.ProfileRepo[c.Profile]
	}
	return nil
}
