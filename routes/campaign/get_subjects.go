package campaign

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

func Subjects() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		var subjects []models.Subject
		if err := database.Get().
			Where("semester = ?", campaign.Semester).
			Find(&subjects).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, []models.Subject{})
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		c.JSON(http.StatusOK, subjects)
	}
}
