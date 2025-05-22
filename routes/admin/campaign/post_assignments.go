package campaign

import (
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"net/http"
	"strconv"
)

type saveAssignmentsInput struct {
	Tutees        []models.TuteeRegistration `json:"tutees"`
	TutorSubjects []models.TutorSubject      `json:"tutorSubjects"`
}

func PostAssignments() gin.HandlerFunc {
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

		var input saveAssignmentsInput
		if err = c.ShouldBindJSON(&input); err != nil {
			_ = c.Error(err)
			return
		}

		for _, ts := range input.TutorSubjects {
			ts.CampaignID = uint(campaignId)
			// si le tutorSubject existe déjà, on le met à jour
			if ts.ID != 0 {
				// on met à jour le max_tutees uniquement
				err = db.Model(&models.TutorSubject{}).
					Where("id = ? AND campaign_id = ?", ts.ID, campaignId).
					Updates(map[string]interface{}{
						"max_tutees": ts.MaxTutees,
					}).Error
				if err != nil {
					apierrors.DatabaseError(c, err)
					return
				}
			} else {
				// sinon, on insère
				if err = db.Create(&ts).Error; err != nil {
					apierrors.DatabaseError(c, err)
					return
				}
			}
		}

		// on fait pareil pour les tutorés
		for _, tr := range input.Tutees {
			tr.CampaignID = uint(campaignId)
			// si le tuteeRegistration existe déjà, on le met à jour
			if tr.ID != 0 {
				// on met à jour le tutor_subject_id uniquement (l'assignation)
				err = db.Model(&models.TuteeRegistration{}).
					Where("id = ?", tr.ID).
					Update("tutor_subject_id", tr.TutorSubjectID).Error
				if err != nil {
					apierrors.DatabaseError(c, err)
					return
				}
			} else {
				// sinon, on insère
				if err = db.Create(&tr).Error; err != nil {
					apierrors.DatabaseError(c, err)
					return
				}
			}
		}

		c.Status(http.StatusOK)
	}
}
