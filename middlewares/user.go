package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
)

// permet d'amorcer le tout premier compte admin sans accéder à la base de données :
// si l'email de l'utilisateur figure dans INITIAL_ADMIN_EMAILS, il est promu admin.
// les admins suivants se gèrent ensuite depuis l'interface (/admin/users).
func promoteInitialAdmin(db *gorm.DB, user *models.User) {
	if user.IsAdmin {
		return
	}

	adminEmails := os.Getenv("INITIAL_ADMIN_EMAILS")
	if adminEmails == "" {
		return
	}

	for _, mail := range strings.Split(adminEmails, ",") {
		if strings.EqualFold(strings.TrimSpace(mail), user.Mail) {
			user.IsAdmin = true
			db.Save(user)
			return
		}
	}
}

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
			promoteInitialAdmin(db, &user)
			c.Set("user", user)
			c.Next()
		}
	}
}
