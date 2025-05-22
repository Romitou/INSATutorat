package auth

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
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

		log.Println(os.Getenv("DEV_MODE") != "true")
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

		log.Println(os.Getenv("DEV_MODE") != "true")
		if user.LoginRequestedAt.IsZero() {
			_ = c.Error(apierrors.Unauthorized)
			return
		}

		// le login token est valide pendant 15 minutes
		log.Println(os.Getenv("DEV_MODE") != "true")
		if os.Getenv("DEV_MODE") != "true" && user.LoginRequestedAt.Add(15*time.Minute).Before(time.Now()) {
			_ = c.Error(apierrors.Unauthorized)
			return
		}

		// on met à jour la session
		session := sessions.Default(c)
		session.Set("user_id", user.ID)
		err := session.Save()
		if err != nil {
			_ = c.Error(err)
		}

		c.Status(http.StatusOK)
	}
}
