package messaging

import (
	"encoding/json"
	"log"
	"natsApi/internal/models"

	"github.com/nats-io/nats.go"
)

func LoadWorker(nc *nats.Conn) {

	var user models.UserMessage

	_, err := nc.Subscribe("user.login", func(m *nats.Msg) {
		if err := json.Unmarshal(m.Data, &user); err != nil {
			log.Fatalf("Error Subscrite Unmarshal %v", err)
		}
		//TODO insert dans la db
		resp := models.UserMessage{
			Username:   user.Username,
			Email:      user.Email,
			Role:       "admin",
			IntraID:    user.IntraID,
			SchoolYear: user.SchoolYear,
			IsActive:   false,
			Db_id:      1,
		}

		respBytes, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error marshal response: %v", err)
		}
		m.Respond(respBytes)
	})

	if err != nil {
		log.Fatalf("Error Subscribe to NATS: %v", err)
	}
}
