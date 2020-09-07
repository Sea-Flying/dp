package repository

import (
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
	"time"
	"voyageone.com/dp/infrastructure/model/customType"
	. "voyageone.com/dp/infrastructure/model/global"
)

var StatusHistoryMetadata = table.Metadata{
	Name:    "app.status_history",
	Columns: []string{"app_name", "status_changed_to", "time"},
	PartKey: []string{"app_name"},
	SortKey: []string{"time"},
}

type StatusHistory struct {
	AppName         string    `json:"app_name"`
	StatusChangedTo string    `json:"status_changed_to"`
	Time            time.Time `json:"time"`
}

var StoppedAppCacheMetadata = table.Metadata{
	Name:    "app.stopped_app_cache",
	Columns: []string{"app_name", "nomad_job_json"},
	PartKey: []string{"app_name"},
	SortKey: nil,
}

type StoppedAppCache struct {
	AppName      string `json:"app_name"`
	NomadJobJson string `json:"nomad_job_json"`
}

func (h *StatusHistory) CreateOrUpdate() error {
	stmt, names := qb.Insert(StatusHistoryMetadata.Name).
		Columns(StatusHistoryMetadata.Columns...).ToCql()
	q := CqlSession.Query(stmt, names).BindStruct(h)
	err := q.ExecRelease()
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func GetByAppNameOrderByTime(appName string, orderDirection qb.Order, pageSize int,
	pageNum int) (total int, ret []StatusHistory, err error) {
	if pageNum < 1 || pageSize < 1 {
		return 0, nil, customType.DPError("invalid pageNum or pageSize when query status_history")
	}
	var statusHistories []StatusHistory
	err = CqlSession.Query(qb.Select("app.status_history").
		Where(qb.Eq("app_name")).
		OrderBy("time", orderDirection).Limit(300).ToCql()).
		Bind(appName).Iter().Select(&statusHistories)
	if err != nil || len(statusHistories) == 0 {
		return
	}
	total = len(statusHistories)
	var left, right int
	if pageNum == 1 {
		left = 0
		right = pageSize - 1
	} else {
		left = pageSize*(pageNum-1) - 1
		right = pageSize*pageNum - 1
	}
	ret = statusHistories[left:right]
	return
}

func (h *StatusHistory) GetByPrimaryKey() error {
	var statusHistoryTable = table.New(StatusHistoryMetadata)
	stmt, names := statusHistoryTable.Get()
	q := CqlSession.Query(stmt, names).BindStruct(*h)
	err := q.GetRelease(h)
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func (c *StoppedAppCache) CreateOrUpdate() error {
	stmt, names := qb.Insert(StoppedAppCacheMetadata.Name).
		Columns(StoppedAppCacheMetadata.Columns...).ToCql()
	q := CqlSession.Query(stmt, names).BindStruct(*c)
	err := q.ExecRelease()
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func (c *StoppedAppCache) GetByPrimaryKey() error {
	var stoppedAppCacheTable = table.New(StoppedAppCacheMetadata)
	stmt, names := stoppedAppCacheTable.Get()
	q := CqlSession.Query(stmt, names).BindStruct(*c)
	err := q.GetRelease(c)
	if err != nil {
		DPLogger.Println(err)
	}
	return err
}

func GetStoppedAppsList() (stoppedApps []string, err error) {
	err = CqlSession.Query(qb.Select("app.stopped_app_cache").
		Columns("app_name").ToCql()).Iter().Select(&stoppedApps)
	return
}

func AddStatusHistory(appId string, status string) (err error) {
	h := StatusHistory{
		AppName:         appId,
		StatusChangedTo: status,
		Time:            time.Now(),
	}
	err = h.CreateOrUpdate()
	return err
}
