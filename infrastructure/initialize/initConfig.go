package initialize

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"voyageone.com/dp/infrastructure/entity/config"
	"voyageone.com/dp/infrastructure/utils"
)

func InitConfig(dpConfig *config.DPConfig) {
	log.Println(utils.GetExecPath())
	_ = cleanenv.ReadConfig("dp.yml", &dpConfig)
}
