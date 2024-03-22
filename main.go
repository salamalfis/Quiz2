package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

var users []User

func main() {
	GinHttp()
}

func GinHttp() {
	// gin => framework HTTP punya golang
	// big community
	engine := gin.New()

	// serve static template
	// engine.LoadHTMLGlob("static/*")
	engine.Static("/static", "./static")

	engine.LoadHTMLGlob("template/*")
	engine.GET("/template/index/:name", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.tmpl", map[string]any{
			"title": ctx.Param("name"),
		})
	})

	// membuat prefix group
	v1 := engine.Group("/api/v1")
	{
		usersGroup := v1.Group("/users")
		{
			// [GET] /api/v1/users
			// filter user by email
			usersGroup.GET("", func(ctx *gin.Context) {
				email := ctx.Query("email")
				if email != "" {
					arr := []User{}
					for _, user := range users {
						// full text search
						if strings.Contains(user.Email, email) {
							arr = append(arr, user)
						}
					}
					ctx.JSON(http.StatusOK, arr)
					return
				}
				ctx.JSON(http.StatusOK, users)
			})

			// [POST] /api/v1/users
			usersGroup.POST("", func(ctx *gin.Context) {
				// binding payload
				user := User{}
				if err := ctx.Bind(&user); err != nil {
					ctx.JSON(http.StatusBadRequest, map[string]any{
						"message": "failed to bind body",
					})
					return
				}
				user.ID = uint(len(users) + 1)
				users = append(users, user)
				ctx.JSON(http.StatusAccepted, map[string]any{
					"message": "user created",
				})
			})

			// [GET] /api/v1/users/:id
			
			usersGroup.GET("/:id", func(ctx *gin.Context) {
				id, err := strconv.Atoi(ctx.Param("id"))
				if err != nil || id <= 0 {
					ctx.JSON(http.StatusBadRequest, map[string]any{
						"message": "invalid ID",
					})
					return
				}
				for _, user := range users {
					if user.ID == uint(id) {
						ctx.JSON(http.StatusOK, user)
						return
					}
				}

				ctx.JSON(http.StatusNotFound, map[string]any{
					"message": "user not found",
				})
			})

			// [PUT] /api/v1/users/:id
			usersGroup.PUT("/:id", func(ctx *gin.Context) {
				id, err := strconv.Atoi(ctx.Param("id"))
				if err != nil || id <= 0 {
					ctx.JSON(http.StatusBadRequest, map[string]any{
						"message": "invalid ID",
					})
					return
				}

				user := User{}
				if err := ctx.Bind(&user); err != nil {
					ctx.JSON(http.StatusBadRequest, map[string]any{
						"message": "failed to bind body",
					})
					return
				}

				for i, u := range users {
					if u.ID == uint(id) {
						users[i] = user
						ctx.JSON(http.StatusAccepted, map[string]any{
							"message": "user updated",
						})
						return
					}
				}

				ctx.JSON(http.StatusNotFound, map[string]any{
					"message": "user not found",
				})
			})

			// [DELETE] /api/v1/users/:id
			usersGroup.DELETE("/:id", func(ctx *gin.Context) {
				id, err := strconv.Atoi(ctx.Param("id"))
				if err != nil || id <= 0 {
					ctx.JSON(http.StatusBadRequest, map[string]any{
						"message": "invalid ID",
					})
					return
				}

				for i, user := range users {
					if user.ID == uint(id) {
						users = append(users[:i], users[i+1:]...)
						ctx.JSON(http.StatusOK, map[string]any{
							"message": "user deleted",
						})
						return
					}
				}

				ctx.JSON(http.StatusNotFound, map[string]any{
					"message": "user not found",
				})
			})
		}

		orderGroup := v1.Group("/orders")
		{
			orderGroup.GET("", func(ctx *gin.Context) {

			})
		}
	}

	engine.Run(":8080")
}

func NetHttp() {
	// /users => API path
	// func(w http.ResponseWriter, r *http.Request) => handler function
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		// validasi request payload (header, body)
		// memanggil logic
		// memberkan response

		// get all users
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(users)
			return
		}
		// create user
		if r.Method == http.MethodPost {
			user := User{}
			// only bind username and email
			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			user.ID = uint(len(users) + 1)
			users = append(users, user)

			w.WriteHeader(http.StatusAccepted)
			return
		}

		

		// mini quiz
		// buatlah method
		// PUT /users/:id untuk edit user by id
		if r.Method == http.MethodPut {
			pathValue, _ := strconv.Atoi(r.PathValue("id"))
			user := User{}
			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			for i, u := range users {
				if u.ID == uint(pathValue) {
					users[i] = user
					w.WriteHeader(http.StatusAccepted)
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Delete /users/:id untuk delete user by id
		if r.Method == http.MethodDelete {
			pathValue, _ := strconv.Atoi(r.PathValue("id"))
			for i, user := range users {
				if user.ID == uint(pathValue) {
					users = append(users[:i], users[i+1:]...)
					w.WriteHeader(http.StatusOK)
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}
	})

	// {id} => path variable
	http.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			pathValue, _ := strconv.Atoi(r.PathValue("id"))
			for _, user := range users {
				if user.ID == uint(pathValue) {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(user)
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
		}
	})



	// :8080 PORT
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}


