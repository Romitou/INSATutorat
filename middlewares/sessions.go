package middlewares

import (
	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"github.com/romitou/insatutorat/database"
	"net/http"
	"os"
)

func SessionHandler() gin.HandlerFunc {
	store := gormsessions.NewStore(database.Get(), true, []byte(os.Getenv("SESSIONS_KEY")))

	opts := sessions.Options{
		Path:     "/",
		Domain:   os.Getenv("DOMAIN"),
		MaxAge:   60 * 60 * 24 * 90, // 3 months
		Secure:   true,
		HttpOnly: true,
	}

	if os.Getenv("DEV_MODE") == "true" {
		opts.Domain = ""
		opts.Secure = false
		opts.SameSite = http.SameSiteLaxMode
	}

	store.Options(opts)

	return sessions.Sessions("insa_tutorat_session", store)
}
