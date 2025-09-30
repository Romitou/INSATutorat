package core

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

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

func SendLoginLink(user models.User, loginToken string) error {
	t, err := template.ParseFiles("mails/build_production/loginLink.html")
	if err != nil {
		return err
	}

	data := defaultData(user)
	data["link"] = os.Getenv("BASE_URL") + "/login?token=" + loginToken

	var htmlContent bytes.Buffer
	err = t.Execute(&htmlContent, data)
	if err != nil {
		return err
	}

	from := os.Getenv("MAIL_SENDER")
	to := user.Mail
	subject := "Tutorat INSA STPI - Lien de connexion"

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlContent.String())

	if smtpDialer == nil {
		return fmt.Errorf("SMTP dialer not initialized. Call SetupMailer() first.")
	}
	if err = smtpDialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
