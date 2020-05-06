package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"voyageone.com/dp/artifactKeeper/model/repository"
	"voyageone.com/dp/artifactKeeper/service"
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
	err := service.CreateOrUpdateClass(class)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("Artifact Class Create Failed，%v", err), c)
	} else {
		response.OkWithMessage("Artifact Class Create Success", c)
	}
}

func CreateEntity(c *gin.Context) {
	var entity repository.Entity
	_ = c.ShouldBindJSON(&entity)
	var repo = repository.Repo{
		Name: entity.RepoName,
	}
	err := service.GetRepoByName(&repo)
	if err != nil || repo.BaseUrl == "" {
		DPLogger.Printf("Create or Update Entity Failed: repoName %s is not existed \n", entity.RepoName)
		DPLogger.Println(err)
		response.FailWithMessage(fmt.Sprintf("Artifact Entity Create Failed，%v", err), c)
		return
	}
	//如果entity的url为空，则生成为默认值
	if entity.Url == "" {
		switch entity.ClassKind {
		case "jar":
			entity.Url = repo.BaseUrl + "/" + entity.Group + "-" + entity.Profile + "/" + entity.ClassName + "-" + entity.Version + ".jar"
		case "docker":
			entity.Url = repo.BaseUrl + "/" + entity.Group + "-" + entity.Profile + "/" + entity.ClassName + ":" + entity.Version
		}
		DPLogger.Printf("Entity Url not specify, use defaulf pattern generate it: %s \n", entity.Url)
	}
	err = service.CreateOrUpdateEntity(entity)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("Artifact Entity Create Failed，%v", err), c)
	} else {
		response.OkWithMessage("Arifact Entity Create Success", c)
	}
}
