package agenda

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/core"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func OverviewAgenda() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)

		campaignId := c.Param("campaignId")
		if campaignId == "" {
			_ = c.Error(apierrors.BadRequest)
			return
		}

		var campaign models.Campaign
		if err := database.Get().Where("id = ?", campaignId).First(&campaign).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.NotFound)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		// on récupère l'agenda de l'utilisateur pour le semestre
		campaignOverview, err := core.GetCampaignOverview(user.SchoolYear+"-STPI"+strconv.Itoa(user.StudyYear), campaign, user.Groups)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, campaignOverview)
	}
}
