package initialize

import (
	"log"
	"os"
	. "voyageone.com/dp/infrastructure/model/global"
)

func InitLogger() {
	DPLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile)
}
