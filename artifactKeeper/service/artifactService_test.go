package service

import (
	"bufio"
	"fmt"
	"github.com/gocql/gocql"
	"io"
	"log"
	"os"
	"testing"
	"time"
	"voyageone.com/dp/artifactKeeper/model/repository"
	"voyageone.com/dp/infrastructure/model/global"
)

func init() {
	cluster := gocql.NewCluster("127.0.0.1:9043")
	cluster.Consistency = gocql.Quorum
	cluster.NumConns = 3
	var err error
	global.CqlSession, err = cluster.CreateSession()
	if err != nil {
		log.Fatal("CQL Session Create Failed!", err)
	}
	global.DPLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile)
}

func TestCreateOrUpdateRepo(t *testing.T) {
	var r = repository.Repo{
		Name:         "local_http",
		Kind:         "http",
		ArtifactKind: "jar",
		WebUrl:       "http://10.0.0.70:4567",
		BaseUrl:      "http://10.0.0.70:4567",
		Description:  "local jar package http repo",
		CreatedTime:  time.Now(),
	}
	err := CreateOrUpdateRepo(r)
	if err != nil {
		t.Error(err)
	}
}

func TestGetRepoByName(t *testing.T) {
	var r = repository.Repo{
		Name: "local_http",
	}
	err := GetRepoByName(&r)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(r)
	}
}

func TestCreateArtifactClass(t *testing.T) {
	var c = repository.Class{
		Name:        "voerp-restapi",
		Group:       "voerp",
		Profile:     "prod",
		Kind:        "jar",
		RepoName:    "hk_oss",
		CreatedTime: time.Now(),
		GitUrl:      "https://github.com",
		CiUrl:       "",
		Description: "",
	}
	err := CreateOrUpdateClass(c)
	t.Log(err)
}

func TestGetClassByPrimaryKey(t *testing.T) {
	var c = repository.Class{
		Group:   "voerp",
		Profile: "staging",
		Name:    "voerp-product",
	}
	err := GetClassByPrimaryKey(&c)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(c)
	}
}

func TestCreateOrUpdateEntity(t *testing.T) {
	var e = repository.Entity{
		Group:         "voerp",
		Profile:       "staging",
		ClassName:     "openvms-restapi-dp",
		ClassKind:     "jar",
		GeneratedTime: time.Now(),
		RepoName:      "local_http",
		Uploader:      "jenkins",
		Version:       "v202003241445",
		Url:           "http://10.0.0.70:4567/voerp-staging/openvms-restapi/openvms-restapi-v202003241445.jar",
	}
	err := CreateOrUpdateEntity(e)
	if err != nil {
		t.Error(err)
	}
}

func TestGetEntityByVersionPartitionKey(t *testing.T) {
	var e = repository.Entity{
		Group:     "voerp",
		Profile:   "dev",
		ClassName: "voerp-product",
		Version:   "v202004161644",
	}
	err := GetEntityByVersionPartitionKey(&e)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(e)
	}
}

func TestOsOpen(t *testing.T) {
	f, err := os.OpenFile("D:\\Develop\\go\\dp\\dp.yml", os.O_RDONLY|os.O_SYNC, 0)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		con, err := r.ReadString('\n') //读取一行
		fmt.Print(con)
		//如何判断文件读取完毕
		if err == io.EOF {
			fmt.Println("文件读取结束")
			break
		}
	}
}
