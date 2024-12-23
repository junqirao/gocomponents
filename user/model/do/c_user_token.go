// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// CUserToken is the golang structure of table c_user_token for DAO operations like Where/Data.
type CUserToken struct {
	g.Meta    `orm:"table:c_user_token, do:true"`
	Id        interface{} //
	UserId    interface{} //
	ClientIp  interface{} //
	ExpiredAt interface{} //
}
