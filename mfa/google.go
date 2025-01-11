package mfa

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"image/png"
	"math"
	"strings"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gogf/gf/v2/text/gstr"
)

type GoogleAuthenticatorConfig struct {
	CodeLength   int    `json:"code_length"`
	SecretLength int    `json:"secret_length"`
	Name         string `json:"name"`
}

type GoogleAuthenticator struct {
	codeLength   int
	secretLength int
}

// NewGoogleAuthenticator constructor
func NewGoogleAuthenticator(codeLength, secretLength int) *GoogleAuthenticator {
	return &GoogleAuthenticator{
		codeLength:   codeLength,
		secretLength: secretLength,
	}
}

// CreateSecret in TOTP
func (g *GoogleAuthenticator) CreateSecret() (string, error) {
	if g.secretLength < 16 || g.secretLength > 128 {
		return "", errors.New("bad secret length")
	}

	validChars := base32LookupTable
	secret := make([]byte, g.secretLength)

	_, err := rand.Read(secret)
	if err != nil {
		return "", errors.New("no source of secure random")
	}

	var result strings.Builder
	for _, b := range secret {
		result.WriteByte(validChars[int(b)&31])
	}

	return result.String(), nil
}

// GetCode in TOTP with secret and time
func (g *GoogleAuthenticator) GetCode(secret string, ts ...int64) (string, error) {
	var timeSlice int64
	if len(ts) > 0 {
		timeSlice = ts[0]
	}
	if timeSlice == 0 {
		timeSlice = time.Now().Unix() / 30
	}

	secretKey, err := base32Decode(secret)
	if err != nil {
		return "", err
	}

	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, uint64(timeSlice))

	hash := hmac.New(sha1.New, secretKey)
	hash.Write(timeBytes)
	hmacHash := hash.Sum(nil)

	offset := hmacHash[len(hmacHash)-1] & 0x0F
	value := int32(binary.BigEndian.Uint32(hmacHash[offset:offset+4]) & 0x7FFFFFFF)

	return fmt.Sprintf("%06d", value%int32(math.Pow10(g.codeLength))), nil
}

// VerifyCode checks if the given code is valid for the provided secret and error range
func (g *GoogleAuthenticator) VerifyCode(secret string, code string, discrepancy int, ts ...int64) bool {
	if len(code) != g.codeLength {
		return false
	}

	var currentTimeSlice int64
	if len(ts) > 0 {
		currentTimeSlice = ts[0]
	}
	if currentTimeSlice == 0 {
		currentTimeSlice = time.Now().Unix() / 30
	}

	for i := -discrepancy; i <= discrepancy; i++ {
		calculatedCode, err := g.GetCode(secret, currentTimeSlice+int64(i))
		if err != nil {
			continue
		}
		if gstr.Equal(calculatedCode, code) {
			return true
		}
	}

	return false
}

// GenerateQRCode png and encode to base64
func (g *GoogleAuthenticator) GenerateQRCode(title string, secret string) (base64QR string, err error) {
	qrContent := fmt.Sprintf("otpauth://totp/%s?secret=%s", title, secret)

	qrCode, err := qr.Encode(qrContent, qr.L, qr.Auto)
	if err != nil {
		return
	}

	qrCode, err = barcode.Scale(qrCode, 250, 250)
	if err != nil {
		return
	}

	var pngData strings.Builder
	if err = png.Encode(&pngData, qrCode); err != nil {
		return
	}

	base64QR = base64.StdEncoding.EncodeToString([]byte(pngData.String()))

	return
}
