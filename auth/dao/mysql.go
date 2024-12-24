package dao

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
)

// ApplicationAuthDao is the data access object for table srv_version.
type ApplicationAuthDao struct {
	db      gdb.DB
	table   string                 // table is the underlying table name of the DAO.
	columns ApplicationAuthColumns // columns contains all the column names of Table for convenient usage.
}

// ApplicationAuthColumns defines and stores column names for table srv_version.
type ApplicationAuthColumns struct {
	Id          string // self increasing id
	Name        string // display name
	Description string // description
	AppId       string // unique app id, for index
	AppKey      string // unique app key
	AppSecret   string // cost time in ms
	CreatedAt   string // created_at
}

// applicationAuthColumns holds the columns for table srv_version.
var applicationAuthColumns = ApplicationAuthColumns{
	Id:          "id",
	Name:        "name",
	Description: "description",
	AppId:       "app_id",
	AppKey:      "app_key",
	AppSecret:   "app_secret",
	CreatedAt:   "created_at",
}

// NewApplicationAuthDao creates and returns a new DAO object for table data access.
func NewApplicationAuthDao(db gdb.DB, table string) *ApplicationAuthDao {
	return &ApplicationAuthDao{
		db:      db,
		table:   table,
		columns: applicationAuthColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ApplicationAuthDao) DB() gdb.DB {
	return dao.db
}

// Table returns the table name of current dao.
func (dao *ApplicationAuthDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ApplicationAuthDao) Columns() ApplicationAuthColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ApplicationAuthDao) Group() string {
	return dao.db.GetGroup()
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ApplicationAuthDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ApplicationAuthDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
