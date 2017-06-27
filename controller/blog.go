package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/makki0205/blog/model"
	_ "fmt"
	"github.com/makki0205/blog/httputil"
)

type Blog struct {
	blogRep model.BlogRepository
	userRep model.UserRepository
}
func NewBlog() Blog{
	return Blog{
		blogRep: model.NewBlogRepository(),
	}
}


func(b *Blog) GetAll(c *gin.Context){
	c.JSON(200, gin.H{
		"hoge": "hoge",
	})
}

func(b *Blog) GetByID(c *gin.Context){
	c.JSON(200, gin.H{
		"hoge": "hoge",
	})
}
func(b *Blog) Create(c *gin.Context){
	//リクエスト処理
	var blog model.Blog
	err := c.BindJSON(&blog)
	if err!=nil {
		panic(err)
	}
	blog.UserID = b.userRep.GetUserId(c)
	// Insert
	ok, err := b.blogRep.Create(blog)
	if err !=nil {
		panic(err)
	}

	c.JSON(200, httputil.NewSuccessResponse(ok))
}
func(b *Blog) Update(c *gin.Context){
	c.JSON(200, gin.H{
		"hoge": "hoge",
	})
}

func(b *Blog) Delete(c *gin.Context){
	c.JSON(200, gin.H{
		"hoge": "hoge",
	})
}



