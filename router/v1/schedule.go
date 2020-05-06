package router

import (
	"github.com/gin-gonic/gin"
	nomadHandler "voyageone.com/dp/scheduler/handler/nomad"
)

func InitScheduleApiGroup(apiRootRouter *gin.RouterGroup) {
	scheduleApiGroup := apiRootRouter.Group("schedule")
	{
		nomadApiGroup := scheduleApiGroup.Group("nomad")
		{
			nomadApiGroup.POST("deployer", nomadHandler.CreateOrUpdateDeployer)
			nomadApiGroup.GET("deployer")

			nomadApiGroup.POST("template", nomadHandler.CreateOrUpdateJobTemplate)
			nomadApiGroup.GET("template")

			nomadApiGroup.POST("job/immediate", nomadHandler.SubbmitJobImmediately)

			nomadApiGroup.GET("job/health/:jobId", nomadHandler.CheckJobLastDeploymentStatus)
		}
	}
}
