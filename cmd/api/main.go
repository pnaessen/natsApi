package main

import (
	"fmt"
	"log"

	// "github.com/nats-io/nats.go"
	"natsApi/internal/config"
	"natsApi/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	if env, err := config.LoadEnv(); err != nil {
		fmt.Println(env)
		log.Fatalf("cannot run the serv %v", err)
	}

	r := gin.Default()

	authHandler := handlers.NewAuthHandler()
	r.GET("/login", authHandler.Login)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("cannot run the serv %v", err)
	}
}
