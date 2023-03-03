package database

import (
	"github.com/JosseMontano/estateInTheCloud/models"
	"github.com/JosseMontano/estateInTheCloud/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	
 	host := utils.DotEnvVariable("HOST")
	user := utils.DotEnvVariable("USER")
	password := utils.DotEnvVariable("PASSWORD")
	dbname := utils.DotEnvVariable("DBNAME")
	port := utils.DotEnvVariable("PORT_DB")
 
	DSN := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port

	database, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})

	if err != nil {
		panic("Could not connect to the database")
	}

	DB = database

	database.AutoMigrate(models.User{})
	database.AutoMigrate(models.TypeRealEstate{})
	database.AutoMigrate(models.Photo{})
	database.AutoMigrate(models.RealEstate{})
	database.AutoMigrate(models.Question{})
}
