// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// CUserToken is the golang structure for table c_user_token.
type CUserToken struct {
	Id        string `json:"id"        orm:"id"         description:""` //
	UserId    string `json:"userId"    orm:"user_id"    description:""` //
	ClientIp  string `json:"clientIp"  orm:"client_ip"  description:""` //
	ExpiredAt int    `json:"expiredAt" orm:"expired_at" description:""` //
}
