package auth

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/apierrors"
)

func Validate() gin.HandlerFunc {
	type query struct {
		Ticket string `form:"ticket" binding:"required"`
	}

	casUrl, err := url.Parse(os.Getenv("CAS_URL"))
	if err != nil {
		log.Fatal("invalid CAS_URL: ", err)
	}

	serviceUrl, err := url.Parse(os.Getenv("SERVICE_URL"))
	if err != nil {
		log.Fatal("invalid SERVICE_URL: ", err)
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
