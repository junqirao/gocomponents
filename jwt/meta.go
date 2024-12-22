package jwt

const (
	HeaderKeyUserId   = "X-UserId"
	HeaderKeyUserName = "X-UserName"
	HeaderKeyUserFrom = "X-UserFrom"
)

type Required struct {
	Authorization string `json:"Authorization" in:"header" v:"required"`
	UserName      string `json:"X-UserName" in:"header" dc:"calculated by middleware"`
	UserId        string `json:"X-UserId" in:"header" dc:"calculated by middleware"`
	UserFrom      string `json:"X-UserFrom" in:"header" dc:"calculated by middleware"`
}
