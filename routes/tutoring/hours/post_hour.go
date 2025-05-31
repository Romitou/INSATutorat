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

type postHourJson struct {
	TuteeId   uint   `json:"tuteeId" binding:"required"`
	StartDate string `json:"startDate" binding:"required"`
	EndDate   string `json:"endDate" binding:"required"`
}

func PostHour() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)

		var input postHourJson
		if err := c.ShouldBindJSON(&input); err != nil {
			_ = c.Error(err)
			return
		}

		tutorSubjectId := c.Param("tutorSubjectId")
		if tutorSubjectId == "" {
			_ = c.Error(apierrors.BadRequest)
			return
		}

		var tutorSubject models.TutorSubject
		if err := database.Get().
			Where("id = ?", tutorSubjectId).
			Preload("Tutees").
			First(&tutorSubject).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.NotFound)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		// on autorise uniquement le tuteur, le tutoré et les admins
		if !user.IsAdmin && tutorSubject.TutorID != user.ID {
			found := false
			for _, tuteeReg := range tutorSubject.Tutees {
				if tuteeReg.TuteeID == input.TuteeId && tuteeReg.TuteeID == user.ID {
					found = true
					break
				}
			}
			if !found {
				_ = c.Error(apierrors.Forbidden)
				return
			}
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

		hour := models.TutorHour{
			TutorSubjectID: tutorSubject.ID,
			TuteeID:        input.TuteeId,
			StartDate:      parsedStartDate,
			EndDate:        parsedEndDate,
		}

		if err = database.Get().Create(&hour).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		// on met à jour le total d'heures du tuteur et du tutoré,
		// on dispose déjà du tutorSubject, on va donc chercher le tutoré
		var tuteeReg models.TuteeRegistration
		if err = database.Get().
			Where("tutee_id = ?", input.TuteeId).
			Where("tutor_subject_id = ?", tutorSubject.ID).
			First(&tuteeReg).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.NotFound)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		// on met à jour le total d'heures du tutoré
		tuteeReg.TotalHours += hour.EndDate.Sub(hour.StartDate).Hours()
		if err = database.Get().Save(&tuteeReg).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		// on met à jour le total d'heures du tuteur
		tutorSubject.TotalHours += hour.EndDate.Sub(hour.StartDate).Hours()
		if err = database.Get().Save(&tutorSubject).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		c.JSON(http.StatusOK, hour)
	}
}
