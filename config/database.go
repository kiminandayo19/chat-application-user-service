package config

import (
	"chat-application/users/internal/domain/models"
	"log"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDbConnection() *gorm.DB {
	env, err := NewEnvDev()
	if err != nil {
		log.Print("[db] - Failed to load env")
	}
	dsn := "host=" + env.DB_HOST + " user=" + env.DB_USERNAME + " password=" + env.DB_PASS + " dbname=" + env.DB_NAME + " port=" + strconv.Itoa(env.DB_PORT) + " sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})

	if err != nil {
		log.Fatal("[db] - Postgres failed connection, reason:", err)
	} else {
		log.Println("[db] - Connected to postgres")
	}

	if err := db.AutoMigrate(&models.Mst_users{}); err != nil {
		log.Fatal("[db] - Migration failed:", err)
		return nil
	}
	log.Print("[db] - Migration success")
	return db
}
