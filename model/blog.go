package model

import (
	"github.com/go-xorm/xorm"
	"time"
)

type Blog struct {
	Id  int `json:"id" xorm:"'id'"`
	UserID  int `json:"user_id" xorm:"'user_id'"`
	Title string `json:"title" xorm:"'title'"`
	Body string `json:"body" xorm:"'body'"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}
type BlogRepository struct {
	engine *xorm.Engine
}
func NewBlogRepository()(BlogRepository){

	return BlogRepository{
		engine:getEngine(),
	};
}
func(self BlogRepository)Create(blog Blog)(int64, error){
	return self.engine.Insert(blog)

}

func (self BlogRepository)Exist(user User)(bool){
	has,err := self.engine.Get(&user)
	if err != nil {
		panic(err)
	}
	return has
}
