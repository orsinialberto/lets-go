package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.api2kafka/model"
	"github.com/IBM/sarama"
)

func main() {
	scheduler()

	// for {
	// 	scheduler()
	// 	time.Sleep(5 * time.Second)
	// }
}

func scheduler() {
	resp, err := http.Get("http://127.0.0.1:8080/players")
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	var players []model.Player
	if err := json.NewDecoder(resp.Body).Decode(&players); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Found %d players to process\n", len(players))

	for _, player := range players {
		produce(&player)
	}
}

func produce(p *model.Player) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	player := *p
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		fmt.Println("Error while creating kafka producer:", err)
		return
	}
	defer producer.Close()

	value, err := player.ToJsonString()
	if err != nil {
		return
	}

	message := &sarama.ProducerMessage{
		Topic: "api_2_kafka",
		Value: sarama.StringEncoder(value),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Message successfully sent (partition: %d, offset: %d)\n", partition, offset)
	}
}
