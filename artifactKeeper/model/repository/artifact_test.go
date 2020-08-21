package repository

import (
	"bufio"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"io"
	"log"
	"os"
	"testing"
	"time"
	"voyageone.com/dp/infrastructure/model/global"
)

//初始化

func init() {
	cluster := gocql.NewCluster("127.0.0.1:9042")
	cluster.Consistency = gocql.Quorum
	cluster.NumConns = 3
	var err error
	global.CqlSession, err = gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Fatal("CQL Session Create Failed!", err)
	}
	global.DPLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile)
}

func TestCreateOrUpdateRepo(t *testing.T) {
	var r = Repo{
		Name:         "test2_http",
		Kind:         "http",
		ArtifactKind: "jar",
		WebUrl:       "http://10.0.0.70:4567",
		BaseUrl:      "http://10.0.0.70:4567",
		Description:  "local jar package http repo",
		CreatedTime:  time.Now(),
	}
	err := r.CreateOrUpdate()
	if err != nil {
		t.Error(err)
	}
}

func TestGetRepoByName(t *testing.T) {
	var r = Repo{
		Name: "test2_http",
	}
	err := r.GetByName()
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(r)
	}
}

func TestCreateArtifactClass(t *testing.T) {
	var c = Class{
		Name:        "voerp-restapi",
		Kind:        "jar",
		RepoName:    "hk_oss",
		CreatedTime: time.Now(),
		GitUrl:      "https://github.com",
		CiUrl:       "",
		Description: "",
	}
	err := c.CreateOrUpdate()
	t.Log(err)
}

func TestGetClassByPrimaryKey(t *testing.T) {
	var c = Class{
		Name: "voerp-product",
	}
	err := c.GetByPrimaryKey()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(c)
	}
}

func TestCreateOrUpdateEntity(t *testing.T) {
	var e = Entity{
		ClassName:     "sf-class",
		ClassKind:     "jar",
		GeneratedTime: time.Now(),
		RepoName:      "local_http",
		Uploader:      "jenkins",
		Version:       "v202007301956",
		Url:           "http://10.0.0.70:4567/voerp-staging/openvms-restapi/openvms-restapi-v202003241445.jar",
	}
	err := e.CreateOrUpdate()
	if err != nil {
		t.Error(err)
	}
}

func TestGetEntityByVersionPartitionKey(t *testing.T) {
	var e = Entity{
		ClassName: "sf-class",
		Version:   "v202007301956",
	}
	err := e.GetByVersionPartitionKey()
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(e)
	}
}

func TestOsOpen(t *testing.T) {
	f, err := os.OpenFile("E:\\Develop\\go\\dp\\dp.yml", os.O_RDONLY|os.O_SYNC, 0)
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
