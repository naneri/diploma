package services

import "github.com/naneri/diploma/internal/user"

func CheckUserPassword(password string, user user.User) bool {
	return password == user.Password
}
