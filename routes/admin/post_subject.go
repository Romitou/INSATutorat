package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
)

func PostSubject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.Subject
		if err := c.ShouldBindJSON(&input); err != nil {
			_ = c.Error(err)
			return
		}

		if err := database.Get().
			Create(&input).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		c.JSON(http.StatusOK, input)
	}
}
