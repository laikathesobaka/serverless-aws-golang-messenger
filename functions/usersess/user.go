package usersess

import (
	"fmt"
	"serverless-chat-app/functions/password"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type User struct {
	Username string
	Password string
}

func NewUser(name, pw string) User {
	return User{
		Username: name,
		Password: password.NewPassword(pw),
	}
}

func (u User) Put(s *session.Session) error {
	db := dynamodb.New(s)
	_, err := db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("users"),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {S: aws.String(u.Username)},
			"password": {S: aws.String(u.Password)},
		},
	})
	return err
}

func GetDBUser(name string, s *session.Session) (User, error) {
	fmt.Println("add pkg")

	db := dynamodb.New(s)
	res, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("users"),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {S: aws.String(name)},
		},
	})
	if err != nil {
		return User{}, err
	}
	if res.Item == nil {
		return User{}, fmt.Errorf("User with username %s not found", name)
	}

	pw, ok := res.Item["password"]
	if !ok {
		return User{}, fmt.Errorf("User has no password %s", name)
	}
	user := User{Username: name, Password: *(pw.S)}
	return user, err
}

func GetDBUserPW(name, pw string, s *session.Session) (User, error) {
	user, err := GetDBUser(name, s)
	// if err != nil {
	// 	return user, err
	// }
	if !password.CheckPassword(pw, user.Password) {
		return User{}, fmt.Errorf("Password doesn't match")
	}
	return user, err
}
