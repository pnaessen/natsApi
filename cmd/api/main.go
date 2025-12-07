package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"natsApi/internal/config"
	"natsApi/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	if env, err := config.LoadEnv(); err != nil {
		fmt.Println(env)
		log.Fatalf("cannot run the serv %v", err)
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connect to NATS: %v", err)
	}
	defer nc.Close()

	r := gin.Default()

	authHandler := handlers.NewAuthHandler(nc)
	r.GET("/login", authHandler.Login)
	r.GET("/callback", authHandler.CallBack)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("cannot run the serv %v", err)
	}
}


	// _, err = nc.Subscribe("foo", func(m *nats.Msg) {
    // fmt.Printf("Reçu sur 'foo': %s\n", string(m.Data))
	// })

	// if err != nil {
	// 	log.Fatalf("Error Subscribe to NATS: %v", err)
	// }

	// err = nc.Publish("foo", []byte("Hello World"))
	// if err != nil {
	// 	log.Fatalf("Error publish NATS: %v", err)
	// }