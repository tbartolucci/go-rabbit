package main

import (
	"encoding/json"
	"log"
	"os"

	"bitsbybit.com/queue-project/internal/common"
	"bitsbybit.com/queue-project/internal/common/message"
)

func main() {
	rabbit := common.CreateAmpqRabbitMQ(common.AMQP_URL)
	messageChannel := rabbit.GetMessageChannel("add")
	defer rabbit.Close()

	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChannel {
			log.Printf("Received a message: %s", d.Body)

			addTask := &message.AddTask{}
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
