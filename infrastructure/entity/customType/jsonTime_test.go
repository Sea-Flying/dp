package customType

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestTimeUnmarshal(t *testing.T) {
	const myTypeStr = `{"id":5, "name":"sqq", "time": "2020-04-16 04:26:22.999"}`
	type myType struct {
		Id   int
		Name string
		Time DpJsonTime
	}
	bytestr := []byte(myTypeStr)

	var i = new(myType)

	err := json.Unmarshal(bytestr, &i)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(i.Time)
	}
}

func TestTimeParse(t *testing.T) {
	//str := "2006-01-02T15:04:05Z07:00"
	str := "2020-04-16 04:26:22.029"
	strbyte := []byte(str)
	tt, _ := time.Parse(DpTimeLayout, string(strbyte))
	fmt.Println(tt)
}
