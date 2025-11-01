package auth

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
	"gopkg.in/cas.v2"
)

func Validate(client cas.Client) gin.HandlerFunc {
	type query struct {
		Ticket string `form:"ticket" binding:"required"`
	}

	return func(c *gin.Context) {
		var q query
		if err := c.ShouldBindQuery(&q); err != nil || q.Ticket == "" {
			_ = c.Error(apierrors.Unauthorized)
			return
		}

		validateURL, err := client.ServiceValidateUrlForRequest(q.Ticket, c.Request)
		if err != nil {
			_ = c.Error(err)
			return
		}

		log.Println("CAS validate URL:")
		log.Println(validateURL)

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

		c.Status(http.StatusOK)
	}
}
