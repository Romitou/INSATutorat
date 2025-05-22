package admin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
	"net/http"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		if err := database.Get().
			Find(&users).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.NotFound)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		// on est sur une route admin, on inclut les détails
		var privateUsers []models.PrivateUser
		for _, user := range users {
			privateUsers = append(privateUsers, user.ToPrivate())
		}

		c.JSON(http.StatusOK, privateUsers)
	}
}
