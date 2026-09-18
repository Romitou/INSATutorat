package models

import "time"

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
