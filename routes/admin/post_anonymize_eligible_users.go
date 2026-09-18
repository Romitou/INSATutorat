package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
)

// AnonymizeEligibleUsers anonymise en une seule opération tous les comptes détectés
// comme ayant quitté le cycle STPI (cf. IsEligibleForAnonymization), plutôt que de
// les traiter un par un depuis l'interface. Opération manuelle, déclenchée par un
// admin, et non réversible.
func AnonymizeEligibleUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		if err := database.Get().
			Where("anonymized_at IS NULL AND is_admin = ?", false).
			Find(&users).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		stpi1GraceMonths, stpi2GraceMonths := anonymizationGraceMonths()
		now := time.Now()

		var anonymized []models.PrivateUser
		err := database.Get().Transaction(func(tx *gorm.DB) error {
			for _, user := range users {
				if !user.IsEligibleForAnonymization(now, stpi1GraceMonths, stpi2GraceMonths) {
					continue
				}

				user.Anonymize()
				if err := tx.Save(&user).Error; err != nil {
					return err
				}
				anonymized = append(anonymized, user.ToPrivate())
			}
			return nil
		})
		if err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"anonymizedCount": len(anonymized),
			"users":           anonymized,
		})
	}
}
