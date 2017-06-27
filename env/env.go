package env

import (
	"github.com/itsjamie/gin-cors"
	"time"
)

var CorsConfig = cors.Config{
	Origins : "*",
	Methods : "GET, PUT, POST, DELETE",
	RequestHeaders : "Origin, Authorization, Content-Type",
	ExposedHeaders : "",
	MaxAge : 50 * time.Second,
	Credentials : true,
	ValidateHeaders : false,
}