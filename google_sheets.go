package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

type SpreadsheetConfig struct {
	service                *spreadsheet.Service
	sheetId                string
	personalSheetTitle     string
	organizationSheetTitle string
}

func NewSpreadsheetConfig() *SpreadsheetConfig {
	data, err := os.ReadFile("service-account-key.json")
	checkError(err)

	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	checkError(err)

	client := conf.Client(context.TODO())
	service := spreadsheet.NewServiceWithClient(client)

	spreadsheetId := getSpreadsheetId()

	return &SpreadsheetConfig{
		service:                service,
		sheetId:                spreadsheetId,
		personalSheetTitle:     "Индивидуальные карточки",
		organizationSheetTitle: "Карточки организаций",
	}
}

func getSpreadsheetId() string {
	id, exists := os.LookupEnv("SPREADSHEET_ID")

	if !exists {
		log.Print("Spreadsheet id not found in .env")
	}

	return id
}

func (ss *SpreadsheetConfig) getSpreadsheetById(sheetId string) (*spreadsheet.Spreadsheet, error) {
	fetchSpreadsheet, err := ss.service.FetchSpreadsheet(sheetId)
	checkError(err)

	return &fetchSpreadsheet, err
}

func (ss *SpreadsheetConfig) getHeadersOfSheet(sheet *spreadsheet.Sheet) ([]spreadsheet.Cell, string) {
	var headersString string

	headers := sheet.Rows[0]

	for _, cell := range headers {
		headersString += cell.Value
	}

	return headers, headersString
}

func (ss *SpreadsheetConfig) getCardByNumber(sheetTitle string, cardNumber int) string {
	spreadsheetById, err := ss.getSpreadsheetById(ss.sheetId)
	checkError(err)

	sheet, err := spreadsheetById.SheetByTitle(sheetTitle)
	checkError(err)

	headersRow, _ := ss.getHeadersOfSheet(sheet)
	cardRow := sheet.Rows[cardNumber]

	var card strings.Builder
	cardFormat := "<b>%s</b>\n%s\n\n"

	for i := 0; i < len(headersRow) && i < len(cardRow); i++ {
		card.WriteString(fmt.Sprintf(cardFormat, headersRow[i].Value, cardRow[i].Value))
	}

	return card.String()
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
