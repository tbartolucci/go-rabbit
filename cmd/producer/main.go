package main

import (
	"encoding/json"

	"bitsbybit.com/queue-project/internal/common"
	"bitsbybit.com/queue-project/internal/common/message"

	"math/rand"
	"time"
)

func main() {

	const QueueName = "add"
	rabbit := common.CreateAmpqRabbitMQ(common.AMQP_URL)
	rabbit.GetChannel()
	rabbit.GetQueue(QueueName)
	defer rabbit.Close()

	rand.Seed(time.Now().UnixNano())

	addTask := message.AddTask{Number1: rand.Intn(999), Number2: rand.Intn(999)}
	body, err := json.Marshal(addTask)
	if err != nil {
		common.HandleError(err, "Error encoding JSON")
	}

	rabbit.Publish(QueueName, body)
}
