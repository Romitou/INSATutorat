package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/core"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type sendLinkJson struct {
	MailAddress string `json:"mail" binding:"required,email"`
}

func SendLink() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input sendLinkJson
		if err := c.ShouldBindJSON(&input); err != nil {
			_ = c.Error(err)
			return
		}

		var user models.User
		result := database.Get().Where(&models.User{
			Mail: input.MailAddress,
		}).First(&user)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.EmailNotRegistered)
				return
			}
			apierrors.DatabaseError(c, result.Error)
			return
		}

		uuidToken, err := uuid.NewV4()
		if err != nil {
			_ = c.Error(err)
			return
		}

		// le login ne sera possible que durant 15 minutes
		user.LoginToken = uuidToken.String()
		user.LoginRequestedAt = time.Now()

		err = database.Get().Save(&user).Error
		if err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		// envoi de l'email
		err = core.SendLoginLink(user, user.LoginToken)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusOK)
	}
}
