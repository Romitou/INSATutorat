package campaign

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func PatchCampaign() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		var input models.Campaign
		if err = c.ShouldBindJSON(&input); err != nil {
			_ = c.Error(err)
			return
		}

		var campaign models.Campaign
		if err = database.Get().
			Where("id = ?", campaignId).
			First(&campaign).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.NotFound)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		// on "sécurise" la mise à jour en ne mettant à jour que les champs nécessaires, càd on garde :
		// id, created_at, updated_at car on MàJ l'input directement dans la base de données
		input.ID = campaign.ID
		input.CreatedAt = campaign.CreatedAt
		input.UpdatedAt = campaign.UpdatedAt

		if err = database.Get().
			Where("id = ?", campaignId).
			Updates(&input).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		c.JSON(http.StatusOK, input)
	}
}
