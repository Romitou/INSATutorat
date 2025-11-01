package auth

import (
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/gorm"
)

type ServiceResponse struct {
	XMLName               xml.Name               `xml:"serviceResponse"`
	AuthenticationSuccess *AuthenticationSuccess `xml:"authenticationSuccess"`
}

type AuthenticationSuccess struct {
	User       string      `xml:"user"`
	Attributes *Attributes `xml:"attributes"`
}

type Attributes struct {
	DisplayName       string   `xml:"displayName"`
	GivenName         string   `xml:"givenName"`
	Mail              string   `xml:"mail"`
	SupannAffectation []string `xml:"supannAffectation"`
	SN                string   `xml:"sn"`
}

func CreateUserFromCas(serviceResp ServiceResponse) (*models.User, error) {
	var newUser models.User
	newUser.CasUsername = serviceResp.AuthenticationSuccess.User
	newUser.FirstName = serviceResp.AuthenticationSuccess.Attributes.GivenName
	newUser.LastName = serviceResp.AuthenticationSuccess.Attributes.SN
	newUser.Mail = serviceResp.AuthenticationSuccess.Attributes.Mail
	newUser.StpiYear = 0
	var stpiGroups []string
	stpiGroups = []string{}
	for _, affil := range serviceResp.AuthenticationSuccess.Attributes.SupannAffectation {
		if strings.Contains(affil, "stpi") {
			stpiGroups = append(stpiGroups, affil)
			if affil == "stpi1" {
				newUser.StpiYear = 1
				newUser.IsTutee = true
			} else if affil == "stpi2" {
				newUser.StpiYear = 2
				newUser.IsTutee = true
			}
			if strings.Contains(affil, "sa2") || strings.Contains(affil, "sa3") { // scolarité aménagée
				newUser.IsTutor = true
				newUser.IsTutee = true
			}
		}
	}

	newUser.Groups = stpiGroups
	return &newUser, nil
}

func Validate() gin.HandlerFunc {
	type query struct {
		Ticket string `form:"ticket" binding:"required"`
	}

	casUrl, parseErr := url.Parse(os.Getenv("CAS_URL"))
	if parseErr != nil {
		log.Fatal("invalid CAS_URL: ", parseErr)
	}

	serviceUrl, parseErr := url.Parse(os.Getenv("SERVICE_URL"))
	if parseErr != nil {
		log.Fatal("invalid SERVICE_URL: ", parseErr)
	}

	return func(c *gin.Context) {
		var q query
		if err := c.ShouldBindQuery(&q); err != nil || q.Ticket == "" {
			_ = c.Error(apierrors.Unauthorized)
			return
		}

		validateURL := casUrl.ResolveReference(&url.URL{
			Path:     "/cas/serviceValidate",
			RawQuery: url.Values{"service": {serviceUrl.String()}, "ticket": {q.Ticket}}.Encode(),
		}).String()

		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, validateURL, nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		httpClient := &http.Client{Timeout: 5 * time.Second}
		resp, err := httpClient.Do(req)
		if err != nil {
			_ = c.Error(err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			_ = c.Error(apierrors.Unauthorized)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			_ = c.Error(err)
			return
		}

		log.Println("CAS response body:")
		log.Println(string(body))

		var serviceResp ServiceResponse
		if err = xml.Unmarshal(body, &serviceResp); err != nil {
			_ = c.Error(err)
			return
		}

		if serviceResp.AuthenticationSuccess == nil {
			_ = c.Error(apierrors.Unauthorized)
			return
		}

		// CAS ticket is valid

		var existingUser models.User
		result := database.Get().Where(&models.User{
			CasUsername: serviceResp.AuthenticationSuccess.User,
		}).First(&existingUser)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				var newUser *models.User
				newUser, err = CreateUserFromCas(serviceResp)
				if err != nil {
					_ = c.Error(err)
					return
				}

				result = database.Get().Create(newUser)
				if result.Error != nil {
					apierrors.DatabaseError(c, result.Error)
					return
				}

				// on met à jour la session
				session := sessions.Default(c)
				session.Set("user_id", existingUser.ID)
				err = session.Save()
				if err != nil {
					_ = c.Error(err)
				}

				c.Status(http.StatusCreated)
				return
			}
			apierrors.DatabaseError(c, result.Error)
			return
		}

		// User exists, update info
		updatedUser, err := CreateUserFromCas(serviceResp)
		if err != nil {
			_ = c.Error(err)
			return
		}

		existingUser.StpiYear = updatedUser.StpiYear
		existingUser.Groups = updatedUser.Groups
		existingUser.IsTutee = updatedUser.IsTutee
		existingUser.IsTutor = updatedUser.IsTutor

		result = database.Get().Save(&existingUser)
		if result.Error != nil {
			apierrors.DatabaseError(c, result.Error)
			return
		}

		// on met à jour la session
		session := sessions.Default(c)
		session.Set("user_id", existingUser.ID)
		err = session.Save()
		if err != nil {
			_ = c.Error(err)
		}

		c.Status(http.StatusOK)
	}
}
