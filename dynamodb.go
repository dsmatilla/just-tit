package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"time"
)

const dynamodbRegion = "eu-west-1"
const dynamodbTable = "JustTit"
const secondsToCache = 60 * 60 * 24

type JustTitCache struct {
	ID        string `json:"id"`
	Result    string `json:"result"`
	Timestamp int64  `json:"timestamp"`
}

func getFromDB(ID string) JustTitCache {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(dynamodbRegion)},
	)
	svc := dynamodb.New(sess)

	result, _ := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(dynamodbTable),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(ID),
			},
		},
	})

	cache := JustTitCache{}
	if result != nil {
		dynamodbattribute.UnmarshalMap(result.Item, &cache)
	}

	return cache
}

func putToDB(ID string, Result string) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(dynamodbRegion)},
	)
	svc := dynamodb.New(sess)

	cache := JustTitCache{ID, Result, time.Now().Unix() + secondsToCache}
	item, _ := dynamodbattribute.MarshalMap(cache)
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(dynamodbTable),
	}
	_, _ = svc.PutItem(input)
}
