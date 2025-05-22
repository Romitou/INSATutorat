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

type postLessonJson struct {
	StartDate string `json:"startDate" binding:"required"`
	EndDate   string `json:"endDate" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

func PostLesson() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)

		var input postLessonJson
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

		lesson := models.TutorLesson{
			TutorSubjectID: tutorSubject.ID,
			StartDate:      parsedStartDate,
			EndDate:        parsedEndDate,
			Content:        input.Content,
		}

		if err = database.Get().Create(&lesson).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		c.JSON(http.StatusOK, lesson)
	}
}
