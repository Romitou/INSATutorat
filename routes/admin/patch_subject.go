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

func PatchSubject() gin.HandlerFunc {
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

		var input models.Subject
		if err = c.ShouldBindJSON(&input); err != nil {
			_ = c.Error(err)
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

		// on "sécurise" la mise à jour en ne mettant à jour que les champs nécessaires, càd on garde :
		// id, created_at, updated_at car on MàJ l'input directement dans la base de données
		input.ID = subject.ID
		input.CreatedAt = subject.CreatedAt
		input.UpdatedAt = subject.UpdatedAt

		if err = database.Get().
			Where("id = ?", subjectId).
			Updates(&input).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		c.JSON(http.StatusOK, input)
	}
}
