package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
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

	CasUsername string      `gorm:"uniqueIndex" json:"-"`
	FirstName   string      `json:"firstName"`
	LastName    string      `json:"lastName"`
	Mail        string      `gorm:"uniqueIndex" json:"-"`
	Groups      StringArray `json:"-"`
	StpiYear    int         `json:"-"` // année d'étude (1, 2)

	IsTutor bool `json:"-"`
	IsTutee bool `json:"-"`
	IsAdmin bool `json:"-"`

	// used for login links
	LoginToken       string    `json:"-"`
	LoginRequestedAt time.Time `json:"-"`

	// dernière connexion réussie (CAS ou magic link) ; distinct de UpdatedAt qui
	// bouge aussi lors d'une modification admin (rôles, anonymisation...)
	LastLoginAt time.Time `json:"-"`

	// non-nil une fois le compte anonymisé (départ présumé du cycle STPI) ; les
	// champs identifiants sont alors purgés mais la ligne est conservée pour
	// préserver l'intégrité des heures/séances/inscriptions historiques
	AnonymizedAt *time.Time `json:"-"`

	Availabilities []SemesterAvailability `json:"-"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (user User) IsEmpty() bool {
	return user.ID == 0
}

// Anonymize purge les champs identifiants d'un utilisateur qui a quitté le cycle
// STPI, tout en conservant la ligne (et donc les heures/séances/inscriptions
// liées) pour les statistiques. Opération non réversible.
func (user *User) Anonymize() {
	now := time.Now()

	user.CasUsername = fmt.Sprintf("anonymise-%d", user.ID)
	user.FirstName = "Utilisateur"
	user.LastName = "anonymisé"
	user.Mail = fmt.Sprintf("anonymise-%d@anonymise.local", user.ID)
	user.Groups = StringArray{}
	user.StpiYear = 0

	user.IsTutor = false
	user.IsTutee = false
	user.IsAdmin = false

	user.LoginToken = ""
	user.LoginRequestedAt = time.Time{}

	user.AnonymizedAt = &now
}

// IsEligibleForAnonymization détermine si un compte peut être considéré comme
// ayant quitté le cycle STPI (2 ans : STPI1 puis STPI2), en combinant la dernière
// connexion et la dernière année d'étude connue :
//   - un compte vu en STPI1 dispose de 2 ans pleins pour terminer STPI2 et en sortir
//   - un compte vu en STPI2 (ou dont l'année n'est pas renseignée, ex: scolarité
//     aménagée) n'a plus qu'1 an avant de quitter le cycle
//
// Les comptes déjà anonymisés ou administrateurs ne sont jamais éligibles : un
// admin peut être un membre du personnel sans lien avec le cycle STPI.
func (user User) IsEligibleForAnonymization(now time.Time, stpi1GraceMonths, stpi2GraceMonths int) bool {
	if user.IsAdmin || user.AnonymizedAt != nil || user.LastLoginAt.IsZero() {
		return false
	}

	graceMonths := stpi2GraceMonths
	if user.StpiYear == 1 {
		graceMonths = stpi1GraceMonths
	}

	return user.LastLoginAt.AddDate(0, graceMonths, 0).Before(now)
}

type PrivateUser struct {
	ID uint `gorm:"primarykey" json:"id"`

	CasUsername string      `gorm:"uniqueIndex" json:"casUsername"`
	FirstName   string      `json:"firstName"`
	LastName    string      `json:"lastName"`
	Mail        string      `gorm:"uniqueIndex" json:"mail"`
	StudyYear   int         `json:"studyYear"`
	Groups      StringArray `json:"groups"`

	IsTutor bool `json:"isTutor"`
	IsTutee bool `json:"isTutee"`
	IsAdmin bool `json:"isAdmin"`

	// used for login links
	// LoginToken       string    `json:"-"`
	// LoginRequestedAt time.Time `json:"-"`

	LastLoginAt  time.Time  `json:"lastLoginAt"`
	AnonymizedAt *time.Time `json:"anonymizedAt"`

	Availabilities []SemesterAvailability `json:"-"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (user User) ToPrivate() PrivateUser {
	return PrivateUser{
		ID:           user.ID,
		CasUsername:  user.CasUsername,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Mail:         user.Mail,
		StudyYear:    user.StpiYear,
		Groups:       user.Groups,
		IsTutor:      user.IsTutor,
		IsTutee:      user.IsTutee,
		IsAdmin:      user.IsAdmin,
		LastLoginAt:  user.LastLoginAt,
		AnonymizedAt: user.AnonymizedAt,
	}
}
