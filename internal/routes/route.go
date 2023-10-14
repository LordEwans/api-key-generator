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

	user.POST("/key=:api-key", db.VerifyKeyAPI(), db.CreateUser())
	user.DELETE("/:api/id=:_id/key=:api-key", db.VerifyKeyAPI(), db.DeleteUser())

	r.GET("verify/:api/key=:api-key", db.VerifyKey())
}
