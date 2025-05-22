package models

import "time"

type SemesterAvailability struct {
	ID uint `gorm:"primarykey" json:"id"`

	Campaign   Campaign `json:"campaign"`
	CampaignID uint     `json:"campaignId"`

	User   User `json:"user"`
	UserID uint `json:"userId"`

	AvailabilityJSON string `json:"availabilityJson"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
