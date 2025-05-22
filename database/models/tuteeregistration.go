package models

import "time"

type Slots = map[time.Weekday][]int

type TuteeRegistration struct {
	ID uint `gorm:"primarykey" json:"id"`

	Tutee   User `json:"tutee"`
	TuteeID uint `json:"tuteeId"`

	Campaign   Campaign `json:"-"`
	CampaignID uint     `json:"-"`

	Subject   Subject `json:"-"`
	SubjectID uint    `json:"subjectId"`

	TutorSubject   TutorSubject `json:"-"`
	TutorSubjectID *uint        `json:"tutorSubjectId"`

	TotalHours float64 `sql:"type:decimal(3,2);" json:"totalHours"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
