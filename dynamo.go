package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log/slog"
	"time"
)

type HistoryRssRecord struct {
	ID              string `dynamodbav:"ID"`
	Record          string `dynamodbav:"Record"`
	UpdateTimestamp int64  `dynamodbav:"UpdateTimestamp"`
}

func (h HistoryRssRecord) GetKey() map[string]types.AttributeValue {
	id, err := attributevalue.Marshal(h.ID)
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
func (t TableBasics) SetRecord(content string) error {
	record := HistoryRssRecord{ID: RecordId, Record: content, UpdateTimestamp: time.Now().Unix()}
	item, err := attributevalue.MarshalMap(record)
	if err != nil {
		return err
	}
	_, err = t.DynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName), Item: item,
	})
	if err != nil {
		return err
	}
	slog.Info("Successful set record")
	return err
}

// GetRecord get rss string content from dynamodb, the key is fixed
func (t TableBasics) GetRecord() (HistoryRssRecord, error) {
	record := HistoryRssRecord{ID: RecordId}
	response, err := t.DynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(t.TableName), Key: record.GetKey(),
	})
	if err != nil {
		slog.Error("Couldn't get info record", "record", RecordId, "reason", err)
	} else {
		if response.Item == nil {
			slog.Info("Content not exists in DynamoDB")
			return record, nil
		} else {
			err = attributevalue.UnmarshalMap(response.Item, &record)
			if err != nil {
				slog.Error("Couldn't unmarshal response", "reason", err)
			}
		}

	}
	return record, err
}

func initDynamoTable() (*TableBasics, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		slog.Error("unable to load SDK config", "reason", err)
		return nil, err
	}

	tableBasics := TableBasics{"rss_cache", dynamodb.NewFromConfig(sdkConfig)}
	return &tableBasics, nil
}
