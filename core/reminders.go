package core

import (
	"log"
	"time"

	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
)

// délai minimum entre deux rappels envoyés au même tuteur, pour ne pas spammer
const hourReminderCooldown = 14 * 24 * time.Hour

// StartHourReminderScheduler lance une vérification périodique (toutes les 24h) des
// tuteurs actifs (campagne en cours, non archivée) n'ayant déclaré aucune heure, et
// leur envoie un rappel par email. À appeler une seule fois au démarrage du serveur.
func StartHourReminderScheduler() {
	go func() {
		// on laisse le temps au serveur de finir de démarrer avant la première vérification
		time.Sleep(1 * time.Minute)
		for {
			checkAndSendHourReminders()
			time.Sleep(24 * time.Hour)
		}
	}()
}

func checkAndSendHourReminders() {
	db := database.Get()
	now := time.Now()

	var tutorSubjects []models.TutorSubject
	if err := db.
		Preload("Tutor").
		Preload("Subject").
		Preload("Tutees.Tutee").
		Joins("JOIN campaigns ON campaigns.id = tutor_subjects.campaign_id").
		Where("campaigns.is_archived = ?", false).
		Where("campaigns.start_date <= ? AND campaigns.end_date >= ?", now, now).
		Where("tutor_subjects.hour_reminder_sent_at IS NULL OR tutor_subjects.hour_reminder_sent_at < ?", now.Add(-hourReminderCooldown)).
		Find(&tutorSubjects).Error; err != nil {
		log.Println("rappel d'heures : erreur lors de la récupération des tuteurs actifs :", err)
		return
	}

	for _, ts := range tutorSubjects {
		var hourCount int64
		if err := db.Model(&models.TutorHour{}).
			Where("tutor_subject_id = ?", ts.ID).
			Count(&hourCount).Error; err != nil {
			log.Println("rappel d'heures : erreur lors du comptage des heures :", err)
			continue
		}
		if hourCount > 0 {
			continue
		}

		// on relance le tuteur (qui n'a plus le droit de déclarer les heures lui-même,
		// cf. routes/tutoring/hours) et chaque tutoré qui, lui, le peut
		if err := SendHourReminderToTutor(ts.Tutor, ts.Subject, ts.ID); err != nil {
			log.Println("rappel d'heures : erreur lors de l'envoi de l'email au tuteur :", err)
			continue
		}

		for _, tuteeReg := range ts.Tutees {
			if err := SendHourReminderToTutee(tuteeReg.Tutee, ts.Tutor, ts.Subject, ts.ID); err != nil {
				log.Println("rappel d'heures : erreur lors de l'envoi de l'email au tutoré :", err)
			}
		}

		sentAt := time.Now()
		if err := db.Model(&models.TutorSubject{}).
			Where("id = ?", ts.ID).
			Update("hour_reminder_sent_at", &sentAt).Error; err != nil {
			log.Println("rappel d'heures : erreur lors de la mise à jour de hour_reminder_sent_at :", err)
		}
	}
}
