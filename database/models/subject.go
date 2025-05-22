package models

import "time"

type Subject struct {
	ID uint `gorm:"primarykey" json:"id"`

	Semester  int    `json:"semester"`
	ShortName string `gorm:"uniqueIndex" json:"shortName"`
	Name      string `json:"name"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// https://agendas.insa-rouen.fr/rss/rss2.0.php?cal=2024-STPI1&cpath=&rssview=month&getdate=20250601
