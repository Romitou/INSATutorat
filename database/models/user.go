package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// si on utilisait une autre base de données type postgresql, on pourrait intégrer le type json directement
// sans passer par une couche applicative que voici :
type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}

	var strValue string

	switch v := value.(type) {
	case []byte: // peut parfois être représenté par un tableau d'octets
		strValue = string(v)
	case string:
		strValue = v
	default:
		return errors.New("type incompatible pour StringArray")
	}

	return json.Unmarshal([]byte(strValue), s)
}

type User struct {
	ID uint `gorm:"primarykey" json:"id"`

	FirstName  string      `json:"firstName"`
	LastName   string      `json:"lastName"`
	Mail       string      `gorm:"uniqueIndex" json:"-"`
	Groups     StringArray `json:"-"`
	StudyYear  int         `json:"-"` // année d'étude (1, 2, 3, 4, 5)
	SchoolYear string      `json:"-"` // année scolaire (2023, 2024, 2025, etc.), volontairement en string pour extensibilité

	IsTutor bool `json:"-"`
	IsTutee bool `json:"-"`
	IsAdmin bool `json:"-"`

	LoginToken       string    `json:"-"`
	LoginRequestedAt time.Time `json:"-"`

	Availabilities []SemesterAvailability `json:"-"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (user User) IsEmpty() bool {
	return user.ID == 0
}

type PrivateUser struct {
	ID uint `gorm:"primarykey" json:"id"`

	FirstName  string      `json:"firstName"`
	LastName   string      `json:"lastName"`
	Mail       string      `gorm:"uniqueIndex" json:"mail"`
	SchoolYear string      `json:"schoolYear"`
	Groups     StringArray `json:"groups"`

	IsTutor bool `json:"isTutor"`
	IsTutee bool `json:"isTutee"`
	IsAdmin bool `json:"isAdmin"`

	LoginToken       string    `json:"-"`
	LoginRequestedAt time.Time `json:"-"`

	Availabilities []SemesterAvailability `json:"-"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (user User) ToPrivate() PrivateUser {
	return PrivateUser{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Mail:       user.Mail,
		Groups:     user.Groups,
		SchoolYear: user.SchoolYear,
		IsTutor:    user.IsTutor,
		IsTutee:    user.IsTutee,
		IsAdmin:    user.IsAdmin,
	}
}
