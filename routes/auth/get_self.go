package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/database/models"
	"net/http"
)

func Self() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)
		// on récupère soi-même :) avec ses données "privées"
		c.JSON(http.StatusOK, user.ToPrivate())
	}
}
