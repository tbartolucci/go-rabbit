package main

import (
	"bitsbybit.com/queue-project/internal/common"
	"encoding/json"

	"math/rand"
	"time"
)

func main() {

	const QueueName = "add"
	rabbit := common.CreateAmpqRabbitMQ()
	rabbit.GetChannel()
	rabbit.GetQueue(QueueName)
	defer rabbit.Close()

	rand.Seed(time.Now().UnixNano())

	addTask := common.AddTask{Number1: rand.Intn(999), Number2 : rand.Intn(999)}
	body, err := json.Marshal(addTask)
	if err != nil {
		common.HandleError(err, "Error encoding JSON")
	}

	rabbit.Publish(QueueName,body)
}
