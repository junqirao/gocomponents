// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CUserDao is the data access object for the table c_user.
type CUserDao struct {
	table   string       // table is the underlying table name of the DAO.
	group   string       // group is the database configuration group name of the current DAO.
	columns CUserColumns // columns contains all the column names of Table for convenient usage.
}

// CUserColumns defines and stores column names for the table c_user.
type CUserColumns struct {
	Id            string //
	Username      string //
	Password      string //
	CreatedAt     string //
	UpdatedAt     string //
	Administrator string //
	Source        string //
	Status        string //
	Extra         string //
}

// cUserColumns holds the columns for the table c_user.
var cUserColumns = CUserColumns{
	Id:            "id",
	Username:      "username",
	Password:      "password",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	Administrator: "administrator",
	Source:        "source",
	Status:        "status",
	Extra:         "extra",
}

// NewCUserDao creates and returns a new DAO object for table data access.
func NewCUserDao() *CUserDao {
	return &CUserDao{
		group:   "default",
		table:   "c_user",
		columns: cUserColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CUserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CUserDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CUserDao) Columns() CUserColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CUserDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CUserDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *CUserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}