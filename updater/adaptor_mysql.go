package updater

import (
	"context"
	"errors"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"

	"github.com/junqirao/gocomponents/updater/dao"
	"github.com/junqirao/gocomponents/updater/embed"
)

type svDao = *dao.SrvVersionDao

type MysqlAdaptor struct {
	logger *glog.Logger
	table  string
	svDao
}

func NewMysqlAdaptor(ctx context.Context, db gdb.DB, logger *glog.Logger, tableName ...string) (m *MysqlAdaptor, err error) {
	table := "srv_version"
	if len(tableName) > 0 && tableName[0] != "" {
		table = tableName[0]
	}

	m = &MysqlAdaptor{
		logger: logger,
		table:  table,
		svDao:  dao.NewSrvVersionDao(db, table),
	}
	err = m.createTableIfNotExists(ctx)
	return
}

func (m *MysqlAdaptor) Store(ctx context.Context, record *Record) (err error) {
	where := g.Map{
		m.Columns().Name: record.Name,
		m.Columns().Type: record.Type,
	}
	count, err := m.Ctx(ctx).Where(where).Count()
	if err != nil {
		return
	}
	if count > 0 {
		_, err = m.Ctx(ctx).Where(where).Update(g.Map{
			m.Columns().Status:    record.Status,
			m.Columns().Cost:      record.Cost,
			m.Columns().CreatedAt: record.CreatedAt,
		})
	} else {
		_, err = m.Ctx(ctx).Insert(g.Map{
			m.Columns().Name:      record.Name,
			m.Columns().Type:      record.Type,
			m.Columns().Status:    record.Status,
			m.Columns().Cost:      record.Cost,
			m.Columns().CreatedAt: record.CreatedAt,
		})
	}
	return
}

func (m *MysqlAdaptor) Load(ctx context.Context, params *RecordQueryParams) (res *RecordQueryResult, err error) {
	res = &RecordQueryResult{
		Records: make([]*Record, 0),
		Total:   0,
	}
	where := g.Map{}
	if params.Name != nil {
		where[m.Columns().Name] = *params.Name
	}
	if params.Type != nil {
		where[m.Columns().Type] = *params.Type
	}
	query := m.Ctx(ctx).Where(where).OmitEmptyWhere()
	res.Total, err = query.Count()
	if err != nil {
		return
	}
	err = query.OrderDesc(m.Columns().CreatedAt).Scan(&res.Records)
	return
}

func (m *MysqlAdaptor) createTableIfNotExists(ctx context.Context) (err error) {
	tables, err := m.DB().Tables(ctx)
	if err != nil {
		return
	}
	for _, table := range tables {
		if table == m.table {
			return
		}
	}
	initSql := embed.GetInitSqlContent()
	if initSql == "" {
		err = errors.New("init sql content is empty")
		return
	}
	_, err = m.DB().Exec(ctx, initSql)
	return
}
