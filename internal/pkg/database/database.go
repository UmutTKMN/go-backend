package database

import (
	"fmt"
	"log"

	"github.com/UmutTKMN/gobackend/configs"
	"github.com/UmutTKMN/gobackend/internal/app/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB paylaşılan veritabanı bağlantısı
var DB *gorm.DB

// Init veritabanı bağlantısını başlatır ve modelleri migrate eder
func Init(config *configs.Config) {
	var err error

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Password,
		config.DB.DBName,
		config.DB.SSLMode,
	)

	// Geliştirme modunda detaylı günlükleme
	gormConfig := &gorm.Config{}
	if config.AppEnv == "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("Veritabanına bağlanılamadı: %v", err)
	}

	log.Println("Veritabanına başarıyla bağlandı")

	// Tabloları otomatik migrate et
	err = DB.AutoMigrate(
		&model.User{},
		&model.UserRequest{},
		&model.UserDocument{},
		&model.UserDevice{},
		&model.UserPreference{},
		&model.UserCommunication{},
		&model.Role{},
		&model.Staff{},
	)
	if err != nil {
		log.Fatalf("Tabloları migrate ederken hata oluştu: %v", err)
	}

	log.Println("Veritabanı tabloları başarıyla oluşturuldu")
}
