package httpservices

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"net/http"
	"time"
)

func GenerateUserCookie(userID uint32, secretKey []byte) http.Cookie {
	uint32userIDBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(uint32userIDBuf[0:], userID)

	hash := hmac.New(sha256.New, secretKey)
	hash.Write(uint32userIDBuf)
	sign := hash.Sum(uint32userIDBuf)
	userCookie := hex.EncodeToString(sign)

	expire := time.Now().Add(10 * time.Minute)
	httpCookie := http.Cookie{Name: "user", Value: userCookie, Path: "/", Expires: expire, MaxAge: 90000}

	return httpCookie
}
