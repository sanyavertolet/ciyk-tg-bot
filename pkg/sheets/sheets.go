package sheets

import (
	model "camus/sanyavertolet/bot/pkg/database/model"
	database "camus/sanyavertolet/bot/pkg/database/repository"
	"context"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"os"
)

const (
	fmtReadRange = "Schedule!A%d:F"
	dateFormat   = "15:04 02.01.2006"
)

type Sheets struct {
	Services      *sheets.Service
	SpreadsheetID string
}

func InitSheets(keyFileName string, spreadsheetID string) (*Sheets, error) {
	creds, err := os.ReadFile(keyFileName)
	if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
		return nil, err
	}

	config, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to create JWT config: %v", err)
		return nil, err
	}

	client := config.Client(context.Background())
	sheetsService, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Google Sheets service: %v", err)
		return nil, err
	}
	return &Sheets{sheetsService, spreadsheetID}, nil
}

func (sh *Sheets) SyncGames(repo *database.Repository) {
	log.Print("Synchronizing games")
	checkpoint, err := repo.GetLastCheckpoint()
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Last detected checkpoint: %d", checkpoint.Line)

	readRange := fmt.Sprintf(fmtReadRange, checkpoint.Line+1)
	response, err := sh.Services.Spreadsheets.Values.Get(sh.SpreadsheetID, readRange).Do()
	if err != nil {
		log.Panic(err)
	}

	var games []model.Game
	for _, row := range response.Values {
		dateString := fmt.Sprintf("%s %s", row[1], row[0])
		date, err := time.ParseInLocation(dateFormat, fmt.Sprintf("%s", dateString), time.Local)
		if err != nil {
			log.Panic(err)
		}

		maxPlayers := 9
		maxPlayers, err = strconv.Atoi(fmt.Sprintf("%s", row[5]))
		if err != nil {
			log.Printf("Could not parse maxPlayers, using 9 as default: %v", err)
			maxPlayers = 9
		}

		game := model.Game{
			Date:       date,
			Name:       fmt.Sprintf("%s %s", row[2], row[3]),
			Place:      fmt.Sprintf("%s", row[4]),
			MaxPlayers: maxPlayers,
		}

		games = append(games, game)
	}

	if len(response.Values) == 0 {
		log.Print("All games are already synchronized.")
		return
	}
	log.Printf("Creating %d games", len(games))
	if err := repo.CreateGames(games); err != nil {
		log.Panic(err)
	}

	savedCheckpoint, err := repo.SaveCheckpoint(checkpoint.Line + len(response.Values))
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Saved checkpoint %d", savedCheckpoint.Line)
}
