package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"example.api2kafka/internal"
)

func init() {
	initVersionLookup()
}

func main() {
	go internal.Scheduler()
	select {}
}

func initVersionLookup() {
	file, err := os.Open(internal.VersionLookup)
	if os.IsNotExist(err) {
		file, err := os.Create(internal.VersionLookup)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		defer file.Close()

		if _, e := file.WriteString("0\n"); err != nil {
			fmt.Println("Error:", e)
			os.Exit(1)
		}
	}

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		internal.LastVersion, _ = strconv.Atoi(scanner.Text())
	} else if err := scanner.Err(); err != nil {
		fmt.Println("Something went wrong while reading version from file")
		os.Exit(1)
	}
}
