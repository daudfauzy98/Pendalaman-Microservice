package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/daudfauzy98/Pendalaman-Microservice/auth-service/database"
	"github.com/daudfauzy98/Pendalaman-Microservice/auth-service/handler"
	"github.com/daudfauzy98/Pendalaman-Microservice/menu-service/config"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		//log.Println(err.Error())
		log.Panic(err)
	} else {
		log.Println(cfg)
	}

	// Assignment variabel db dan err dengan nilai dari nilai kembalian
	// Method initDB(). *Harus 2 variabel yang dideklarasikan
	db, err := initDB(cfg.Database)
	authHandler := handler.AuthDB{Db: db}
	router := mux.NewRouter()

	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("DB connection success!")
	}

	router.Handle("/auth/validate", http.HandlerFunc(authHandler.ValidateAuth))
	router.Handle("/auth/signup", http.HandlerFunc(authHandler.SignUp))
	router.Handle("/auth/login", http.HandlerFunc(authHandler.Login))

	fmt.Println("Auth service listen on 8001")
	log.Panic(http.ListenAndServe(":8001", router))
}

// Ambil konfigurasi database yang akan digunakan dari file config.yml
func getConfig() (config.Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	viper.SetConfigName("config.yml")

	if err := viper.ReadInConfig(); err != nil {
		return config.Config{}, err
	}

	var cfg config.Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return config.Config{}, err
	}

	return cfg, nil
}

func initDB(cfg config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName, cfg.Config)
	log.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&database.Auth{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
