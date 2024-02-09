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
	if err := os.MkdirAll(internal.VersionLookupPath, os.ModePerm); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	_, err := os.Stat(internal.VersionLookupPath + internal.VersionLookupFilename)
	if os.IsNotExist(err) {
		if _, err := os.Create(internal.VersionLookupPath + internal.VersionLookupFilename); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}

	file, err := os.Open(internal.VersionLookupPath + internal.VersionLookupFilename)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		version := scanner.Text()
		if version == "" {
			internal.LastVersion = 0
		}
		internal.LastVersion, _ = strconv.Atoi(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Something went wrong while reading version from file")
		os.Exit(1)
	}
}
