package database

import (
	"github.com/romitou/insatutorat/database/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var database *gorm.DB

func Connect() {
	// le dsn est documenté https://gorm.io/docs/connecting_to_the_database.html
	mysqlDsn := os.Getenv("MYSQL_DSN")

	logLevel := logger.Error
	switch os.Getenv("LOG_LEVEL") {
	case "silent":
		logLevel = logger.Silent
	case "debug":
		logLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(mysqlDsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Println(err)
	}

	// on migre les modèles automatiquement
	err = db.AutoMigrate(
		&models.Campaign{},
		&models.SemesterAvailability{},
		&models.Subject{},
		&models.TutorHour{},
		&models.TutorLesson{},
		&models.TutorSubject{},
		&models.User{},
		&models.TuteeRegistration{},
	)
	if err != nil {
		log.Println(err)
	}

	database = db
}

func Get() *gorm.DB {
	return database
}
