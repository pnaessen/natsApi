package main

import (
	"fmt"
	"log"

	"natsApi/internal/config"
	"natsApi/internal/database"
	"natsApi/internal/handlers"
	"natsApi/internal/messaging"
	repository "natsApi/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

func main() {

	if env, err := config.LoadEnv(); err != nil {
		fmt.Println(env)
		log.Fatalf("cannot run the serv %v", err)
	}

	db := database.InitDB()

	userRepo := repository.NewUserRepository(db)

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connect to NATS: %v", err)
	}
	defer nc.Close()

	messaging.LoadWorker(nc, userRepo)
	r := gin.Default()

	authHandler := handlers.NewAuthHandler(nc)
	r.GET("/login", authHandler.Login)
	r.GET("/callback", authHandler.CallBack)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("cannot run the serv %v", err)
	}
}
