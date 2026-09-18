package campaign

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

// setArchived factorise l'archivage/désarchivage : une simple mise à jour de colonne
// (plutôt qu'un Updates(&struct)) pour éviter le piège des zero-values de GORM, qui
// ignorerait silencieusement un "false" lors du désarchivage.
func setArchived(c *gin.Context, archived bool) {
	campaignIdStr := c.Param("campaignId")
	if campaignIdStr == "" {
		_ = c.Error(apierrors.BadRequest)
		return
	}
	campaignId, err := strconv.Atoi(campaignIdStr)
	if err != nil {
		_ = c.Error(apierrors.BadRequest)
		return
	}

	db := database.Get()

	var camp models.Campaign
	if err = db.Where("id = ?", campaignId).First(&camp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Error(apierrors.NotFound)
			return
		}
		apierrors.DatabaseError(c, err)
		return
	}

	if err = db.Model(&models.Campaign{}).
		Where("id = ?", campaignId).
		Update("is_archived", archived).Error; err != nil {
		apierrors.DatabaseError(c, err)
		return
	}

	camp.IsArchived = archived
	c.JSON(http.StatusOK, camp)
}

func ArchiveCampaign() gin.HandlerFunc {
	return func(c *gin.Context) {
		setArchived(c, true)
	}
}

func UnarchiveCampaign() gin.HandlerFunc {
	return func(c *gin.Context) {
		setArchived(c, false)
	}
}
