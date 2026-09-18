package admin

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
)

func DeleteSubject() gin.HandlerFunc {
	return func(c *gin.Context) {
		subjectIdStr := c.Param("id")
		if subjectIdStr == "" {
			_ = c.Error(apierrors.BadRequest)
			return
		}
		subjectId, err := strconv.Atoi(subjectIdStr)
		if err != nil {
			_ = c.Error(apierrors.BadRequest)
			return
		}

		var subject models.Subject
		if err = database.Get().
			Where("id = ?", subjectId).
			First(&subject).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.NotFound)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		// on interdit la suppression si la matière est déjà utilisée dans une campagne,
		// pour ne pas casser les données existantes (tutorat, inscriptions)
		var tutorSubjectCount int64
		if err = database.Get().Model(&models.TutorSubject{}).
			Where("subject_id = ?", subjectId).
			Count(&tutorSubjectCount).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		var tuteeRegistrationCount int64
		if err = database.Get().Model(&models.TuteeRegistration{}).
			Where("subject_id = ?", subjectId).
			Count(&tuteeRegistrationCount).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		if tutorSubjectCount > 0 || tuteeRegistrationCount > 0 {
			_ = c.Error(apierrors.ResourceInUse)
			return
		}

		if err = database.Get().Delete(&subject).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
