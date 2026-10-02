package models

import (
	"strconv"
	"time"
)

type Campaign struct {
	ID uint `gorm:"primarykey" json:"id"`

	Semester   int    `json:"semester"`
	SchoolYear string `json:"schoolYear"`

	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`

	RegistrationStatus    string    `json:"registrationStatus"`
	RegistrationStartDate time.Time `json:"registrationStartDate"`
	RegistrationEndDate   time.Time `json:"registrationEndDate"`

	// une campagne archivée n'est plus proposée aux tuteurs/tutorés (tableaux de bord),
	// mais reste consultable et modifiable depuis l'admin : ce n'est pas une suppression
	IsArchived bool `json:"isArchived"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// SchoolYearStart renvoie l'année de rentrée de l'année scolaire de la campagne (du 1er septembre
// au 31 août), ex. 2025 pour une campagne qui commence entre le 01/09/2025 et le 31/08/2026
func (c Campaign) SchoolYearStart() int {
	start := c.StartDate
	// les dates sont envoyées à minuit heure de Paris : on évite qu'un 1er septembre devienne un 31 août en UTC
	if loc, err := time.LoadLocation("Europe/Paris"); err == nil {
		start = start.In(loc)
	}
	if start.Month() >= time.September {
		return start.Year()
	}
	return start.Year() - 1
}

// ComputeSchoolYear remplit SchoolYear (ex. "2025-2026") à partir de la date de début
func (c *Campaign) ComputeSchoolYear() {
	year := c.SchoolYearStart()
	c.SchoolYear = strconv.Itoa(year) + "-" + strconv.Itoa(year+1)
}

// IsRegistrationOpen indique si les tuteurs/tutorés peuvent encore s'inscrire ou modifier leurs disponibilités
func (c Campaign) IsRegistrationOpen(now time.Time) bool {
	if c.IsArchived || c.RegistrationStatus != "OPEN" {
		return false
	}
	if !c.RegistrationStartDate.IsZero() && now.Before(c.RegistrationStartDate) {
		return false
	}
	// la date de fin est saisie comme un jour (minuit) : on laisse les inscriptions ouvertes toute la journée
	if !c.RegistrationEndDate.IsZero() && !now.Before(c.RegistrationEndDate.AddDate(0, 0, 1)) {
		return false
	}
	return true
}
