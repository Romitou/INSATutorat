package lessons

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
	"net/http"
)

func DeleteLesson() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)

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

		// implémenter ici des vérifications ?
		// étape intermédiaire laissée intentionnellement

		if err := database.Get().
			Where("id = ?", lesson.ID).
			Delete(&lesson).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		c.Status(http.StatusOK)
	}
}
