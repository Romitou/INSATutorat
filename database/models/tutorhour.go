package models

import "time"

type TutorHour struct {
	ID uint `gorm:"primarykey" json:"id"`

	TutorSubject   TutorSubject `json:"-"`
	TutorSubjectID uint         `json:"-"`

	Tutee   User `json:"-"`
	TuteeID uint `json:"-"`

	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
