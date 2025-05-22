package tutor

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
	"net/http"
)

type subjectWithMaxTutees struct {
	models.Subject
	MaxTutees int `json:"maxTutees"`
}

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
			Where("school_year = ?", user.SchoolYear).
			First(&campaign).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.NotFound)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		var existingRegistrations []models.TutorSubject
		if err := database.Get().
			Where("tutor_id = ?", user.ID).
			Where("campaign_id = ?", campaign.ID).
			Preload("Subject").
			Find(&existingRegistrations).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				apierrors.DatabaseError(c, err)
				return
			}
		}

		// on construit la liste des matières avec le nombre maximum de tutorés
		subjects := make([]subjectWithMaxTutees, 0, len(existingRegistrations))
		for _, existingRegistration := range existingRegistrations {
			subjects = append(subjects, subjectWithMaxTutees{
				Subject:   existingRegistration.Subject,
				MaxTutees: existingRegistration.MaxTutees,
			})
		}

		c.JSON(http.StatusOK, subjects)
	}
}
