package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"voyageone.com/dp/artifactKeeper/model/repository"
	"voyageone.com/dp/artifactKeeper/service"
	"voyageone.com/dp/infrastructure/response"
)

func CreateClass(c *gin.Context) {
	var class repository.Class
	_ = c.ShouldBindJSON(&class)
	err := service.CreateOrUpdateClass(class)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("Artifact Class Create Failed，%v", err), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

func CreateEntity(c *gin.Context) {
	var entity repository.Entity
	_ = c.ShouldBindJSON(&entity)
	err := service.CreateOrUpdateEntity(entity)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("Artifact Entity Create Failed，%v", err), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}
