// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CUser is the golang structure for table c_user.
type CUser struct {
	Id            string      `json:"id"            orm:"id"            description:""` //
	Username      string      `json:"username"      orm:"username"      description:""` //
	Password      string      `json:"password"      orm:"password"      description:""` //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"    description:""` //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"    description:""` //
	Administrator int         `json:"administrator" orm:"administrator" description:""` //
	Source        string      `json:"source"        orm:"source"        description:""` //
	Status        int         `json:"status"        orm:"status"        description:""` //
	Extra         string      `json:"extra"         orm:"extra"         description:""` //
}
