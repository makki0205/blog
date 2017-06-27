package middleware

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/makki0205/blog/model"
)

var AuthMiddleware = GinJWTMiddleware{
	Realm:      "api",
	Key:        []byte("hogehoge"),
	Timeout:    time.Hour * 24,
	MaxRefresh: time.Hour * 24 * 14,
	Authenticator: authenticatorHandler,
	Unauthorized: func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	},
	TimeFunc: time.Now,
}

func authenticatorHandler(username string, password string, c *gin.Context)(string, bool)  {
	user := model.NewUserRepository()
	return username, user.Exist(model.User{Username:username,Password:password})
}