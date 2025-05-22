package core

import (
	"errors"
	"fmt"
	"github.com/mmcdole/gofeed/rss"
	"github.com/romitou/insatutorat/database/models"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"
)

type InsaPeriod int

const (
	M1 InsaPeriod = iota // 8h - 9h30
	M2                   // 9h45 - 11h15
	M3                   // 11h30 - 13h
	A1                   // 13h15 - 14h45
	A2                   // 15h - 16h30
	A3                   // 16h45 - 18h15
	A4                   // 18h30 - 20h
)

func GetStartEndDate(period InsaPeriod) (time.Time, time.Time) {
	var hour, minute int
	switch period {
	case M1:
		hour, minute = 8, 0
	case M2:
		hour, minute = 9, 45
	case M3:
		hour, minute = 11, 30
	case A1:
		hour, minute = 13, 15
	case A2:
		hour, minute = 15, 0
	case A3:
		hour, minute = 16, 45
	case A4:
		hour, minute = 18, 30
	}

	startTime := time.Date(0, 1, 1, hour, minute, 0, 0, time.UTC)
	endTime := startTime.Add(time.Hour + 30*time.Minute)

	return startTime, endTime
}

func GetInsaPeriods(startDate time.Time, endDate time.Time) []InsaPeriod {
	// on normalise startDate et endDate : on ne garde que l'heure et minute
	startHM := time.Date(0, 1, 1, startDate.Hour(), startDate.Minute(), 0, 0, time.UTC)
	endHM := time.Date(0, 1, 1, endDate.Hour(), endDate.Minute(), 0, 0, time.UTC)

	periods := make([]InsaPeriod, 0)
	for i := M1; i <= A4; i++ {
		start, end := GetStartEndDate(i)

		if (start.Equal(startHM) || start.After(startHM)) &&
			(end.Equal(endHM) || end.Before(endHM)) {
			periods = append(periods, i)
		}
	}
	return periods
}

var daysOfWeek = map[string]time.Weekday{
	"SUNDAY":    time.Sunday,
	"MONDAY":    time.Monday,
	"TUESDAY":   time.Tuesday,
	"WEDNESDAY": time.Wednesday,
	"THURSDAY":  time.Thursday,
	"FRIDAY":    time.Friday,
	"SATURDAY":  time.Saturday,
}

func ParseWeekday(v string) (time.Weekday, error) {
	if d, ok := daysOfWeek[v]; ok {
		return d, nil
	}

	return time.Sunday, fmt.Errorf("mauvias jour '%s'", v)
}

type AgendaItem struct {
	Title     string    `json:"title"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Groups    []string  `json:"groups"`
	Location  string    `json:"location"`
}

// les appels à l'api de l'insa sont coûteux, alors on met en place un cache
// pour éviter de surcharger le serveur avec des requêtes inutiles
var agendaCache = make(map[string]map[string][]AgendaItem)
var agendaCacheUpdates = make(map[string]map[string]time.Time)

func GetMonthAgenda(agenda string, date time.Time) ([]AgendaItem, error) {
	dateFormat := date.Format("20060102")
	if cachedDate, exists := agendaCacheUpdates[agenda]; exists {
		if dateCached, cacheExists := cachedDate[dateFormat]; cacheExists {
			if time.Since(dateCached) < time.Hour {
				return agendaCache[agenda][dateFormat], nil
			}
		}
	}

	request, err := http.NewRequest(
		"GET",
		"https://agendas.insa-rouen.fr/rss/rss2.0.php?cal="+agenda+"&cpath=&rssview=month&getdate="+dateFormat,
		nil)
	if err != nil {
		return nil, err
	}

	httpClient := http.DefaultClient
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("code de réponse non valide : " + response.Status)
	}

	rssParser := rss.Parser{}
	feed, err := rssParser.Parse(response.Body)
	if err != nil {
		return nil, err
	}

	var timeSlots []AgendaItem
	for _, item := range feed.Items {
		description := item.Description

		// sépare toutes les lignes (<br/> en html)
		parts := strings.Split(description, "<br/>")

		var groups []string
		for _, part := range parts {
			// on retire les espaces inutiles
			part = strings.TrimSpace(part)

			// on ne garde que les lignes qui commencent par "STPI" (STPI11, STPI12, etc.)
			if strings.HasPrefix(part, "STPI") {
				groups = append(groups, part)
			}
		}

		startString := item.Extensions["ev"]["startdate"][0].Value
		startDate, parseErr := time.Parse("2006-01-02T15:04:05", startString)
		if parseErr != nil {
			log.Printf("date invalide : %v", err)
			continue
		}

		endString := item.Extensions["ev"]["enddate"][0].Value
		endDate, dateErr := time.Parse("2006-01-02T15:04:05", endString)
		if dateErr != nil {
			log.Printf("date invalide : %v", err)
			continue
		}

		if item.Extensions["ev"]["location"] == nil {
			continue
		}

		location := item.Extensions["ev"]["location"][0].Value

		subject := strings.Split(item.Title, ": ")
		if len(subject) < 2 {
			log.Println("format invalide : '" + item.Title)
			continue
		}

		timeSlot := AgendaItem{
			Title:     subject[1],
			StartDate: startDate,
			EndDate:   endDate,
			Groups:    groups,
			Location:  location,
		}

		timeSlots = append(timeSlots, timeSlot)
	}

	// on crée le cache si inexistant
	if _, exists := agendaCache[agenda]; !exists {
		agendaCache[agenda] = make(map[string][]AgendaItem)
	}
	if _, exists := agendaCacheUpdates[agenda]; !exists {
		agendaCacheUpdates[agenda] = make(map[string]time.Time)
	}

	// on met le cache à jour
	agendaCache[agenda][dateFormat] = timeSlots
	agendaCacheUpdates[agenda][dateFormat] = time.Now()

	return timeSlots, nil
}

type OverviewDay struct {
	Day     string           `json:"day"`
	Periods []OverviewPeriod `json:"periods"`
}

type OverviewPeriod struct {
	Period int            `json:"period"`
	Items  map[string]int `json:"items"`
}

// generateMonthsBetween crée un tableau de chaque mois entre deux dates
func generateMonthsBetween(start, end time.Time) []time.Time {
	var months []time.Time
	for d := start; d.Before(end); d = d.AddDate(0, 1, 0) {
		months = append(months, d)
	}
	return months
}

// hasCommonGroup vérifie s'il y a une valeur commune entre deux listes
func hasCommonGroup(itemGroups, targetGroups []string) bool {
	for _, g := range targetGroups {
		if slices.Contains(itemGroups, g) {
			return true
		}
	}
	return false
}

func GetCampaignOverview(agenda string, campaign models.Campaign, groups []string) ([]OverviewDay, error) {
	start, end := campaign.StartDate, campaign.EndDate

	if start.IsZero() || end.IsZero() || start.After(end) {
		return nil, errors.New("dates de début/fin invalides")
	}

	months := generateMonthsBetween(start, end)
	agendaItems := make([]AgendaItem, 0)

	// Récupération et filtrage combinés
	for _, month := range months {
		// on récupère TOUS les items du mois
		items, err := GetMonthAgenda(agenda, month)
		if err != nil {
			return nil, err
		}

		// on filtre les items en fonction du groupe souhaité
		for _, item := range items {
			if hasCommonGroup(item.Groups, groups) {
				agendaItems = append(agendaItems, item)
			}
		}
	}

	// organisation par jour et période
	// structure assez complexe mais représentée comme suit :
	// chaque jour de la semaine (ex: Lundi) contient une map de périodes (ex: M1, M2, etc.)
	// chaque période contient une map de cours (ex: "CM-M8", "TD-I4", etc.) avec le nombre d'occurrences
	tempOverview := make(map[time.Weekday]map[InsaPeriod]map[string]int)

	for _, item := range agendaItems {
		weekday := item.StartDate.Weekday()
		periods := GetInsaPeriods(item.StartDate, item.EndDate)

		// si le jour n'existe pas, on l'initialise
		if _, ok := tempOverview[weekday]; !ok {
			tempOverview[weekday] = make(map[InsaPeriod]map[string]int)
		}

		// pour chaque période du cours, on l'ajoute à la map
		for _, period := range periods {
			// si la période n'existe pas, on l'initialise
			if _, ok := tempOverview[weekday][period]; !ok {
				tempOverview[weekday][period] = make(map[string]int)
			}
			tempOverview[weekday][period][item.Title]++
		}
	}

	// construction finale de l'overview
	overview := make([]OverviewDay, 0)
	for _, day := range []time.Weekday{
		time.Monday, time.Tuesday, time.Wednesday,
		time.Thursday, time.Friday, time.Saturday, time.Sunday,
	} { // on passe les jours dans l'ordre pour plus de clarté dans la sortie de l'api
		if periods, ok := tempOverview[day]; ok {
			dayPeriods := make([]OverviewPeriod, 0)
			for period, courseCounts := range periods {
				dayPeriods = append(dayPeriods, OverviewPeriod{
					Period: int(period),
					Items:  courseCounts,
				})
			}
			overview = append(overview, OverviewDay{
				Day:     strings.ToUpper(day.String()),
				Periods: dayPeriods,
			})
		}
	}

	return overview, nil
}
