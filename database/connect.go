package database

import (
	"github.com/JosseMontano/estateInTheCloud/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	host := "localhost"
	user := "postgres"
	password := "8021947cbba"
	dbname := "realEstatePrueba1"
	port := "5432"

	DSN := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port

	database, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})

	if err != nil {
		panic("Could not connect to the database")
	}

	DB = database

	database.AutoMigrate(models.User{})
	database.AutoMigrate(models.TypeRealState{})
	database.AutoMigrate(models.Photo{})
	database.AutoMigrate(models.RealEstate{})
}
