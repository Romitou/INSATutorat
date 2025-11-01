package availabilities

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/romitou/insatutorat/core"

	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
)

func PostAvailabilities() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)

		campaignId := c.Param("campaignId")
		if campaignId == "" {
			_ = c.Error(apierrors.BadRequest)
			return
		}

		var campaign models.Campaign
		if err := database.Get().
			Where("id = ?", campaignId).
			Where("school_year = ?", os.Getenv("SCHOOL_YEAR")).
			First(&campaign).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(apierrors.NotFound)
				return
			}
			apierrors.DatabaseError(c, err)
			return
		}

		var slotsJson models.Slots
		if err := c.ShouldBindJSON(&slotsJson); err != nil {
			_ = c.Error(err)
			return
		}

		// les slots sont de type map[time.Weekday][]int
		slots := make(models.Slots, 5)

		// se référer à core pour avoir le détail des constantes,
		// initialisation des cases à 0 et -1 selon le JSON
		for i := time.Monday; i <= time.Friday; i++ {
			for j := core.M1; j <= core.A4; j++ {
				slots[i] = append(slots[i], 0)
				if slotsJson[i][j] == -1 {
					slots[i][j] = -1
				}
			}
		}

		// on récupère l'agenda de l'utilisateur pour le semestre
		campaignOverview, err := core.GetCampaignOverview(os.Getenv("SCHOOL_YEAR")+"-STPI"+strconv.Itoa(user.StpiYear), campaign, user.Groups)
		if err != nil {
			_ = c.Error(err)
			return
		}

		// pour chaque jour de la semaine, on va remplir les slots
		for _, overviewDay := range campaignOverview {
			var day time.Weekday
			day, err = core.ParseWeekday(overviewDay.Day)
			if err != nil {
				_ = c.Error(err)
				return
			}
			// pour chaque slot de la journée, on va remplir les slots en sommant le nombre de cours durant le semestre
			for _, overviewSlot := range overviewDay.Periods {
				sumItems := 0
				for _, item := range overviewSlot.Items {
					sumItems += item
				}
				slots[day][overviewSlot.Period] = sumItems
			}
		}

		// on a donc un agenda rempli avec les disponibilités de l'utilisateur et dont les créneaux occupés par des cours
		// ont été "overwrite" par le nombre de cours durant le semestre, évitant des saisies invalides

		// on récupère le résultat en JSON
		availabilityJSON, err := json.Marshal(slots)
		if err != nil {
			_ = c.Error(err)
			return
		}

		var semesterAvailability models.SemesterAvailability
		if err = database.Get().
			Where("user_id = ?", user.ID).
			Where("campaign_id = ?", campaign.ID).
			Find(&semesterAvailability).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		// création de la disponibilité si elle n'existe pas
		if semesterAvailability.ID == 0 {
			semesterAvailability = models.SemesterAvailability{
				UserID:           user.ID,
				CampaignID:       campaign.ID,
				AvailabilityJSON: string(availabilityJSON),
			}
			if err = database.Get().Create(&semesterAvailability).Error; err != nil {
				apierrors.DatabaseError(c, err)
				return
			}
		} else {
			// sinon, on met à jour la disponibilité
			semesterAvailability.AvailabilityJSON = string(availabilityJSON)
			if err = database.Get().Save(&semesterAvailability).Error; err != nil {
				apierrors.DatabaseError(c, err)
				return
			}
		}

		c.Status(http.StatusOK)
	}
}
