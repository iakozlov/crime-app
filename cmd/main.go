package main

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/iakozlov/crime-app"
	"github.com/iakozlov/crime-app/pkg/handler"
	"github.com/iakozlov/crime-app/pkg/repository"
	"github.com/iakozlov/crime-app/pkg/service"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing congigs: %s", err.Error())
	}
	db, err := repository.NewPostgresDB(repository.Config{
		Host: "localhost",
		Port: "5432",
		Username: "postgres",
		Password: "qwerty",
		DBName: "postgres",
		SSLMode: "disable",
	})
	if err != nil{
		log.Fatalf("failed to initializaed db: %s", err.Error())
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(crime.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occurred while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
