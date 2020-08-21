package initialize

import (
	"fmt"
	"testing"
	"voyageone.com/dp/infrastructure/model/global"
)

func TestInitNomadClient(t *testing.T) {
	InitConfig("E:\\Develop\\go\\dp\\dp.yml")
	c, err := initNomadClient(global.DPConfig.Nomad)
	if err != nil {
		t.Error(err)
	}
	node, _, _ := c.Nodes().Info("0a782816-f8c7-d0e3-159c-ac4324638951", nil)
	fmt.Println(node)
}
