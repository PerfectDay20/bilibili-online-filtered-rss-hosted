package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"time"
)

type HistoryRssRecord struct {
	ID              string `dynamodbav:"ID"`
	Record          string `dynamodbav:"Record"`
	UpdateTimestamp int64  `dynamodbav:"UpdateTimestamp"`
}

func (record HistoryRssRecord) GetKey() map[string]types.AttributeValue {
	id, err := attributevalue.Marshal(record.ID)
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{"ID": id}
}

type TableBasics struct {
	TableName      string
	DynamoDbClient *dynamodb.Client
}

const RecordId = "dummy"

// SetRecord add rss string content to dynamodb, the key is fixed
func (basics TableBasics) SetRecord(content string) error {
	record := HistoryRssRecord{ID: RecordId, Record: content, UpdateTimestamp: time.Now().Unix()}
	item, err := attributevalue.MarshalMap(record)
	if err != nil {
		return err
	}
	_, err = basics.DynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(basics.TableName), Item: item,
	})
	if err != nil {
		return err
	}
	log.Println("Successful set record")
	return err
}

// GetRecord get rss string content from dynamodb, the key is fixed
func (basics TableBasics) GetRecord() (HistoryRssRecord, error) {
	record := HistoryRssRecord{ID: RecordId}
	response, err := basics.DynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(basics.TableName), Key: record.GetKey(),
	})
	if err != nil {
		log.Printf("Couldn't get info about %v. Here's why: %v\n", RecordId, err)
	} else {
		if response.Item == nil {
			log.Println("Content not exists in DynamoDB")
			return record, nil
		} else {
			err = attributevalue.UnmarshalMap(response.Item, &record)
			if err != nil {
				log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
			}
		}

	}
	return record, err
}

func initTable() *TableBasics {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	tableBasics := TableBasics{"rss_cache", dynamodb.NewFromConfig(sdkConfig)}
	return &tableBasics
}
