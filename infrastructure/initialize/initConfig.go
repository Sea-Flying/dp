package initialize

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"voyageone.com/dp/infrastructure/entity/global"
	"voyageone.com/dp/infrastructure/utils"
)

func InitConfig() {
	log.Println(utils.GetExecPath())
	_ = cleanenv.ReadConfig("D:\\Develop\\go\\dp\\dp.yml", global.DPConfig)
}
