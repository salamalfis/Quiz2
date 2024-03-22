package handler

import (

    "net/http"
 
    "strings"

    "github.com/gin-gonic/gin"
)

type User struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

var users []User

func GetUsers(ctx *gin.Context) {
    // Get all users or filter by email
    email := ctx.Query("email")
    if email != "" {
        arr := []User{}
        for _, user := range users {
            if strings.Contains(user.Email, email) {
                arr = append(arr, user)
            }
        }
        ctx.JSON(http.StatusOK, arr)
        return
    }
    ctx.JSON(http.StatusOK, users)
}

