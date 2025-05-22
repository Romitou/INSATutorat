package assignments

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
	"net/http"
)

type TuteeAssignment struct {
	ID      *uint          `json:"id"`
	Subject models.Subject `json:"subject"`
	Tutor   models.User    `json:"tutor"`

	TotalHours int `json:"totalHours"`
}

type tuteeAssignmentElement struct {
	models.Campaign
	Assignments []TuteeAssignment `json:"assignments"`
}

func TuteeAssignments() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)

		var openCampaigns []models.Campaign
		if err := database.Get().
			Where("school_year = ?", user.SchoolYear).
			Find(&openCampaigns).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, []tuteeAssignmentElement{})
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		// on récupère les inscriptions du tutoré
		var tuteeRegistrations []models.TuteeRegistration
		if err := database.Get().
			Where("tutee_id = ?", user.ID).
			Preload("TutorSubject").
			Preload("TutorSubject.Tutor").
			Preload("TutorSubject.Subject").
			Preload("TutorSubject.Campaign").
			Find(&tuteeRegistrations).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, []tuteeAssignmentElement{})
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		assignmentMap := make(map[uint][]TuteeAssignment, len(tuteeRegistrations))
		// pour chaque inscription, on l'ajoute au tableau d'assignations correspondant à la campagne,
		// afin de retourner un tableau d'assignations par campagne
		for _, tuteeRegistration := range tuteeRegistrations {
			assignmentMap[tuteeRegistration.TutorSubject.CampaignID] = append(
				assignmentMap[tuteeRegistration.TutorSubject.CampaignID],
				TuteeAssignment{
					ID:      tuteeRegistration.TutorSubjectID,
					Subject: tuteeRegistration.TutorSubject.Subject,
					Tutor:   tuteeRegistration.TutorSubject.Tutor,
				},
			)
		}

		// on construit la réponse en ajoutant les assignations et la campagne associée
		response := make([]tuteeAssignmentElement, 0, len(openCampaigns))
		for _, campaign := range openCampaigns {
			response = append(response, tuteeAssignmentElement{
				Campaign:    campaign,
				Assignments: assignmentMap[campaign.ID],
			})
		}

		c.JSON(http.StatusOK, response)
	}
}
