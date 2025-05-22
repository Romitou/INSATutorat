package assignments

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
)

type tutorAssignmentElement struct {
	models.Campaign
	Assignments []models.TutorSubjectDetailed `json:"assignments"`
}

func TutorAssignments() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)

		var openCampaigns []models.Campaign
		if err := database.Get().
			Where("school_year = ?", user.SchoolYear).
			Find(&openCampaigns).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, []tutorAssignmentElement{})
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		// on récupère les matières du tuteur
		var tutorSubjects []models.TutorSubject
		if err := database.Get().
			Where("tutor_id = ?", user.ID).
			Preload("Tutees").
			Preload("Tutees.Tutee").
			Preload("Subject").
			Find(&tutorSubjects).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, []tutorAssignmentElement{})
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		// on construit une réponse pour chaque campagne
		response := make([]tutorAssignmentElement, 0, len(openCampaigns))
		for _, campaign := range openCampaigns {
			// on ne garde que les matières de la campagne
			var detailedTutorSubjects []models.TutorSubjectDetailed
			for _, tutorSubject := range tutorSubjects {
				if tutorSubject.CampaignID != campaign.ID {
					continue
				}
				detailedTutorSubjects = append(detailedTutorSubjects, tutorSubject.ToDetailed())
			}

			// une fois les matières filtrées, on les ajoute à la réponse avec la campagne
			response = append(response, tutorAssignmentElement{
				Campaign:    campaign,
				Assignments: detailedTutorSubjects,
			})
		}

		c.JSON(http.StatusOK, response)
	}
}
