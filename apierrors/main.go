package apierrors

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/core"
)

type PrivateError struct {
	ErrorCode string `json:"error,omitempty"`
	EventID   string `json:"eventId,omitempty"`
}

func (e PrivateError) Error() string {
	return e.ErrorCode
}

type PublicError struct {
	HttpCode  int    `json:"code,omitempty"`
	ErrorCode string `json:"error,omitempty"`
	Help      string `json:"tip,omitempty"`
}

func (e PublicError) Error() string {
	return e.ErrorCode
}

func DatabaseError(ctx *gin.Context, err error) {
	_ = ctx.Error(PrivateError{
		ErrorCode: "DATABASE_ERROR",
		EventID:   LogError(ctx, err),
	})
}

func LogError(_ *gin.Context, err error) string {
	eventId := core.RandString(8)
	// permet de logger l'erreur dans un fichier avec un identifiant / corréler les erreurs
	// avec un identifiant également retourné au client
	// fait pour implémenter un système d'error catching type sentry
	fmt.Printf("Event ID: %s - Error: %s\n", eventId, err.Error())
	return eventId
}
