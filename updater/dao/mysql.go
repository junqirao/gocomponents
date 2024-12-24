package dao

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
)

// SrvVersionDao is the data access object for table srv_version.
type SrvVersionDao struct {
	db      gdb.DB
	table   string            // table is the underlying table name of the DAO.
	columns SrvVersionColumns // columns contains all the column names of Table for convenient usage.
}

// SrvVersionColumns defines and stores column names for table srv_version.
type SrvVersionColumns struct {
	Id        string // self increasing id
	Name      string // unique name
	Type      string // 0: function, 1: sql
	Status    string // 0: failed, 1: success
	Cost      string // cost time in ms
	CreatedAt string // created_at
}

// srvVersionColumns holds the columns for table srv_version.
var srvVersionColumns = SrvVersionColumns{
	Id:        "id",
	Name:      "name",
	Type:      "type",
	Status:    "status",
	Cost:      "cost",
	CreatedAt: "created_at",
}

// NewSrvVersionDao creates and returns a new DAO object for table data access.
func NewSrvVersionDao(db gdb.DB, table string) *SrvVersionDao {
	return &SrvVersionDao{
		db:      db,
		table:   table,
		columns: srvVersionColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SrvVersionDao) DB() gdb.DB {
	return dao.db
}

// Table returns the table name of current dao.
func (dao *SrvVersionDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SrvVersionDao) Columns() SrvVersionColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SrvVersionDao) Group() string {
	return dao.db.GetGroup()
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SrvVersionDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SrvVersionDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
