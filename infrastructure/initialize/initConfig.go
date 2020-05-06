package initialize

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"voyageone.com/dp/infrastructure/model/global"
	"voyageone.com/dp/infrastructure/utils"
)

func InitConfig(configPath string) {
	log.Println(utils.GetExecPath())
	_ = cleanenv.ReadConfig(configPath, &global.DPConfig)
}
