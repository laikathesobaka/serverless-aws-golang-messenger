package chatsess

import (
	"html"
	"serverless-chat-app/functions/timestamp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const dateFormat = "02-01-2002"

type Chat struct {
	CreatedDate      string
	CreatedTimestamp time.Time
	Username         string
	Text             string
}

func NewChat(username, text string) Chat {
	return Chat{
		CreatedDate:      time.Now().Format(dateFormat),
		CreatedTimestamp: time.Now(),
		Username:         username,
		Text:             html.EscapeString(text),
	}
}

func (c Chat) Put(s *session.Session) error {
	db := dynamodb.New(s)
	_, err := db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("chats"),
		Item: map[string]*dynamodb.AttributeValue{
			"created_date":      {S: aws.String(c.CreatedDate)},
			"created_timestamp": {N: timestamp.TimeToDB(c.CreatedTimestamp)},
			"username":          {S: aws.String(c.Username)},
			"text":              {S: aws.String(c.Text)},
		},
	})
	return err
}

func GetChats(s *session.Session) ([]Chat, error) {
	db := dynamodb.New(s)
	now := aws.String(time.Now().Format(dateFormat))
	cond := expression.Key("created_date").Equal(expression.Value(now))
	expr, err := expression.NewBuilder().WithKeyCondition(cond).Build()
	if err != nil {
		return []Chat{}, err
	}
	params := &dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName:                 aws.String("chats"),
	}
	dbRes, err := db.Query(params)
	if err != nil {
		return []Chat{}, err
	}
	res := []Chat{}
	err = dynamodbattribute.UnmarshalListOfMaps(dbRes.Items, &res)
	if err != nil {
		return []Chat{}, err
	}
	return res, nil
}

func GetChatsAfter(createdDate string, t time.Time, s *session.Session) ([]Chat, error) {
	db := dynamodb.New(s)
	now := aws.String(time.Now().Format(dateFormat))
	cond := expression.Key("created_date").Equal(expression.Value(now))
	expr, err := expression.NewBuilder().WithKeyCondition(cond).Build()
	params := &dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ExclusiveStartKey: map[string]*dynamodb.AttributeValue{
			"created_date":      {S: aws.String(createdDate)},
			"created_timestamp": {N: timestamp.TimeToDB(t)},
		},
		TableName: aws.String("chats"),
	}
	dbRes, err := db.Query(params)
	if err != nil {
		return []Chat{}, err
	}
	res := []Chat{}
	err = dynamodbattribute.UnmarshalListOfMaps(dbRes.Items, &res)
	if err != nil {
		return []Chat{}, err

	}
	return res, nil
}
