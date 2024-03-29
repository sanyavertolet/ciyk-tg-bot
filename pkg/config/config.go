package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	TelegramBotToken        string `json:"telegramBotToken"`
	DatabaseFileName        string `json:"databaseFileName"`
	GoogleSheetsKeyFileName string `json:"googleSheetsKeyFileName"`
	GoogleSpreadsheetID     string `json:"googleSpreadsheetId"`
}

func LoadConfig(filePath string) (*Config, error) {
	config := Config{DatabaseFileName: "ciyk.db"}

	configFile, err := os.Open(filePath)
	if err != nil {
		return &config, err
	}

	defer func(configFile *os.File) { _ = configFile.Close() }(configFile)

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&config); err != nil {
		return &config, err
	}

	return &config, nil
}
