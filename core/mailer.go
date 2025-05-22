package core

import (
	"bytes"
	"github.com/mailjet/mailjet-apiv3-go/v4"
	"github.com/romitou/insatutorat/database/models"
	"html/template"
	"os"
)

var client *mailjet.Client

// var emailStyle string

func SetupMailer() {
	client = mailjet.NewMailjetClient(os.Getenv("MAILJET_API_KEY"), os.Getenv("MAILJET_API_SECRET"))
	// si besoin d'injecter du CSS dans le mail, on peut le faire ici
	//file, err := os.ReadFile("mails/dist/index.css")
	//if err != nil {
	//	log.Println("Error reading CSS file:", err)
	//	return
	//}
	//emailStyle = string(file)
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
	data["link"] = os.Getenv("BASE_URL") + "/login#token=" + loginToken

	var htmlContent bytes.Buffer
	err = t.Execute(&htmlContent, data)
	if err != nil {
		return err
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: os.Getenv("MAIL_SENDER"),
				Name:  "InsaTutorat",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: user.Mail,
					Name:  user.LastName + " " + user.FirstName,
				},
			},
			Subject:  "Tutorat INSA STPI - Lien de connexion",
			HTMLPart: htmlContent.String(),
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err = client.SendMailV31(&messages)
	if err != nil {
		return err
	}
	return nil
}
