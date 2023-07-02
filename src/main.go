package main

import (
	"context"
	"encoding/json"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/personalizeevents"
)

var trackingID string = "40b1a63a-1027-454a-af5a-86b5f563ad52"

func extractSKU(url string) *string {
	// Define a regular expression to match the product code
	re := regexp.MustCompile(`/(R[SCZ]\d{2}-\d{4}[A-Z0-9-]+)(?:/|$)`)
	// Find the first match of the regular expression in the URL
	match := re.FindStringSubmatch(url)
	if len(match) > 1 {
		return &match[1]
	}
	return nil

}

// Parse event data for EventCategory==PageView
func parsePageView(data MyData) *personalizeevents.PutEventsInput {

	url := data.EventLabel
	containSKU := strings.Contains(url, "RZ") || strings.Contains(url, "RS") || strings.Contains(url, "RC")
	if containSKU {
		sku := extractSKU(url)
		if sku != nil {
			// Create a new PutEventsInput
			putEventsInput := personalizeevents.PutEventsInput{
				TrackingId: aws.String(trackingID),
				UserId:     aws.String(data.UserID),
				SessionId:  aws.String(data.SessionID),
				EventList: []*personalizeevents.Event{
					{
						EventType: aws.String("View"),
						SentAt:    aws.Time(time.Unix(int64(data.UnixTimestamp), 0)),
						ItemId:    aws.String((*sku)[:13]),
					},
				},
			}
			return &putEventsInput
		}
	}
	return nil
}

// Parse event data for EventCategory==Recommendation
func parseRecommendation(data MyData) *personalizeevents.PutEventsInput {
	if data.EventAction != "Load" {
		return nil
	}

	skuList := data.EventDimensions.RecommendedItems

	var eventList []*personalizeevents.Event
	for _, sku := range *skuList {
		event := &personalizeevents.Event{
			EventType: aws.String("View"),
			SentAt:    aws.Time(time.Unix(int64(data.UnixTimestamp), 0)),
			ItemId:    aws.String(sku[:13]),
		}
		eventList = append(eventList, event)
	}
	putEventsInput := personalizeevents.PutEventsInput{
		TrackingId: aws.String(trackingID),
		UserId:     aws.String(data.UserID),
		SessionId:  aws.String(data.SessionID),
		EventList:  eventList,
	}
	return &putEventsInput

}

func writeEventTracker(putEventsInput *personalizeevents.PutEventsInput) {

	// fmt.Println(putEventsInput)
	// Create a new session to PersonalizeEvents
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Printf("Error creating session: %v", err)
		// continue
	}
	svc := personalizeevents.New(sess)

	// fmt.Println("PRINT OUT EVENTS\n ")

	//Put the events into the Personalize event tracker
	_, err = svc.PutEvents(putEventsInput)
	if err != nil {
		log.Printf("Error putting events: %v", err)
	}
}

func Handler(ctx context.Context, kinesisEvent events.KinesisEvent) error {
	// fmt.Println(kinesisEvent)
	for _, record := range kinesisEvent.Records {

		var kinesisData MyData
		err := json.Unmarshal(record.Kinesis.Data, &kinesisData)
		if err != nil {
			log.Printf("Error unmarshalling Kinesis data: %v", err)
			continue
		}

		var putEventsInput *personalizeevents.PutEventsInput
		switch kinesisData.EventCategory {
		case "Pageview":
			putEventsInput = parsePageView(kinesisData)
		case "Recommendation":
			putEventsInput = parseRecommendation(kinesisData)
		default:
			log.Printf("Unknown event category: %s", kinesisData.EventCategory)
		}

		if putEventsInput != nil {
			// fmt.Printf("%+v\n", *putEventsInput)
			writeEventTracker(putEventsInput)
		}
	}
	return nil
}

func main() {
	lambda.Start(Handler)

	// testData := generateKinesisInput()
	// for _, d := range testData {
	// 	Handler(context.Background(), d)
	// }

}
