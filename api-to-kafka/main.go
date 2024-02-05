package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"example.api2kafka/model"
	"github.com/IBM/sarama"
)

var last_version int

func init() {

	if _, err := os.Stat("data/version_lookup.txt"); os.IsNotExist(err) {
		file, err := os.Create("data/version_lookup.txt")
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		defer file.Close()

		if _, e := file.WriteString("0\n"); err != nil {
			fmt.Println("Error:", e)
			os.Exit(1)
		}
	} else if err == nil {
		file, err := os.OpenFile("data/version_lookup.txt", os.O_RDONLY, 0666)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			last_version, err = strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Something went wrong while reading version from file")
			os.Exit(1)
		}
	} else {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func main() {
	for {
		scheduler()
		time.Sleep(5 * time.Second)
	}
}

func scheduler() {
	url := "http://127.0.0.1:8080/players?size=10&from=" + strconv.Itoa(last_version+1)
	resp, err := http.Get(url)
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
		last_version = player.Version
	}

	if err := updateVersion(last_version); err != nil {
		fmt.Println("Error:", err)
		return
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

	fmt.Println("Processing player:", value)
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

func updateVersion(version int) error {
	file, err := os.OpenFile("data/version_lookup.txt", os.O_WRONLY|os.O_TRUNC|os.O_SYNC, 0666)
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
