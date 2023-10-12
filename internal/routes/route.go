package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lordewans/api-key-generator/internal/database"
)

var (
	db, _ = database.ConnectDB()
)

func UseRoute(r *gin.Engine) {
	user := r.Group("/user")

	user.POST("/create", db.CreateUser())
	user.DELETE("/:_id")

	r.GET("api/:api-key")
}
