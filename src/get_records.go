package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

func extractKinesisEvent() {
	// Create a new AWS session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create a Kinesis client
	svc := kinesis.New(sess)

	// Get a shard iterator for the beginning of the shard
	streamName := "ecommerce_unauth_dev"
	shardId := "shard-000000000000"
	shardIteratorInput := &kinesis.GetShardIteratorInput{
		StreamName:        aws.String(streamName),
		ShardId:           aws.String(shardId),
		ShardIteratorType: aws.String("LATEST"),
	}
	shardIteratorOutput, err := svc.GetShardIterator(shardIteratorInput)
	if err != nil {
		panic(err)
	}
	shardIterator := shardIteratorOutput.ShardIterator

	// Set up the input parameters
	params := &kinesis.GetRecordsInput{
		ShardIterator: shardIterator,
		StreamARN:     aws.String("arn:aws:kinesis:us-east-1:<account-id>:stream/ecommerce_unauth_dev"),
		Limit:         aws.Int64(100),
	}

	// Get the records from the stream
	resp, err := svc.GetRecords(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print out the records
	for _, r := range resp.Records {
		fmt.Println(string(r.Data))
	}
}
