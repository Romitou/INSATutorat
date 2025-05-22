package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database/models"
)

func AdminHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			_ = c.Error(apierrors.Unauthorized)
			return
		}

		user, ok := userInterface.(models.User)
		if !ok {
			_ = c.Error(errors.New("user is not a valid user")) // ne dois jamais se produire
			return
		}

		if !user.IsAdmin {
			_ = c.Error(apierrors.Forbidden)
			return
		}

		c.Next()
	}
}
