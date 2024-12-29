package audit

import (
	"context"
	"errors"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"

	// db driver
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	"github.com/junqirao/gocomponents/audit/dao"
	"github.com/junqirao/gocomponents/audit/embed"
)

type auditDao = *dao.AuditDao

type MysqlAdaptor struct {
	logger *glog.Logger
	table  string
	auditDao
}

func NewMysqlAdaptor(db gdb.DB, logger *glog.Logger, tableName ...string) (m *MysqlAdaptor, err error) {
	table := "audit_record"
	if len(tableName) > 0 && tableName[0] != "" {
		table = tableName[0]
	}

	m = &MysqlAdaptor{
		logger:   logger,
		table:    table,
		auditDao: dao.NewAuditDao(db, table),
	}
	err = m.createTableIfNotExists(context.Background())
	return
}

func (m *MysqlAdaptor) Store(ctx context.Context, record *Record) (err error) {
	_, err = m.Ctx(ctx).Insert(g.Map{
		m.Columns().Module:    record.Module,
		m.Columns().Event:     record.Event,
		m.Columns().From:      record.From,
		m.Columns().Content:   gconv.String(record.Content),
		m.Columns().CreatedAt: record.CreatedAt,
	})
	m.logger.Cat("audit_record").Infof(ctx, "%s.%s(from=%s): %v", record.Module, record.Event, record.From, gconv.String(record.Content))
	return
}

func (m *MysqlAdaptor) Load(ctx context.Context, params *RecordQueryParams) (res *RecordQueryResult, err error) {
	res = &RecordQueryResult{
		Records: make([]*Record, 0),
		Total:   0,
	}
	where := g.Map{
		m.Columns().Module: params.Module,
		m.Columns().Event:  params.Event,
		m.Columns().From:   params.From,
	}
	if params.Start != nil {
		where[m.Columns().CreatedAt+" >"] = *params.Start
	}
	if params.End != nil {
		where[m.Columns().CreatedAt+" <"] = *params.End
	}
	query := m.Ctx(ctx).Where(where).OmitEmptyWhere()
	res.Total, err = query.Count()
	if err != nil {
		return
	}
	err = query.Page(params.PageNum, params.PageSize).OrderDesc(m.Columns().CreatedAt).Scan(&res.Records)
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
