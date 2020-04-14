package handler

import (
	"github.com/gin-gonic/gin"
	"voyageone.com/dp/artifactKeeper/service"
	"voyageone.com/dp/model/repository/artifact"
)

func CreateArtifactClass(c *gin.Context) {
	var artifactClass artifact.Class
	_ = c.ShouldBindJSON(&artifactClass)
	err := service.CreateArtifactClass(artifactClass)
	if err != nil {
		println(err)
	}
}
