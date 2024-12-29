package dao

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
)

// AuditDao is the data access object for table audit.
type AuditDao struct {
	db      gdb.DB
	table   string       // table is the underlying table name of the DAO.
	columns AuditColumns // columns contains all the column names of Table for convenient usage.
}

// AuditColumns defines and stores column names for table audit.
type AuditColumns struct {
	Id        string // self increasing id
	Module    string // module
	Event     string // event name
	From      string // identity
	Content   string // content
	CreatedAt string // created_at
}

// auditColumns holds the columns for table audit.
var auditColumns = AuditColumns{
	Id:        "id",
	Module:    "module",
	Event:     "event",
	From:      "from",
	Content:   "content",
	CreatedAt: "created_at",
}

// NewAuditDao creates and returns a new DAO object for table data access.
func NewAuditDao(db gdb.DB, table string) *AuditDao {
	return &AuditDao{
		db:      db,
		table:   table,
		columns: auditColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AuditDao) DB() gdb.DB {
	return dao.db
}

// Table returns the table name of current dao.
func (dao *AuditDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AuditDao) Columns() AuditColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AuditDao) Group() string {
	return dao.db.GetGroup()
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AuditDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *AuditDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
