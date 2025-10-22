package main

import (
	"log"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/romitou/insatutorat/core"
	"github.com/romitou/insatutorat/database"
	"github.com/romitou/insatutorat/middlewares"
	"github.com/romitou/insatutorat/routes/admin"
	adminCampaign "github.com/romitou/insatutorat/routes/admin/campaign"
	"github.com/romitou/insatutorat/routes/assignments"
	"github.com/romitou/insatutorat/routes/auth"
	"github.com/romitou/insatutorat/routes/campaign"
	"github.com/romitou/insatutorat/routes/campaign/agenda"
	"github.com/romitou/insatutorat/routes/campaign/availabilities"
	"github.com/romitou/insatutorat/routes/campaign/tutee"
	"github.com/romitou/insatutorat/routes/campaign/tutor"
	"github.com/romitou/insatutorat/routes/tutoring"
	"github.com/romitou/insatutorat/routes/tutoring/hours"
	"github.com/romitou/insatutorat/routes/tutoring/lessons"
	"gopkg.in/cas.v2"
)

func main() {
	// permets de charger le fichier .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	// connexion au client mail
	core.SetupMailer()
	// connexion à la base de données
	database.Connect()

	// middlewares étant utilisés dans certaines routes
	corsMiddleware := middlewares.CorsHandler()
	errorsMiddleware := middlewares.ErrorHandler()
	sessionMiddleware := middlewares.SessionHandler()
	userMiddleware := middlewares.UserHandler()
	adminMiddleware := middlewares.AdminHandler()

	casUrl, err := url.Parse(os.Getenv("CAS_URL"))
	if err != nil {
		log.Fatal("invalid CAS_URL: ", err)
		return
	}

	casClient := cas.NewClient(&cas.Options{
		URL: casUrl,
	})

	// définition du routeur principal
	router := gin.Default()

	// utilisation des middlewares généraux
	router.Use(errorsMiddleware)
	router.Use(corsMiddleware)
	router.Use(sessionMiddleware)

	// logique d'authentification
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/login", auth.Login())
		authRouter.POST("/send-link", auth.SendLink())
		authRouter.GET("/self", userMiddleware, auth.Self())
		authRouter.GET("/logout", auth.Logout())

		authRouter.POST("/validate", auth.Validate(*casClient))
	}

	// récapitulatifs des affectations (page principale)
	assignmentsRouter := router.Group("/assignments", userMiddleware)
	{
		assignmentsRouter.GET("/tutee", assignments.TuteeAssignments())
		assignmentsRouter.GET("/tutor", assignments.TutorAssignments())
	}

	// routes des campagnes de tutorat
	campaignRouter := router.Group("/campaign/:campaignId", userMiddleware)
	{
		campaignRouter.GET("/agenda", agenda.OverviewAgenda())

		campaignRouter.GET("/availabilities", availabilities.GetAvailabilities())
		campaignRouter.POST("/availabilities", availabilities.PostAvailabilities())

		campaignRouter.GET("/subjects", campaign.Subjects())

		tuteeRouter := campaignRouter.Group("/tutee")
		{
			tuteeRouter.GET("/registrations", tutee.GetRegistrations())
			tuteeRouter.POST("/registrations", tutee.PostRegistrations())
		}

		tutorRouter := campaignRouter.Group("/tutor")
		{
			tutorRouter.GET("/registrations", tutor.GetRegistrations())
			tutorRouter.POST("/registrations", tutor.PostRegistrations())
		}

	}

	// routes d'administration, muni du middleware admin (ordre important)
	adminRouter := router.Group("/admin", userMiddleware, adminMiddleware)
	{
		adminRouter.GET("/subjects", admin.GetSubjects())
		adminRouter.GET("/users", admin.GetUsers())

		adminRouter.GET("/campaigns", admin.GetCampaigns())
		adminRouter.POST("/campaigns", admin.PostCampaign())

		adminRouter.PATCH("/campaign/:campaignId", adminCampaign.PatchCampaign())
		acRouter := adminRouter.Group("/campaign/:campaignId")
		{
			acRouter.GET("/overview", adminCampaign.GetCampaign())
			acRouter.GET("/users", adminCampaign.GetUsers())

			acRouter.GET("/assignments", adminCampaign.GetAssignments())
			acRouter.POST("/assignments", adminCampaign.PostAssignments())

			acRouter.DELETE("/assignments/tutor", adminCampaign.DeleteTutorAssignment())
			acRouter.DELETE("/assignments/tutee", adminCampaign.DeleteTuteeAssignment())

			acRouter.GET("/generate-assignments", adminCampaign.GenerateAssignments())
		}
	}

	// routes pour la gestion des espaces tutorat
	tutRouter := router.Group("/tutoring/:tutorSubjectId", userMiddleware)
	{
		tutRouter.GET("/summary", tutoring.GetSummary())

		// séances
		tutRouter.POST("/lessons", lessons.PostLesson())
		tutRouter.PATCH("/lesson/:lessonId", lessons.PatchLesson())
		tutRouter.DELETE("/lesson/:lessonId", lessons.DeleteLesson())

		// heures
		tutRouter.POST("/hours", hours.PostHour())
		tutRouter.PATCH("/hour/:hourId", hours.PatchHour())
		tutRouter.DELETE("/hour/:hourId", hours.DeleteHour())
	}

	// démarrage du routeur, utilise le PORT défini dans les variables d'environnement
	err = router.Run()
	if err != nil {
		log.Fatal("error starting server: ", err)
		return
	}
}
