package tests

import (
	"L0_azat/internal/config"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func StartMsgSpam(cfg *config.Config, delay time.Duration) {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	subject := "L0-subject"

	for {
		msg := generateMsg()
		//jsonMsg, err := json.Marshal(msg)
		jsonMsg, err := json.MarshalIndent(msg, "", "\t")
		if err != nil {
			log.Println("Error encoding JSON:", err)
			continue
		}

		err = nc.Publish(subject, jsonMsg)
		if err != nil {
			log.Println("Error publishing message:", err)
		} else {
			fmt.Println("Message sent: \n", string(jsonMsg))
			//fmt.Println("Message sent")
		}

		time.Sleep(delay)
	}
}
