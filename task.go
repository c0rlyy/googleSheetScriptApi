package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	credentialsPath = "./credentials.json"
	pathToConfig    = "./config.json"
)

func LoadSheet(fp string) (*sheets.Service, error) {

	if _, err := os.Stat(fp); os.IsNotExist(err) {
		return nil, fmt.Errorf("nie znaleziono pliku %s", fp)
	}

	srv, err := sheets.NewService(context.Background(), option.WithCredentialsFile(fp))
	if err != nil {
		return nil, fmt.Errorf("error while truing to get Sheet service %v", err)
	}

	return srv, nil
}

func Task(req RequestBody) {
	sheetsService, err := LoadSheet(credentialsPath)
	if err != nil {
		log.Fatal(err)
	}

	config, err := NewConfig(pathToConfig)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := sheetsService.Spreadsheets.Values.Get(config.FirstSheetID, config.FirstSheetRange).Do()
	if err != nil {
		log.Println(err)
		return
	}

	sheetVals := resp.Values
	googleSheet := Sheet{values: sheetVals}
	sheetHeaders := googleSheet.values[0]

	sheetContainter, err := googleSheet.ConvertToMap()
	if err != nil {
		log.Println(err)
		return
	}

	targetValue, err := FilterSheetContainter(&sheetContainter, config, req.KeyValue, req.SerialNumberValue)
	if err != nil {
		log.Println(err)
		return
	}

	googleSheet.ConvertToValues(&sheetContainter, sheetHeaders)

	resp, err = sheetsService.Spreadsheets.Values.Get(config.SecondSheetID, config.SecondSheetRange).Do()
	if err != nil {
		log.Println(err)
		return
	}

	secondSheet := resp.Values

	secondGoogleSheet := Sheet{values: secondSheet}
	sheetHeaders = secondGoogleSheet.values[0]

	secondContainter, _ := secondGoogleSheet.ConvertToMap()
	err = ReplaceSheetKey(&secondContainter, config, targetValue, req.KeyValue)
	if err != nil {
		log.Println(err)
		return
	}

	secondGoogleSheet.ConvertToValues(&secondContainter, sheetHeaders)

	payload := sheets.ValueRange{Values: googleSheet.values}
	res, err := sheetsService.Spreadsheets.Values.Update(config.FirstSheetID, config.FirstSheetRange, &payload).ValueInputOption("USER_ENTERED").Do()

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(res.UpdatedData)

	payload = sheets.ValueRange{Values: secondGoogleSheet.values}
	res, err = sheetsService.Spreadsheets.Values.Update(config.SecondSheetID, config.SecondSheetRange, &payload).ValueInputOption("USER_ENTERED").Do()

	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(res.UpdatedData)
}
