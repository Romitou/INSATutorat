package tutee

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
	"net/http"
)

type postRegisterJson struct {
	Subjects []uint `json:"subjects" binding:"required"`
}

func PostRegistrations() gin.HandlerFunc {
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

		var semesterAvailability models.SemesterAvailability
		if err := database.Get().
			Where("user_id = ?", user.ID).
			Where("campaign_id = ?", campaignId).
			Find(&semesterAvailability).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.BadRequest)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		if semesterAvailability.ID == 0 {
			_ = c.Error(apierrors.BadRequest)
			return
		}

		var registerJson postRegisterJson
		if err := c.ShouldBindJSON(&registerJson); err != nil {
			_ = c.Error(err)
			return
		}

		// on récupère les matières passées en JSON
		var subjects []models.Subject
		if err := database.Get().
			Where("id IN ?", registerJson.Subjects).
			Find(&subjects).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.NotFound)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		if len(subjects) == 0 {
			_ = c.Error(apierrors.NotFound)
			return
		}

		// validation des matières, on vérifie qu'elles sont toutes dans la même campagne
		for _, subject := range subjects {
			if subject.Semester != campaign.Semester {
				_ = c.Error(apierrors.BadRequest)
				return
			}
		}

		if len(subjects) != len(registerJson.Subjects) {
			_ = c.Error(apierrors.BadRequest)
			return
		}

		var existingRegistrations []models.TuteeRegistration
		if err := database.Get().
			Where("tutee_id = ?", user.ID).
			Where("campaign_id = ?", campaign.ID).
			Find(&existingRegistrations).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				apierrors.DatabaseError(c, err)
				return
			}
		}

		// pour chaque sujet, on vérifie s'il existe déjà une inscription
		for _, subject := range subjects {
			registration := models.TuteeRegistration{
				TuteeID:    user.ID,
				CampaignID: campaign.ID,
				SubjectID:  subject.ID,
			}

			var existingRegistration models.TuteeRegistration
			if err := database.Get().
				Where("tutee_id = ?", user.ID).
				Where("campaign_id = ?", campaign.ID).
				Where("subject_id = ?", subject.ID).
				Find(&existingRegistration).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					apierrors.DatabaseError(c, err)
					return
				}
			}

			// s'il n'existe pas, on l'ajoute
			if existingRegistration.ID == 0 {
				if err := database.Get().Create(&registration).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						_ = c.Error(apierrors.NotFound)
						return
					}
					apierrors.DatabaseError(c, err)
					return
				}
			}
		}

		// maintenant, on supprime les inscriptions qui ne sont plus dans le JSON
		for _, existingRegistration := range existingRegistrations {
			stillInSubjects := false
			for _, subject := range subjects {
				if existingRegistration.SubjectID == subject.ID {
					stillInSubjects = true
				}
			}
			if !stillInSubjects {
				// la matière n'est plus dans le JSON, on la supprime
				if err := database.Get().Delete(&existingRegistration).Error; err != nil {
					if !errors.Is(err, gorm.ErrRecordNotFound) {
						apierrors.DatabaseError(c, err)
						return
					}
				}
			}
		}

		c.Status(http.StatusOK)
	}
}
