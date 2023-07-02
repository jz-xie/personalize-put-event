package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

type MyData struct {
	EventID         string          `json:"EventID"`
	SessionID       string          `json:"SessionID"`
	UTM             string          `json:"UTM"`
	IP              string          `json:"IP"`
	UserID          string          `json:"UserID"`
	UnixTimestamp   float64         `json:"UnixTimestamp"`
	EventCategory   string          `json:"EventCategory"`
	EventAction     string          `json:"EventAction"`
	EventLabel      string          `json:"EventLabel"`
	EventValue      int             `json:"EventValue"`
	EventDimensions EventDimensions `json:"EventDimensions"`
}

type EventDimensions struct {
	RecommendedItems *[]string `json:"RecommendedItems"`
	Items            *[]string `json:"Items"`
	Cart             *[]string `json:"Cart"`
}

func readJSONFile(filename string) ([]MyData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []MyData
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 104857600)
	for scanner.Scan() {
		var d MyData
		err := json.Unmarshal(scanner.Bytes(), &d)
		if err != nil {
			return nil, err
		}
		data = append(data, d)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

// generate kinesis test input
func generateKinesisInput() []events.KinesisEvent {
	data, err := readJSONFile("events_input_test.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}
	inputs := make([]events.KinesisEvent, 0)
	for _, record := range data {
		content, _ := json.Marshal(record)
		event := events.KinesisEvent{
			Records: []events.KinesisEventRecord{
				{
					Kinesis: events.KinesisRecord{
						Data:           content,
						PartitionKey:   "partition-key-1",
						SequenceNumber: "1234567890",
					},
					EventSource: "aws:kinesis",
				},
			},
		}
		inputs = append(inputs, event)
	}
	return inputs
}
