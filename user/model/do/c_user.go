// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CUser is the golang structure of table c_user for DAO operations like Where/Data.
type CUser struct {
	g.Meta        `orm:"table:c_user, do:true"`
	Id            interface{} //
	Username      interface{} //
	Password      interface{} //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
	Administrator interface{} //
	Source        interface{} //
	Status        interface{} //
	Extra         interface{} //
}
