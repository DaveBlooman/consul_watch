package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/twinj/uuid"
)

// CheckData info
type CheckData struct {
	CheckID     string `json:"CheckID"`
	Name        string `json:"Name"`
	Node        string `json:"Node"`
	Notes       string `json:"Notes"`
	Output      string `json:"Output"`
	ServiceID   string `json:"ServiceID"`
	ServiceName string `json:"ServiceName"`
	Status      string `json:"Status"`
}

func main() {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	var data []CheckData

	err = json.Unmarshal(bytes, &data)

	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(session.New(), aws.NewConfig().WithRegion("eu-west-1"))

	now := time.Now()
	secs := strconv.FormatInt(now.Unix(), 10)

	var putreq []*dynamodb.WriteRequest

	for _, checkItem := range data {

		if checkItem.ServiceName == "" {
			continue
		}

		putreq = append(putreq, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: map[string]*dynamodb.AttributeValue{
					"id": {
						S: aws.String(uuid.NewV4().String()),
					},
					"timestamp": {
						S: aws.String(secs),
					},
					"CheckID": {
						S: aws.String(checkItem.CheckID),
					},
					"Status": {
						S: aws.String(checkItem.Status),
					},
					"ServiceName": {
						S: aws.String(checkItem.ServiceName),
					},
					"ServiceID": {
						S: aws.String(checkItem.ServiceID),
					},
					"Node": {
						S: aws.String(checkItem.Node),
					},
					"Output": {
						S: aws.String(checkItem.Output),
					},
				},
			},
		},
		)
	}

	for i := 25; i < len(putreq); i += 25 {
		params := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				"dashboard": putreq[i-25 : i],
			},
		}

		_, err = svc.BatchWriteItem(params)
		if err != nil {
			panic(err)
		}
	}
}
