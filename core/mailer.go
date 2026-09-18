package core

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"

	"github.com/go-gomail/gomail"
	"github.com/romitou/insatutorat/database/models"
)

var smtpDialer *gomail.Dialer

func SetupMailer() {
	smtpHost := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	var smtpPort int
	_, err := fmt.Sscanf(port, "%d", &smtpPort)
	if err != nil {
		fmt.Println("Invalid SMTP_PORT, defaulting to 587")
		smtpPort = 587
	}
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	smtpDialer = gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
}

func defaultData(user models.User) map[string]interface{} {
	return map[string]interface{}{
		"user": user,
	}
}

// sendTemplatedMail factorise le chargement d'un template mails/build_production/*.html
// (généré par `npm run build` dans mails/, cf. mails/emails/) et l'envoi via SMTP.
func sendTemplatedMail(templateFile, to, subject string, data map[string]interface{}) error {
	t, err := template.ParseFiles("mails/build_production/" + templateFile)
	if err != nil {
		return err
	}

	var htmlContent bytes.Buffer
	if err = t.Execute(&htmlContent, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("MAIL_SENDER"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlContent.String())

	if smtpDialer == nil {
		return fmt.Errorf("SMTP dialer not initialized. Call SetupMailer() first.")
	}
	return smtpDialer.DialAndSend(m)
}

// tutoringSpaceLink construit le lien vers l'espace de tutorat (heures/séances) d'une
// affectation tuteur/tutoré donnée
func tutoringSpaceLink(tutorSubjectID uint) string {
	return os.Getenv("BASE_URL") + "/tutoring/" + strconv.FormatUint(uint64(tutorSubjectID), 10)
}

func SendLoginLink(user models.User, loginToken string) error {
	data := defaultData(user)
	data["link"] = os.Getenv("BASE_URL") + "/login?token=" + loginToken

	if os.Getenv("DEV_MODE") == "true" {
		log.Println("MAGIC LINK:", data["link"])
	}

	return sendTemplatedMail("loginLink.html", user.Mail, "Tutorat INSA STPI - Lien de connexion", data)
}

// SendAssignmentNotificationToTutee prévient un tutoré qu'il vient d'être mis en
// relation avec un tuteur pour une matière donnée
func SendAssignmentNotificationToTutee(tutee, tutor models.User, subject models.Subject, tutorSubjectID uint) error {
	data := map[string]interface{}{
		"tutee":   tutee,
		"tutor":   tutor,
		"subject": subject,
		"link":    tutoringSpaceLink(tutorSubjectID),
	}
	return sendTemplatedMail("assignmentNotificationTutee.html", tutee.Mail, "Tutorat INSA STPI - Vous avez un nouveau tuteur", data)
}

// SendAssignmentNotificationToTutor prévient un tuteur qu'un nouveau tutoré vient de
// lui être affecté pour une matière donnée
func SendAssignmentNotificationToTutor(tutor, tutee models.User, subject models.Subject, tutorSubjectID uint) error {
	data := map[string]interface{}{
		"tutor":   tutor,
		"tutee":   tutee,
		"subject": subject,
		"link":    tutoringSpaceLink(tutorSubjectID),
	}
	return sendTemplatedMail("assignmentNotificationTutor.html", tutor.Mail, "Tutorat INSA STPI - Un nouveau tutoré vous a été affecté", data)
}

// SendHourReminderToTutor prévient un tuteur qu'aucune heure n'a encore été déclarée
// pour une de ses matières. Seul le tutoré (ou un admin) peut déclarer les heures
// (cf. routes/tutoring/hours) : ce mail lui demande donc de relancer son tutoré,
// qui reçoit en parallèle son propre rappel via SendHourReminderToTutee.
func SendHourReminderToTutor(tutor models.User, subject models.Subject, tutorSubjectID uint) error {
	data := map[string]interface{}{
		"tutor":   tutor,
		"subject": subject,
		"link":    tutoringSpaceLink(tutorSubjectID),
	}
	return sendTemplatedMail("hourReminderTutor.html", tutor.Mail, "Tutorat INSA STPI - Relancez votre tutoré pour la déclaration des heures", data)
}

// SendHourReminderToTutee rappelle à un tutoré qu'il n'a encore déclaré aucune heure
// pour une de ses matières. C'est lui (ou un admin) qui doit renseigner les heures
// effectuées avec son tuteur.
func SendHourReminderToTutee(tutee, tutor models.User, subject models.Subject, tutorSubjectID uint) error {
	data := map[string]interface{}{
		"tutee":   tutee,
		"tutor":   tutor,
		"subject": subject,
		"link":    tutoringSpaceLink(tutorSubjectID),
	}
	return sendTemplatedMail("hourReminderTutee.html", tutee.Mail, "Tutorat INSA STPI - Pensez à déclarer vos heures", data)
}
