package campaign

import (
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"net/http"
	"strconv"
)

func DeleteTutorAssignment() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.Get()

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

		var input models.TutorSubject
		if err = c.ShouldBindJSON(&input); err != nil {
			_ = c.Error(err)
			return
		}

		if input.ID == 0 {
			_ = c.Error(apierrors.BadRequest)
			return
		}

		if err = db.
			Where("id = ? AND campaign_id = ?", input.ID, campaignId).
			Delete(&models.TutorSubject{}).Error; err != nil {
			apierrors.DatabaseError(c, err)
		}

		c.Status(http.StatusOK)
	}
}

func DeleteTuteeAssignment() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.Get()

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

		var input models.TuteeRegistration
		if err = c.ShouldBindJSON(&input); err != nil {
			_ = c.Error(err)
			return
		}

		if input.ID == 0 {
			_ = c.Error(apierrors.BadRequest)
			return
		}

		if err = db.
			Where("id = ? AND campaign_id = ?", input.ID, campaignId).
			Delete(&models.TuteeRegistration{}).Error; err != nil {
			apierrors.DatabaseError(c, err)
		}

		c.Status(http.StatusOK)
	}
}
