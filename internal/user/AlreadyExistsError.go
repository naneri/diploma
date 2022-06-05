package user

import "fmt"

type AlreadyExistsError struct {
	login string
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("User with login %s already exists.", e.login)
}
