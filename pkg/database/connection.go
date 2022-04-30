package database

import (
	"github.com/Reynadi531/creatica2022-be/pkg/utils"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func InitDatabase() *gorm.DB {
	dsn := utils.GenerateDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Error().Err(err).Msg("failed intialized database")
		return nil
	}

	DB = db
	log.Info().Msg("success established connection with database")

	return DB
}
