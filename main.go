package main

import (
	"github.com/gin-gonic/gin"
	"github.com/makki0205/blog/controller"
	"github.com/itsjamie/gin-cors"
	"github.com/makki0205/blog/env"
	"github.com/makki0205/blog/middleware"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "yes")
	})

	r.Use(cors.Middleware(env.CorsConfig))
	r.Use(gin.Logger())

	user := r.Group("/api")

	uctr := controller.NewUser()
	user.POST("/signup", uctr.SignUp)
	user.POST("/signin", middleware.AuthMiddleware.LoginHandler)


	blog := r.Group("/api")
	blog.Use(middleware.AuthMiddleware.MiddlewareFunc("user"))

	bctr := controller.NewBlog()
	blog.GET("/blog", bctr.GetAll)
	blog.GET("/blog/:id", bctr.GetByID)
	blog.POST("/blog", bctr.Create)
	blog.PUT("/blog/:id", bctr.Update)
	blog.DELETE("/blog", bctr.Delete)

	r.Run(":3001")
}