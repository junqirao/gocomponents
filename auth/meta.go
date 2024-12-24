package auth

const (
	HeaderKeyAppId  = "AppId"
	HeaderKeyAppKey = "AppKey"
)

type Required struct {
	Authorization string `json:"Authorization" in:"header" v:"required"`
	AppId         string `json:"AppId" in:"header" dc:"calculated by middleware"`
	AppKey        string `json:"AppKey" in:"header" dc:"calculated by middleware"`
}
