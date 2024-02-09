package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"example.api2kafka/model"
	"github.com/IBM/sarama"
)

const (
	VersionLookupPath     = "data/"
	VersionLookupFilename = "version_lookup.txt"
	playerUrl             = "http://127.0.0.1:8080/players?size=10&from="
	kafkaBroker           = "localhost:9092"
	topic                 = "player"
)

var LastVersion int

func Scheduler() {
	for {
		consume()
		time.Sleep(5 * time.Second)
	}
}

func consume() {
	resp, err := http.Get(playerUrl + strconv.Itoa(LastVersion+1))
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: %v\n", err)
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
		go produce(&player)
		LastVersion = player.Version
	}

	if err := updateVersion(LastVersion); err != nil {
		fmt.Println("Error:", err)
	}
}

func produce(p *model.Player) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	player := *p
	producer, err := sarama.NewSyncProducer([]string{kafkaBroker}, config)
	if err != nil {
		fmt.Println("Error while creating kafka producer:", err)
		return
	}
	defer producer.Close()

	value, err := player.ToJsonString()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Processing player:", value)
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(value),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Message successfully sent (partition: %d, offset: %d)\n", partition, offset)
	}
}

func updateVersion(version int) error {
	file, err := os.OpenFile(VersionLookupPath+VersionLookupFilename, os.O_WRONLY|os.O_TRUNC|os.O_SYNC, 0666)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	if _, err := file.WriteString(strconv.Itoa(version) + "\n"); err != nil {
		fmt.Println("Error:", err)
		return err
	}

	defer file.Close()

	return nil
}
