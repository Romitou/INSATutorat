package tutee

import (
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
)

func GetRegistrations() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)

		campaignId := c.Param("campaignId")
		if campaignId == "" {
			_ = c.Error(apierrors.BadRequest)
			return
		}

		var campaign models.Campaign
		if err := database.Get().
			Where("id = ?", campaignId).
			Where("school_year = ?", os.Getenv("SCHOOL_YEAR")).
			First(&campaign).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.NotFound)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		var existingRegistrations []models.TuteeRegistration
		if err := database.Get().
			Where("tutee_id = ?", user.ID).
			Where("campaign_id = ?", campaign.ID).
			Preload("Subject").
			Find(&existingRegistrations).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				apierrors.DatabaseError(c, err)
				return
			}
		}

		subjects := make([]models.Subject, 0, len(existingRegistrations))
		for _, existingRegistration := range existingRegistrations {
			subjects = append(subjects, existingRegistration.Subject)
		}

		c.JSON(http.StatusOK, subjects)
	}
}
