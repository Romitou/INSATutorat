package campaign

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/core"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
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
				var existing models.TuteeRegistration
				if err = db.Preload("Tutee").Where("id = ?", tr.ID).First(&existing).Error; err != nil {
					apierrors.DatabaseError(c, err)
					return
				}

				// on met à jour le tutor_subject_id uniquement (l'assignation)
				err = db.Model(&models.TuteeRegistration{}).
					Where("id = ?", tr.ID).
					Update("tutor_subject_id", tr.TutorSubjectID).Error
				if err != nil {
					apierrors.DatabaseError(c, err)
					return
				}

				// on ne notifie que si l'affectation vient d'être créée ou changée, pas à
				// chaque enregistrement de la page (qui renvoie systématiquement toute la liste)
				if newlyAssigned(existing.TutorSubjectID, tr.TutorSubjectID) {
					notifyNewAssignment(c, db, existing.Tutee, *tr.TutorSubjectID)
				}
			} else {
				// sinon, on insère
				if err = db.Create(&tr).Error; err != nil {
					apierrors.DatabaseError(c, err)
					return
				}

				if tr.TutorSubjectID != nil {
					var fullTutee models.User
					if err = db.First(&fullTutee, tr.TuteeID).Error; err == nil {
						notifyNewAssignment(c, db, fullTutee, *tr.TutorSubjectID)
					}
				}
			}
		}

		c.Status(http.StatusOK)
	}
}

// newlyAssigned détermine si une affectation vient d'être créée ou modifiée, càd si
// le tutor_subject_id passe de nil à une valeur, ou change de valeur
func newlyAssigned(previous, next *uint) bool {
	if next == nil {
		return false
	}
	return previous == nil || *previous != *next
}

// notifyNewAssignment envoie un email au tutoré et au tuteur concernés par une
// nouvelle affectation. Un échec d'envoi est loggué mais ne fait pas échouer la
// requête : l'affectation en base de données reste la source de vérité.
func notifyNewAssignment(c *gin.Context, db *gorm.DB, tutee models.User, tutorSubjectID uint) {
	var ts models.TutorSubject
	if err := db.Preload("Tutor").Preload("Subject").First(&ts, tutorSubjectID).Error; err != nil {
		apierrors.LogError(c, err)
		return
	}

	if err := core.SendAssignmentNotificationToTutee(tutee, ts.Tutor, ts.Subject, ts.ID); err != nil {
		apierrors.LogError(c, err)
	}
	if err := core.SendAssignmentNotificationToTutor(ts.Tutor, tutee, ts.Subject, ts.ID); err != nil {
		apierrors.LogError(c, err)
	}
}
