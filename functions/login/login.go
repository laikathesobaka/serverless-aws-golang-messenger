package login

import (
	"crypto/rand"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Login struct {
	SessionID string
	Username  string
}

func NewLogin(name string) Login {
	b := make([]byte, 20)
	rand.Read(b)
	return Login{
		SessionID: fmt.Sprintf("%x", b),
		Username:  name,
	}
}

func GetLogin(sessionID string, s *session.Session) (Login, error) {
	db := dynamodb.New(s)
	res, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("sessions"),
		Key:       map[string]*dynamodb.AttributeValue{"id": {S: aws.String(sessionID)}},
	})
	if err != nil {
		return Login{}, err
	}
	if res.Item == nil {
		return Login{}, fmt.Errorf("Session not found")
	}
	username, ok := res.Item["username"]
	if !ok {
		return Login{}, fmt.Errorf("Session with username not found")
	}
	return Login{SessionID: sessionID, Username: *(username.S)}, nil
}

func (l Login) Put(s *session.Session) error {
	db := dynamodb.New(s)
	_, err := db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("sessions"),
		Item: map[string]*dynamodb.AttributeValue{
			"id":       {S: aws.String(l.SessionID)},
			"username": {S: aws.String(l.Username)},
		},
	})
	return err
}
