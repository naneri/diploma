package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/http"
)

type UserID string

const UserIDContextKey = "UserID"

var secretKey = []byte("secret key")
var userID uint32

func IDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			data   []byte
			err    error
			idSign []byte
		)

		// parse cookie
		cookie, err := r.Cookie("user")

		// can I make this code prettier?
		if err != nil {
			fmt.Println("error getting cookie: " + err.Error())
			http.Error(w, "authentication error", http.StatusUnauthorized)
		} else {
			data, err = hex.DecodeString(cookie.Value)

			if err != nil {
				fmt.Println("error decoding cookie: " + err.Error())
				http.Error(w, "authentication error", http.StatusUnauthorized)
			} else {
				userID = binary.BigEndian.Uint32(data[:4])
				h := hmac.New(sha256.New, secretKey)
				h.Write(data[:4])
				idSign = h.Sum(nil)

				// if parse correctly, add the cookie to context
				if !hmac.Equal(idSign, data[4:]) {
					fmt.Println("wrong sign")
					http.Error(w, "authentication error", http.StatusUnauthorized)
				}
			}
		}

		ctx := r.Context()
		req := r.WithContext(context.WithValue(ctx, UserID(UserIDContextKey), userID))
		*r = *req
		// else grant user the signed cookie with Unique identifier
		next.ServeHTTP(w, r)
	})
}
