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

// AnonymizeUser purge les champs identifiants d'un compte (nom, email, identifiant
// CAS...), typiquement utilisé lorsqu'un utilisateur a quitté le cycle STPI. La
// ligne est conservée pour préserver l'historique (heures, séances, inscriptions).
// Opération manuelle, déclenchée par un admin, et non réversible.
func AnonymizeUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdStr := c.Param("id")
		if userIdStr == "" {
			_ = c.Error(apierrors.BadRequest)
			return
		}
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			_ = c.Error(apierrors.BadRequest)
			return
		}

		var targetUser models.User
		if err = database.Get().
			Where("id = ?", userId).
			First(&targetUser).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.NotFound)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		if targetUser.AnonymizedAt != nil {
			_ = c.Error(apierrors.AlreadyAnonymized)
			return
		}

		// un compte admin peut être un membre du personnel sans lien avec le cycle
		// STPI : il doit d'abord être rétrogradé (Admin > Utilisateurs) avant de
		// pouvoir être anonymisé, pour éviter une erreur de manipulation
		if targetUser.IsAdmin {
			_ = c.Error(apierrors.Forbidden)
			return
		}

		targetUser.Anonymize()

		if err = database.Get().Save(&targetUser).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		c.JSON(http.StatusOK, targetUser.ToPrivate())
	}
}
