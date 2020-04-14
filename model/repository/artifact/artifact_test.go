package artifact

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"log"
	"testing"
)

var session *gocql.Session

//初始化
func init() {
	cluster := gocql.NewCluster("127.0.0.1:9043")
	cluster.Keyspace = "artifacts"
	cluster.Consistency = gocql.Consistency(1)
	cluster.NumConns = 3
	var err error
	session, err = cluster.CreateSession()
	if err != nil {
		log.Panic("start", err)
		return
	}
}

//插入数据
func TestScylla_Insert(t *testing.T) {
	query := fmt.Sprintf(`INSERT INTO artifact_class (group, name, repo, profile, id, ts) VALUES ('group1','name1', 'repo1', 'profile1', now(), toUnixTimestamp(now()) )`)
	err := session.Query(query).Exec()
	if err != nil {
		fmt.Println(err)
	}
}

func TestGocqlxInsert(t *testing.T) {
	type A struct {
		col_1 string
		col_2 string
	}
	a := A{"my_name", "my_group"}
	stmt, names := qb.Insert("artifacts.table_b").Columns("col_1", "col_2").ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindStruct(a)
	if err := q.ExecRelease(); err != nil {
		t.Fatal(err)
	}
}
