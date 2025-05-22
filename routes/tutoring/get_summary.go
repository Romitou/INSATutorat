package tutoring

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
	"net/http"
)

type tuteeWithHours struct {
	models.User
	Hours []models.TutorHour `json:"hours"`
}

type summary struct {
	Subject models.Subject       `json:"subject"`
	Tutor   models.User          `json:"tutor"`
	Lessons []models.TutorLesson `json:"lessons"`
	Tutees  []tuteeWithHours     `json:"tutees"`
}

func GetSummary() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)

		tutorSubjectId := c.Param("tutorSubjectId")
		if tutorSubjectId == "" {
			_ = c.Error(apierrors.BadRequest)
			return
		}

		var tutorSubject models.TutorSubject
		if err := database.Get().
			Where("id = ?", tutorSubjectId).
			Preload("Tutees").
			Preload("Tutees.Tutee").
			Preload("Subject").
			Preload("Tutor").
			First(&tutorSubject).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.NotFound)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		// on autorise uniquement le tuteur, les tutorés et les admins
		if !user.IsAdmin && tutorSubject.TutorID != user.ID {
			isTutee := false
			for _, tuteeReg := range tutorSubject.Tutees {
				if tuteeReg.TuteeID == user.ID {
					isTutee = true
					break
				}
			}
			if !isTutee {
				_ = c.Error(apierrors.Forbidden)
				return
			}
		}

		var tutorHours []models.TutorHour
		query := database.Get().
			Where("tutor_subject_id = ?", tutorSubject.ID)

		// si l'utilisateur n'est pas admin et n'est pas le tuteur, on filtre par tutee_id
		// sinon, on récupère toutes les heures
		if !user.IsAdmin && tutorSubject.TutorID != user.ID {
			query = query.Where("tutee_id = ?", user.ID)
		}

		if err := query.Find(&tutorHours).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, []models.TutorHour{})
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		// on construit la liste des tutorés avec leurs heures
		tuteesWithHours := make([]tuteeWithHours, 0)
		for _, tutee := range tutorSubject.Tutees {
			tuteeHours := make([]models.TutorHour, 0)
			for _, hour := range tutorHours {
				if hour.TuteeID == tutee.TuteeID {
					tuteeHours = append(tuteeHours, hour)
				}
			}
			tuteesWithHours = append(tuteesWithHours, tuteeWithHours{
				User:  tutee.Tutee,
				Hours: tuteeHours,
			})
		}

		var lessons []models.TutorLesson
		if err := database.Get().
			Where("tutor_subject_id = ?", tutorSubject.ID).
			Find(&lessons).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, []models.TutorLesson{})
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		lessonsWithDetails := summary{
			Subject: tutorSubject.Subject,
			Tutor:   tutorSubject.Tutor,
			Lessons: lessons,
			Tutees:  tuteesWithHours,
		}

		c.JSON(http.StatusOK, lessonsWithDetails)
	}
}
