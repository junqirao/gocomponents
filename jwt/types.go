package jwt

type GenerateTokenRequest struct {
	// user basic info
	UserId   string `json:"user_id"`
	From     string `json:"from"`
	UserName string `json:"user_name"`
	// token info
	Subject   string      `json:"subject"`
	Issuer    string      `json:"issuer"`
	ExpiredIn int64       `json:"expired_in"` // seconds
	Key       interface{} `json:"key"`
}

type ParseTokenRequest struct {
	Token string      `json:"token"`
	Key   interface{} `json:"Key"`
}
