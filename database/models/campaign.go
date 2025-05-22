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

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
