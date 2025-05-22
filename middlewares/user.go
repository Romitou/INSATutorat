package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"net/http"
)

func UserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.Get()
		session := sessions.Default(c)
		userID := session.Get("user_id")
		var user models.User

		switch id := userID.(type) {
		case uint:
			db.Where(&models.User{
				ID: id,
			}).First(&user)
		default:
			// userID absent ou mauvais type
			_ = c.Error(apierrors.Unauthorized)
			c.Status(http.StatusForbidden)
			c.Abort()
			return
		}

		if user.IsEmpty() {
			_ = c.Error(apierrors.Unauthorized)
			c.Status(http.StatusForbidden)
			c.Abort()
			return
		} else {
			c.Set("user", user)
			c.Next()
		}
	}
}
