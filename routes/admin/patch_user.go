package admin

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
)

type patchUserInput struct {
	IsAdmin bool `json:"isAdmin"`
	IsTutor bool `json:"isTutor"`
	IsTutee bool `json:"isTutee"`
}

func PatchUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdStr := c.Param("id")
		if userIdStr == "" {
			_ = c.Error(apierrors.BadRequest)
			return
		}
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			_ = c.Error(apierrors.BadRequest)
			return
		}

		var input patchUserInput
		if err = c.ShouldBindJSON(&input); err != nil {
			_ = c.Error(err)
			return
		}

		var targetUser models.User
		if err = database.Get().
			Where("id = ?", userId).
			First(&targetUser).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.NotFound)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		// un admin ne peut pas se retirer lui-même ses droits admin, pour éviter un verrouillage
		currentUser := c.MustGet("user").(models.User)
		if targetUser.ID == currentUser.ID && input.IsAdmin != targetUser.IsAdmin {
			_ = c.Error(apierrors.Forbidden)
			return
		}

		targetUser.IsAdmin = input.IsAdmin
		targetUser.IsTutor = input.IsTutor
		targetUser.IsTutee = input.IsTutee

		if err = database.Get().Save(&targetUser).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		c.JSON(http.StatusOK, targetUser.ToPrivate())
	}
}
