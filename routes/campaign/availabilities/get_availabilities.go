package availabilities

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
	"net/http"
)

func GetAvailabilities() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)

		campaignId := c.Param("campaignId")
		if campaignId == "" {
			_ = c.Error(apierrors.BadRequest)
			return
		}

		var semesterAvailability models.SemesterAvailability
		if err := database.Get().
			Where("user_id = ?", user.ID).
			Where("campaign_id = ?", campaignId).
			First(&semesterAvailability).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, &models.Slots{})
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		// on désérialise le JSON
		var slots models.Slots
		err := json.Unmarshal([]byte(semesterAvailability.AvailabilityJSON), &slots)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, slots)
	}
}
