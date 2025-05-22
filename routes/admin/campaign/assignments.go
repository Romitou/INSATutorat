package campaign

import (
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"net/http"
)

type assignmentsJson struct {
	Tutees        []models.TuteeRegistration `json:"tutees"`
	TutorSubjects []models.TutorSubject      `json:"tutorSubjects"`
}

func GetAssignments() gin.HandlerFunc {
	return func(c *gin.Context) {

		var tuteeRegs []models.TuteeRegistration
		var tutorRegs []models.TutorSubject

		db := database.Get()

		campaignId := c.Param("campaignId")
		if campaignId == "" {
			_ = c.Error(apierrors.BadRequest)
			return
		}

		if err := db.
			Where("campaign_id = ?", campaignId).
			Preload("Tutee").
			Find(&tuteeRegs).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		if err := db.
			Where("campaign_id = ?", campaignId).
			Preload("Tutor").
			Find(&tutorRegs).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		c.JSON(http.StatusOK, &assignmentsJson{
			Tutees:        tuteeRegs,
			TutorSubjects: tutorRegs,
		})
	}
}
