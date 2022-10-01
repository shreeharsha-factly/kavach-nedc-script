package config

import (
	"fmt"
	"log"
	"time"

	"github.com/factly/x/loggerx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var ProdDB *gorm.DB
var LocalDB *gorm.DB

// SetupDB is database setup
func SetupKavachProdDB() {
	fmt.Println("connecting to prod kavach database ...")

	dbString := fmt.Sprint("host=", "localhost", " ",
		"user=", "postgres", " ",
		"password=", "password", " ",
		"dbname=", "dbname", " ",
		"port=", "5433", " ",
		"sslmode=", "disable")

	dialector := postgres.Open(dbString)

	var err error
	ProdDB, err = gorm.Open(dialector, &gorm.Config{
		Logger: loggerx.NewGormLogger(logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Info,
			Colorful:      true,
		})})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to prod kavach database ...")
}

// SetupDB is database setup
func LocalSetupDB() {
	fmt.Println("connecting to local kavach database ...")

	dbString := fmt.Sprint("host=", "localhost", " ",
		"user=", "postgres", " ",
		"password=", "password", " ",
		"dbname=", "dbname", " ",
		"port=", "5433", " ",
		"sslmode=", "disable")

	dialector := postgres.Open(dbString)

	var err error
	LocalDB, err = gorm.Open(dialector, &gorm.Config{
		Logger: loggerx.NewGormLogger(logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Info,
			Colorful:      true,
		}),
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to local kavach database ...")
}
