package admin

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
)

// délais d'inactivité (en mois) au-delà desquels un compte est considéré comme
// ayant quitté le cycle STPI, en fonction de la dernière année d'étude connue.
// Configurables via env pour ne pas figer une durée dans le code.
func anonymizationGraceMonths() (stpi1, stpi2 int) {
	stpi1 = 24
	stpi2 = 12

	if v, err := strconv.Atoi(os.Getenv("ANONYMIZATION_INACTIVITY_MONTHS_STPI1")); err == nil {
		stpi1 = v
	}
	if v, err := strconv.Atoi(os.Getenv("ANONYMIZATION_INACTIVITY_MONTHS_STPI2")); err == nil {
		stpi2 = v
	}

	return stpi1, stpi2
}

type AdminUserView struct {
	models.PrivateUser
	EligibleForAnonymization bool `json:"eligibleForAnonymization"`
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		if err := database.Get().
			Find(&users).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.NotFound)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		stpi1GraceMonths, stpi2GraceMonths := anonymizationGraceMonths()
		now := time.Now()

		// on est sur une route admin, on inclut les détails
		var views []AdminUserView
		for _, user := range users {
			views = append(views, AdminUserView{
				PrivateUser:              user.ToPrivate(),
				EligibleForAnonymization: user.IsEligibleForAnonymization(now, stpi1GraceMonths, stpi2GraceMonths),
			})
		}

		c.JSON(http.StatusOK, views)
	}
}
