package middlewares

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/romitou/insatutorat/apierrors"
	"net/http"
)

// https://blog.depa.do/post/gin-validation-errors-handling
// https://github.com/gin-gonic/gin/issues/430

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		// si le contexte contient des erreurs, on les traite
		if len(ctx.Errors) > 0 {
			var eventIds []string
			for _, ctxErr := range ctx.Errors {

				// on regarde si ce sont des erreurs de validation
				var fieldsError validator.ValidationErrors
				ok := errors.As(ctxErr.Err, &fieldsError)
				if ok {
					validationJson := make(map[string]interface{})
					for _, fieldError := range fieldsError {
						validationJson[fieldError.Field()] = validationErrorToText(fieldError)
					}

					ctx.JSON(http.StatusBadRequest, gin.H{
						"errorCode": "VALIDATION_ERROR",
						"errors":    validationJson,
					})
					return
				}

				// on regarde si c'est une erreur de type json (par exemple, une erreur de parsing)
				var syntaxError *json.SyntaxError
				ok = errors.As(ctxErr.Err, &syntaxError)
				if ok {
					ctx.AbortWithStatusJSON(apierrors.InvalidJson.HttpCode, gin.H{
						"errorCode": apierrors.InvalidJson.ErrorCode,
						"help":      apierrors.InvalidJson.Help,
					})
					return
				}

				// on regarde si c'est une erreur "publique"/"commune"
				var publicError apierrors.PublicError
				ok = errors.As(ctxErr.Err, &publicError)
				if ok {
					ctx.AbortWithStatusJSON(publicError.HttpCode, gin.H{
						"errorCode": publicError.ErrorCode,
						"help":      publicError.Help,
					})
					return
				}

				// sinon, il s'agit d'une erreur interne
				var privateError apierrors.PrivateError
				ok = errors.As(ctxErr.Err, &privateError)
				if ok {
					eventIds = append(eventIds, privateError.EventID)
				}
			}

			if len(eventIds) > 0 {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"errorCode":   "INTERNAL_SERVER_ERROR",
					"identifiers": eventIds,
				})
			}

		}
	}
}

// validationErrorToText convertit une erreur de validation en texte lisible
func validationErrorToText(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fieldError.Field())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", fieldError.Field(), fieldError.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s", fieldError.Field(), fieldError.Param())
	case "email":
		return fmt.Sprintf("Invalid email format")
	case "len":
		return fmt.Sprintf("%s must be %s characters long", fieldError.Field(), fieldError.Param())
	}
	return fmt.Sprintf("%s is not valid", fieldError.Field())
}
