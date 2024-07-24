package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	FirstSheetID           string `json:"FirstSheetID"`
	SecondSheetID          string `json:"SecondSheetID"`
	FirstSheetRange        string `json:"FirstSheetRange"`
	SecondSheetRange       string `json:"SecondSheetRange"`
	TargetValueHeader      string `json:"TargetValueHeader"`
	SerialKeyHeader        string `json:"SerialKeyHeader"`
	KeyIDHeader            string `json:"KeyIDHeader"`
	SecondSheetKeyIDHeader string `json:"SecondSheetKeyIDHeader"`
}

func NewConfig(fp string) (*Config, error) {
	jsonFile, err := os.Open(fp)
	if err != nil {
		return nil, fmt.Errorf("error while trying to read the file: %w", err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("error while reading config file: %w", err)
	}

	var config Config
	err = json.Unmarshal(byteValue, &config)

	if err != nil {
		return nil, fmt.Errorf("error while converting config file to struct: %w", err)
	}

	return &config, nil
}
