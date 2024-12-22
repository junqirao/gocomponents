package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	testKey = []byte("test_key")
)

func TestGenToken(t *testing.T) {
	token, err := GenerateToken(&GenerateTokenRequest{
		UserId:    "test_user_id",
		Subject:   "subject",
		Issuer:    "issuer",
		ExpiredIn: 60 * 60 * 24,
		Key:       testKey,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
}

func TestParseToken(t *testing.T) {
	token, err := GenerateToken(&GenerateTokenRequest{
		UserId:    "test_user_id",
		Subject:   "subject",
		Issuer:    "issuer",
		ExpiredIn: 60 * 60 * 24,
		Key:       testKey,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
	parsed, err := ParseToken(&ParseTokenRequest{
		Token: token,
		Key:   testKey,
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	bytes, err := json.MarshalIndent(parsed, "", "    ")
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println(string(bytes))
}

func TestParseExpiredToken(t *testing.T) {
	token, err := GenerateToken(&GenerateTokenRequest{
		UserId:    "test_user_id",
		Subject:   "subject",
		Issuer:    "issuer",
		ExpiredIn: 3,
		Key:       testKey,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
	t.Log("waiting for 5 seconds")
	time.Sleep(time.Second * 5)
	_, err = ParseToken(&ParseTokenRequest{
		Token: token,
		Key:   testKey,
	})
	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Fatal("unexpected err: ", err)
	}
}
