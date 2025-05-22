package campaign

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/core"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
)

// tuteeRegistrationWithAvailability complète les données de l'inscription d'un tutoré avec ses disponibilités
type tuteeRegistrationWithAvailability struct {
	Registration models.TuteeRegistration
	Availability models.Slots
}

// tutorSubjectWithAvailability complète les données de l'inscription d'un tuteur avec ses disponibilités ET le quota restant
type tutorSubjectWithAvailability struct {
	TutorSubject models.TutorSubject
	Availability models.Slots
	QuotaLeft    int
}

// emptySlots crée une structure de slots vide (= indispo partout) pour les disponibilités
func emptySlots() models.Slots {
	slots := make(models.Slots, 6)
	for i := time.Monday; i <= time.Friday; i++ {
		for j := core.M1; j <= core.A4; j++ {
			slots[i] = append(slots[i], -1)
		}
	}
	return slots
}

func GenerateAssignments() gin.HandlerFunc {
	return func(c *gin.Context) {
		var logs []string // logs de la génération
		var tuteeRegs []models.TuteeRegistration
		var tutorRegs []models.TutorSubject
		var subjects []models.Subject
		var availabilities []models.SemesterAvailability

		db := database.Get()
		campaignId := c.Param("campaignId")
		if campaignId == "" {
			_ = c.Error(apierrors.BadRequest)
			return
		}

		logs = append(logs, "Début de la génération pour la campagne "+campaignId)
		start := time.Now()
		logs = append(logs, "Date : "+start.Format("2006-01-02 15:04:05"))

		// récupération des inscriptions des tutorés
		if err := db.Where("campaign_id = ?", campaignId).
			Where("tutor_subject_id IS NULL").
			Preload("Tutee").
			Find(&tuteeRegs).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		// récupération des inscriptions des tuteurs
		if err := db.Where("campaign_id = ?", campaignId).
			Preload("Tutor").
			Preload("Tutees").
			Find(&tutorRegs).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		// récupération des matières du semestre
		if err := db.Find(&subjects).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		// récupération des disponibilités des tutorés et tuteurs
		if err := db.Where("campaign_id = ?", campaignId).
			Find(&availabilities).Error; err != nil {
			apierrors.DatabaseError(c, err)
			return
		}

		logs = append(logs, fmt.Sprintf("%d tutorés, %d tuteurs, %d matières, %d disponibilités récupérées", len(tuteeRegs), len(tutorRegs), len(subjects), len(availabilities)))

		// on crée un mapping direct entre l'id de la matière et la matière
		subjectMap := make(map[uint]models.Subject)
		for _, s := range subjects {
			subjectMap[s.ID] = s
		}

		// on crée les inscriptions des tutorés complétés avec leurs disponibilités
		tuteeParsed := make([]*tuteeRegistrationWithAvailability, len(tuteeRegs))
		for i, t := range tuteeRegs {
			availability := findAvailability(availabilities, t.TuteeID)
			t.Subject = subjectMap[t.SubjectID]
			tuteeParsed[i] = &tuteeRegistrationWithAvailability{Registration: t, Availability: availability}
		}

		// on crée les inscriptions des tuteurs complétés avec leurs disponibilités
		tutorParsed := make([]*tutorSubjectWithAvailability, len(tutorRegs))
		for i, t := range tutorRegs {
			availability := findAvailability(availabilities, t.TutorID)
			t.Subject = subjectMap[t.SubjectID]
			tutorParsed[i] = &tutorSubjectWithAvailability{
				TutorSubject: t,
				Availability: availability,
				QuotaLeft:    t.MaxTutees - len(t.Tutees), // on enlève le nombre de tutorés déjà affectés
			}
		}

		// on lance l'appariement. notons le passage de pointeurs pour éviter de faire des copies
		matchLogs := runMatching(tuteeParsed, tutorParsed, subjects)
		logs = append(logs, matchLogs...)

		// parmi les tutorés, on ne garde que ceux qui ont été affectés
		assigned := make([]models.TuteeRegistration, 0)
		for _, t := range tuteeParsed {
			if t.Registration.TutorSubjectID != nil {
				assigned = append(assigned, t.Registration)
			}
		}

		logs = append(logs, "")
		logs = append(logs, fmt.Sprintf("Total des affectations réussies : %d", len(assigned)))
		logs = append(logs, fmt.Sprintf("Durée de la génération : %s", time.Since(start).String()))
		c.JSON(200, gin.H{"affectedTutees": assigned, "logs": logs})
	}
}

// findAvailability recherche les disponibilités d'un utilisateur dans la liste des disponibilités
func findAvailability(avails []models.SemesterAvailability, userID uint) models.Slots {
	for _, a := range avails {
		if a.UserID == userID {
			var slots models.Slots
			if err := json.Unmarshal([]byte(a.AvailabilityJSON), &slots); err == nil {
				return slots
			}
		}
	}
	// si pas de disponibilité trouvée, on renvoie une structure vide, mais ne devrait pas arriver...
	return emptySlots()
}

// tuteeSlot est une STRUCTURE intermédiaire pour le matching. ne comporte qu'une référence vers l'inscription du tutoré,
// mais peut être amené à être modifié pour ajouter des données
type tuteeSlot struct {
	OriginalRegistration *tuteeRegistrationWithAvailability // nil si pas d'affectation (dummy)
}

// tutorSlot est une STRUCTURE intermédiaire pour le matching. ne comporte qu'une référence vers l'inscription du tuteur,
// mais peut être amené à être modifié pour ajouter des données
type tutorSlot struct {
	OriginalTutorSubject *tutorSubjectWithAvailability
}

// runMatching effectue l'appariement entre les tutorés et les tuteurs
func runMatching(tutees []*tuteeRegistrationWithAvailability, tutors []*tutorSubjectWithAvailability, subjects []models.Subject) []string {
	var logs []string

	// on parcourt les matières et on effectue l'appariement pour chaque matière
	for _, subject := range subjects {
		logs = append(logs, "")
		logs = append(logs, fmt.Sprintf("%s - %s", subject.ShortName, subject.Name))

		// on récupère les tutorés de cette matière
		var tuteesSlots []*tuteeSlot
		for _, t := range tutees {
			if t.Registration.Subject.ID == subject.ID && t.Registration.TutorSubjectID == nil {
				tuteesSlots = append(tuteesSlots, &tuteeSlot{OriginalRegistration: t})
			}
		}

		// on récupère les tuteurs de cette matière
		var tutorsSubj []*tutorSubjectWithAvailability
		for _, t := range tutors {
			if t.TutorSubject.Subject.ID == subject.ID {
				tutorsSubj = append(tutorsSubj, t)
			}
		}

		// grâce aux quotas, on crée une liste de slots pour les tuteurs
		tutorsSlots := []*tutorSlot{}
		for _, tutor := range tutorsSubj {
			for i := 0; i < tutor.QuotaLeft; i++ {
				tutorsSlots = append(tutorsSlots, &tutorSlot{OriginalTutorSubject: tutor})
			}
		}

		tuteesSize := len(tuteesSlots)
		tutorsSize := len(tutorsSlots)

		logs = append(logs, fmt.Sprintf("→ %d places, %d demandes", tutorsSize, tuteesSize))

		if tuteesSize == 0 {
			logs = append(logs, "→ Toutes les demandes sont déjà satisfaites.")
			continue
		} else if tutorsSize == 0 {
			logs = append(logs, "× Aucune place disponible pour cette matière.")
			continue
		}

		// impossible d'affecter plus de tutorés que de places disponibles,
		// voir rapport de TIP
		if tuteesSize > tutorsSize {
			logs = append(logs, "× Nombre de demandes supérieur au nombre de places disponibles, ajustez les quotas.")
			continue
		}

		// on ajoute des tutorés fictifs pour compléter le nombre de places,
		// voir rapport de TIP
		if tuteesSize < tutorsSize {
			for i := 0; i < tutorsSize-tuteesSize; i++ {
				tuteesSlots = append(tuteesSlots, &tuteeSlot{OriginalRegistration: nil}) // fictif
			}
		}

		// on construit les matrices (= tableaux) de préférences
		tuteePref := buildPreferenceMatrix(tuteesSlots, tutorsSlots)
		tutorPref := buildPreferenceMatrix(tutorsSlots, tuteesSlots)

		// on effectue l'appariement avec l'algorithme de Gale-Shapley
		matching := core.GaleShapley(tuteePref, tutorPref)

		matchCount := 0
		for tuteeIdx, tutorIdx := range matching {
			if tutorIdx == -1 {
				continue
			}

			tutee := tuteesSlots[tuteeIdx]
			tutor := tutorsSlots[tutorIdx]

			// on ne garde que les appariements valides, donc pas de fictifs
			if tutee.OriginalRegistration == nil || tutor.OriginalTutorSubject == nil {
				continue
			}

			tutee.OriginalRegistration.Registration.TutorSubjectID = &tutor.OriginalTutorSubject.TutorSubject.ID
			tutor.OriginalTutorSubject.TutorSubject.Tutees = append(
				tutor.OriginalTutorSubject.TutorSubject.Tutees,
				tutee.OriginalRegistration.Registration)
			matchCount++
		}

		logs = append(logs, fmt.Sprintf("✔ %d appariements effectués", matchCount))
	}
	return logs
}

// [A, B any] est un type générique
// buildPreferenceMatrix construit une matrice de préférences à partir de deux tableaux génériques d'objets
// les types A et B peuvent être quelconques, mais dans ce contexte, ils sont soit *tuteeSlot soit *tutorSlot
// la matrice de sortie est un tableau de tableau d'entiers : chaque ligne représente un élément de `from`
// et contient les indices des éléments de `to`, triés par ordre décroissant selon le score de préférence
func buildPreferenceMatrix[A, B any](from []*A, to []*B) [][]int {
	matrix := make([][]int, len(from))
	// on parcourt chaque élément de `from` et on calcule les scores de préférence pour chaque élément de `to`
	for i, f := range from {
		var scored []struct {
			Index int     // indice de l'élément dans `to`
			Score float64 // score de préférence calculé
		}

		for j, t := range to {
			var score float64

			// type assertion dynamique pour gérer les deux types possibles,
			// comme on ne connait pas le type à l'avance
			switch a := any(f).(type) {
			case *tuteeSlot: // cas d'un tutoré
				b := any(t).(*tutorSlot)

				if a.OriginalRegistration == nil || b.OriginalTutorSubject == nil {
					score = -1.0 // score fictif
				} else {
					score = core.AvailabilityScore(a.OriginalRegistration.Availability, b.OriginalTutorSubject.Availability)
				}
			case *tutorSlot: // cas d'un tuteur
				b := any(t).(*tuteeSlot)

				if a.OriginalTutorSubject == nil || b.OriginalRegistration == nil {
					score = -1.0 // score fictif
				} else {
					score = core.AvailabilityScore(a.OriginalTutorSubject.Availability, b.OriginalRegistration.Availability)
				}
			}

			// on ajoute le score et l'indice à la liste
			scored = append(scored, struct {
				Index int
				Score float64
			}{j, score})
		}

		// on trie les préfrénces par score décroissant
		sort.SliceStable(scored, func(i, j int) bool {
			return scored[i].Score > scored[j].Score // tri décroissant
		})

		// on remplit la ligne de la matrice
		for _, s := range scored {
			matrix[i] = append(matrix[i], s.Index)
		}
	}

	return matrix
}
