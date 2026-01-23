package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type configResponse struct {
	AuthMethod string `json:"authMethod"`
	CasUrl     string `json:"casUrl"`
	ServiceUrl string `json:"serviceUrl"`
}

func GetConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		authMethod := os.Getenv("AUTH_METHOD")
		if authMethod == "" {
			authMethod = "CAS"
		}

		c.JSON(http.StatusOK, configResponse{
			AuthMethod: authMethod,
			CasUrl:     os.Getenv("CAS_URL"),
			ServiceUrl: os.Getenv("SERVICE_URL"),
		})
	}
}
