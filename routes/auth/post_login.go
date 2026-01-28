package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
)

type loginJson struct {
	LoginToken string `json:"token" binding:"required"`
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input loginJson
		if err := c.ShouldBindJSON(&input); err != nil {
			_ = c.Error(err)
			return
		}

		var user models.User
		result := database.Get().Where(&models.User{
			LoginToken: input.LoginToken,
		}).First(&user)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.Unauthorized)
				return
			}
			apierrors.DatabaseError(c, result.Error)
			return
		}

		if user.LoginRequestedAt.IsZero() {
			_ = c.Error(apierrors.Unauthorized)
			return
		}

		// le login token est valide pendant 15 minutes
		// La vérification peut être désactivée avec CHECK_TOKEN_EXPIRATION=false
		checkExpiration := os.Getenv("CHECK_TOKEN_EXPIRATION") != "false"
		if checkExpiration && user.LoginRequestedAt.Add(15*time.Minute).Before(time.Now()) {
			_ = c.Error(apierrors.Unauthorized)
			return
		}

		// on met à jour la session
		session := sessions.Default(c)
		session.Clear()
		session.Set("user_id", user.ID)
		err := session.Save()
		if err != nil {
			_ = c.Error(err)
		}

		c.Status(http.StatusOK)
	}
}
