package main

import (
	"github.com/iakozlov/crime-app"
	"github.com/iakozlov/crime-app/pkg/handler"
	"github.com/iakozlov/crime-app/pkg/repository"
	"github.com/iakozlov/crime-app/pkg/service"
	"log"
)

func main()  {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(crime.Server)
	if err := srv.Run("8000", handlers.InitRoutes()); err != nil{
		log.Fatalf("error occurred while running http server: %s", err.Error())
	}
}