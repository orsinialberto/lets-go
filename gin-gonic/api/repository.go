package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"example.com/gin-gonic/model"
)

type People []model.Player

func (p People) Len() int           { return len(p) }
func (p People) Less(i, j int) bool { return p[i].Version < p[j].Version } // Sort by Version
func (p People) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func SavePlayer(player string) error {
	file, err := os.OpenFile(model.PlayersFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(player + "\n"); err != nil {
		fmt.Println("Error:", err)
		return err
	}

	return nil
}

func FindPlayerById(pId string) (model.Player, error) {
	var p model.Player
	file, err := os.OpenFile(model.PlayersFilePath, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Error:", err)
		return p, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if err := json.Unmarshal([]byte(line), &p); err != nil {
			fmt.Println("Error:", err)
			return p, err
		}
		if p.Id == pId {
			return p, nil
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
		return p, err
	}

	return p, nil
}

func ReadPlayersFromVersionLimit(from int, size int) ([]model.Player, error) {
	file, err := os.OpenFile(model.PlayersFilePath, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer file.Close()

	var players People
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		var p model.Player
		if err := json.Unmarshal([]byte(line), &p); err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		players = append(players, p)
	}

	sort.Sort(players)

	var result []model.Player
	var count int
	for _, p := range players {
		if count >= size {
			return result, nil
		}

		if p.Version >= from {
			result = append(result, p)
			count++
		}
	}

	return result, nil
}

func DeletePlayerById(pId string) error {

	filenameTmp, err := copyFileWithoutId(pId)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	if err := os.Remove(model.PlayersFilePath); err != nil {
		fmt.Println("Error:", err)
		return err
	}

	if err := os.Rename(filenameTmp, model.PlayersFilePath); err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}

func GetLastVersion() (int, error) {
	file, err := os.OpenFile(model.VersionFilePath, os.O_CREATE|os.O_RDWR|os.O_SYNC, 0666)
	if err != nil {
		fmt.Println("Error:", err)
		return 0, err
	}
	defer file.Close()

	var version int
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		version, err = strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Error:", err)
			return 0, err
		}
		fmt.Println("Last version is:", version)
	} else {
		fmt.Println("Empty file, version is 0")
		version = 0
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
		return 0, err
	}

	return version, nil
}

func UpdateVersion(version int) error {
	file1, err := os.OpenFile(model.VersionFilePath, os.O_WRONLY|os.O_TRUNC|os.O_SYNC, 0666)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	if _, err := file1.WriteString(strconv.Itoa(version) + "\n"); err != nil {
		fmt.Println("Error:", err)
		return err
	}

	defer file1.Close()

	return nil
}

func copyFileWithoutId(pId string) (string, error) {
	file, err := os.OpenFile(model.PlayersFilePath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	filenameTmp := model.PlayersFilePath + ".tmp"

	fileTmp, err := os.Create(filenameTmp)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	defer fileTmp.Close()

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, pId) {
			if _, err := fileTmp.WriteString(line + "\n"); err != nil {
				fmt.Println("Error:", err)
				return "", err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	return filenameTmp, nil
}
