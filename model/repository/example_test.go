package repository

import (
	"github.com/scylladb/gocqlx"
	. "github.com/scylladb/gocqlx/gocqlxtest"
	"github.com/scylladb/gocqlx/qb"
	"testing"
)

func TestExample(t *testing.T) {
	session := CreateSession(t)
	defer session.Close()

	//	const personSchema = `
	//CREATE TABLE IF NOT EXISTS artifacts.person (
	//    first_name text,
	//    last_name text,
	//    email list<text>,
	//    PRIMARY KEY(first_name, last_name)
	//)`
	//
	//	if err := ExecStmt(session, personSchema); err != nil {
	//		t.Fatal("create table:", err)
	//	}

	// Person represents a row in person table.
	// Field names are converted to camel case by default, no need to add special tags.
	// If you want to disable a field add `db:"-"` tag, it will not be persisted.
	type A struct {
		ColA string
		ColB string
	}

	p := A{
		"Song3",
		"Flying",
	}

	// Insert, bind data from struct.
	{
		stmt, names := qb.Insert("artifacts.table_name").Columns("col_a", "col_b").ToCql()
		q := gocqlx.Query(session.Query(stmt), names).BindStruct(p)

		if err := q.ExecRelease(); err != nil {
			t.Fatal(err)
		}
	}

}
