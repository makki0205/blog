package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/makki0205/blog/model"
	"github.com/makki0205/blog/httputil"
)

type User struct {

}
func NewUser() User{
	return User{}
}
func(u *User) SignUp(c *gin.Context){
	var user model.User
	if err :=c.BindJSON(&user); err != nil{
		c.JSON(400, httputil.NewErrorResponse(err))
		return
	}
	userRep := model.NewUserRepository()

	err := userRep.Create(user.Username, user.Password)
	if err != nil{
		c.JSON(400, httputil.NewErrorResponse(err))
		return
	}
	c.JSON(200, httputil.NewSuccessResponse("ok"))
}

func(u *User) SignIn(c *gin.Context){
	c.JSON(200, gin.H{
		"hoge": "hoge",
	})
}
