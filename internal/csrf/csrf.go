package csrf

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var Tokens *HashToken

// nolint:gochecknoinits
func init() {
	Tokens = NewHMACHashToken(os.Getenv("CSRF_SECRET"))
}

type HashToken struct {
	Secret []byte
}

func NewHMACHashToken(secret string) *HashToken {
	if err := godotenv.Load(".env"); err != nil {
		return &HashToken{Secret: []byte(secret)}
	}
	return &HashToken{Secret: []byte(secret)}
}

func (tk *HashToken) Create(session string, tokenExpTime int64) (string, error) {
	h := hmac.New(sha256.New, tk.Secret)
	data := fmt.Sprintf("%s:%d", session, tokenExpTime)
	if _, err := h.Write([]byte(data)); err != nil {
		return "", err
	}
	token := hex.EncodeToString(h.Sum(nil)) + ":" + strconv.FormatInt(tokenExpTime, 10)
	return token, nil
}

func (tk *HashToken) Check(session string, inputToken string) (bool, error) {
	tokenData := strings.Split(inputToken, ":")
	if len(tokenData) != 2 {
		return false, fmt.Errorf("bad token data")
	}

	tokenExp, err := strconv.ParseInt(tokenData[1], 10, 64)
	if err != nil {
		return false, fmt.Errorf("bad token time")
	}

	if tokenExp < time.Now().Unix() {
		return false, fmt.Errorf("token expired")
	}

	h := hmac.New(sha256.New, tk.Secret)
	data := fmt.Sprintf("%s:%d", session, tokenExp)
	if _, err = h.Write([]byte(data)); err != nil {
		return false, err
	}
	expectedMAC := h.Sum(nil)
	messageMAC, err := hex.DecodeString(tokenData[0])
	if err != nil {
		return false, fmt.Errorf("can't hex decode token")
	}

	return hmac.Equal(messageMAC, expectedMAC), nil
}
