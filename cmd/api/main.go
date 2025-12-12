package main

import (
	//"fmt"
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

	env, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("cannot run the serv %v", err)
	}

	db := database.InitDB(env)

	userRepo := repository.NewUserRepository(db)

	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer nc.Close()

	messaging.LoadWorker(nc, userRepo)
	r := gin.Default()

	authHandler := handlers.NewAuthHandler(nc, env)
	userHandler := handlers.NewUserHandler(nc)

	r.GET("/login", authHandler.Login)
	r.GET("/callback", authHandler.CallBack)

	//ex: /users/pnaessen/admin || /users/pnaessen/instructor
	r.PATCH("/users/:username/role", userHandler.UpdateRole)
	//ex: /users/pnaessen/info || /users/cassie/info
	r.GET("/users/:username/info", userHandler.GetUserInfo)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("cannot run the serv %v", err)
	}
}
