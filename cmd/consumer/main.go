package main

import (
	"bitsbybit.com/queue-project/internal/common"
	"encoding/json"
	"log"
	"os"
)

func main() {
	rabbit := common.CreateAmpqRabbitMQ()
	messageChannel := rabbit.GetMessageChannel("add")
	defer rabbit.Close()

	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChannel {
			log.Printf("Received a message: %s", d.Body)

			addTask := &common.AddTask{}
			err := json.Unmarshal(d.Body, addTask)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			log.Printf("%d + %d = %d", addTask.Number1, addTask.Number2, addTask.Number1+addTask.Number2)

			if err := d.Ack(false); err != nil {
				log.Printf("Error ACKing message: %s", err)
			} else {
				log.Printf("ACKed message")
			}
		}
	}()

	<-stopChan
}
