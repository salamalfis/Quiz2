package repository

import (
    "github.com/salamalfis/Golang-DTS/helper/handler"
)

func GetUserByID(id uint) (handler.User, error) {
    for _, user := range handler.users {
        if user.ID == id {
            return user, nil
        }
    }
    return handler.User{}, ErrUserNotFound
}

// Other repository functions go here...

var ErrUserNotFound = errors.New("user not found")
