package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/junqirao/gocomponents/response"
)

type wrappedErrors interface {
	Unwrap() []error
}

type Claims struct {
	jwt.RegisteredClaims

	UserId   string `json:"uid,omitempty"`
	UserName string `json:"nam,omitempty"`
	From     string `json:"fro,omitempty"`
}

// GenerateToken ...
func GenerateToken(req *GenerateTokenRequest) (string, error) {
	ts := time.Now()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  []string{req.UserId},
			ExpiresAt: jwt.NewNumericDate(ts.Add(time.Second * time.Duration(req.ExpiredIn))), // 过期时间
			IssuedAt:  jwt.NewNumericDate(ts),
			Issuer:    req.Issuer, // 签发人
			NotBefore: jwt.NewNumericDate(ts),
			Subject:   req.Subject,
		},
		UserId:   req.UserId,
		UserName: req.UserName,
		From:     req.From,
	}).SignedString(req.Key)
}

// ParseToken ...
func ParseToken(req *ParseTokenRequest) (claims *Claims, err error) {
	claims = new(Claims)
	token, err := jwt.ParseWithClaims(req.Token, claims, func(_ *jwt.Token) (interface{}, error) {
		return req.Key, nil
	})
	if err != nil {
		return
	}
	if token.Valid {
		return claims, nil
	}

	for _, e := range err.(wrappedErrors).Unwrap() {
		if e != nil {
			return nil, response.CodeUnauthorized.WithMessage(e.Error())
		}
	}
	return
}
