package model

type (
	CheckUsernameReq struct {
		Username string `json:"username" v:"username@required|length:3,20"`
	}

	CreateUserReq struct {
		Username string `json:"username" v:"username@required|length:3,20"`
		Password string `json:"password" v:"password@required|length:6,50"`
	}
)
