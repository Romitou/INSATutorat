package models

import "time"

type TutorSubject struct {
	ID uint `gorm:"primarykey" json:"id"`

	Campaign   Campaign `json:"-"`
	CampaignID uint     `json:"-"`

	Subject   Subject `json:"-"`
	SubjectID uint    `json:"subjectId"`

	Tutor   User `json:"tutor"`
	TutorID uint `json:"-"`

	MaxTutees int                 `json:"maxTutees"`
	Tutees    []TuteeRegistration `json:"-"`

	TotalHours float64 `sql:"type:decimal(3,2);" json:"totalHours"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type TutorSubjectDetailed struct {
	ID uint `gorm:"primarykey" json:"id"`

	Campaign   Campaign `json:"-"`
	CampaignID uint     `json:"-"`

	Subject   Subject `json:"subject"`
	SubjectID uint    `json:"subjectId"`

	Tutor   User `json:"-"`
	TutorID uint `json:"-"`

	MaxTutees int                 `json:"maxTutees"`
	Tutees    []TuteeRegistration `json:"tutees"`

	TotalHours float64 `sql:"type:decimal(3,2);" json:"totalHours"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (t TutorSubject) ToDetailed() TutorSubjectDetailed {
	return TutorSubjectDetailed{
		ID:         t.ID,
		CampaignID: t.CampaignID,
		Subject:    t.Subject,
		SubjectID:  t.SubjectID,
		TutorID:    t.TutorID,
		MaxTutees:  t.MaxTutees,
		Tutees:     t.Tutees,
		TotalHours: t.TotalHours,
	}
}
