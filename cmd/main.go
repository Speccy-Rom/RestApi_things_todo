package main

import (
	"github.com/speccy-rom/RestApi_things_todo"
	"github.com/speccy-rom/RestApi_things_todo/pkg/handler"
	"github.com/speccy-rom/RestApi_things_todo/pkg/repository"
	"github.com/speccy-rom/RestApi_things_todo/pkg/service"
	"log"
)

func main() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	if err := srv.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
