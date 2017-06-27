package model

import (
	"github.com/go-xorm/xorm"
	"github.com/pkg/errors"
	"github.com/gin-gonic/gin"
)

type User struct {
	Id  int `json:"id" xorm:"'id'"`
	Username string `json:"username" xorm:"'username'"`
	Password string `json:"password" xorm:"'password'"`
}
type UserRepository struct {
	engine *xorm.Engine
}
func NewUserRepository()(UserRepository){

	return UserRepository{
		engine:getEngine(),
	};
}
func(self UserRepository)Create(email string,pass string )(error){
	user := User{
		Username: email,
	}
	if self.Exist(user) {
		return errors.New("ユーザーネームはすでに使用されています")
	}
	user.Password = pass
	_, err  :=self.engine.Insert(&user)
	return err
}

func (self UserRepository)Exist(user User)(bool){
	has,err := self.engine.Get(&user)
	if err != nil {
		panic(err)
	}
	return has
}

func(self UserRepository)getUserId(username string)(int){
	user := User{
		Username:username,
	}
	has,err := getEngine().Get(&user)
	if has {
		return user.Id
	}
	if err != nil {
		panic(err)
	}
	return 0
}

func(self UserRepository)GetUserId(c *gin.Context)(int){
	username,_ := c.Get("userID")
	username_string := username.(string)
	userId := self.getUserId(username_string)
	return userId
}