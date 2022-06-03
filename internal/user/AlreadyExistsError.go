package user

import "fmt"

type alreadyExists struct {
	login string
}

func (e *alreadyExists) Error() string {
	return fmt.Sprintf("User with login %s already exists.", e.login)
}
