package hours

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type patchHourJson struct {
	ID        uint   `json:"id" binding:"required"`
	StartDate string `json:"startDate" binding:"required"`
	EndDate   string `json:"endDate" binding:"required"`
}

func PatchHour() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)

		var input patchHourJson
		if err := c.ShouldBindJSON(&input); err != nil {
			_ = c.Error(err)
			return
		}

		tutorSubjectId := c.Param("tutorSubjectId")
		if tutorSubjectId == "" {
			_ = c.Error(apierrors.BadRequest)
			return
		}

		hourId := c.Param("hourId")
		if hourId == "" {
			_ = c.Error(apierrors.BadRequest)
			return
		}

		var tutorSubject models.TutorSubject
		if err := database.Get().
			Where("id = ?", tutorSubjectId).
			First(&tutorSubject).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.NotFound)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		var hour models.TutorHour
		if err := database.Get().
			Where("id = ?", hourId).
			Where("tutor_subject_id = ?", tutorSubject.ID).
			Find(&hour).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.NotFound)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		// on autorise seulement les admins et le tuteur/tutoré à modifier l'heure
		if !user.IsAdmin && tutorSubject.TutorID != user.ID && hour.TuteeID != user.ID {
			_ = c.Error(apierrors.Forbidden)
			return
		}

		parsedStartDate, err := time.Parse("2006-01-02T15:04:05.000-07:00", input.StartDate)
		if err != nil {
			_ = c.Error(err)
			return
		}

		parsedEndDate, err := time.Parse("2006-01-02T15:04:05.000-07:00", input.EndDate)
		if err != nil {
			_ = c.Error(err)
			return
		}

		// on calculera un delta de temps pour mettre à jour le nombre d'heures,
		// on stocke donc la durée originale avant de la modifier
		originalDuration := hour.EndDate.Sub(hour.StartDate).Hours()

		hour.StartDate = parsedStartDate
		hour.EndDate = parsedEndDate

		if err = database.Get().Save(&hour).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		// on met à jour le total d'heures du tuteur et du tutoré,
		// on dispose déjà du tutorSubject, on va donc chercher le tutoré
		var tuteeReg models.TuteeRegistration
		if err = database.Get().
			Where("tutee_id = ?", hour.TuteeID).
			Where("tutor_subject_id = ?", tutorSubject.ID).
			First(&tuteeReg).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.NotFound)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		// calcule le delta de temps
		newDuration := hour.EndDate.Sub(hour.StartDate).Hours()
		durationDelta := newDuration - originalDuration

		// on met à jour le total d'heures du tutoré
		tuteeReg.TotalHours += durationDelta
		if err = database.Get().Save(&tuteeReg).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		// on met à jour le total d'heures du tuteur
		tutorSubject.TotalHours += durationDelta
		if err = database.Get().Save(&tutorSubject).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		c.JSON(http.StatusOK, hour)
	}
}
