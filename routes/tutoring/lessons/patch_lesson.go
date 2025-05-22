package lessons

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

type patchLessonJson struct {
	ID        uint   `json:"id" binding:"required"`
	StartDate string `json:"startDate" binding:"required"`
	EndDate   string `json:"endDate" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

func PatchLesson() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)

		var input patchLessonJson
		if err := c.ShouldBindJSON(&input); err != nil {
			_ = c.Error(err)
			return
		}

		tutorSubjectId := c.Param("tutorSubjectId")
		if tutorSubjectId == "" {
			_ = c.Error(apierrors.BadRequest)
			return
		}

		lessonId := c.Param("lessonId")
		if lessonId == "" {
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

		// on autorise uniquement le tuteur et les admins
		if !user.IsAdmin && tutorSubject.TutorID != user.ID {
			_ = c.Error(apierrors.Forbidden)
			return
		}

		var lesson models.TutorLesson
		if err := database.Get().
			Where("id = ?", lessonId).
			Where("tutor_subject_id = ?", tutorSubject.ID).
			Find(&lesson).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.NotFound)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		if lesson.ID == 0 {
			_ = c.Error(apierrors.BadRequest)
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

		// on "sécurise" la mise à jour en ne modifiant que les champs nécessaires
		lesson.StartDate = parsedStartDate
		lesson.EndDate = parsedEndDate
		lesson.Content = input.Content

		if err = database.Get().Save(&lesson).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		c.JSON(http.StatusOK, lesson)
	}
}
