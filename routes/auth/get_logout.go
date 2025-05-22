package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete("user_id")
		session.Delete("user")
		err := session.Save()
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.Status(http.StatusOK)
	}
}
