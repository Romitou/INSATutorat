package campaign

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.Get()

		campaignIdStr := c.Param("campaignId")
		if campaignIdStr == "" {
			_ = c.Error(apierrors.BadRequest)
			return
		}
		campaignId, err := strconv.Atoi(campaignIdStr)
		if err != nil {
			_ = c.Error(apierrors.BadRequest)
			return
		}

		var campaign models.Campaign
		if err = database.Get().
			Where("id = ?", campaignId).
			First(&campaign).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.NotFound)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		var users []models.User
		if err = db.Find(&users).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		// on est sur une route admin, on inclut les détails
		privateUsers := make([]models.PrivateUser, 0, len(users))
		for _, user := range users {
			privateUsers = append(privateUsers, user.ToPrivate())
		}

		c.JSON(200, privateUsers)
	}
}
