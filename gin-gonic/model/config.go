package model

type Config struct {
	Database Database `json:"database"`
}

type Database struct {
	FileName string `json:"fileName"`
}
