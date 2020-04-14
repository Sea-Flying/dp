package init

import (
	"github.com/ilyakaznacheev/cleanenv"
	"voyageone.com/dp/model"
)

func InitConfig(dpConfig *model.DPConfig) {
	err := cleanenv.ReadConfig("dp.yml", &dpConfig)
	if err != nil {
		panic("dp load config failed!")
	}
}
