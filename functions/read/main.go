package main

import (
	"bytes"
	"context"
	"encoding/json"
	"serverless-chat-app/functions/chatsess"
	"serverless-chat-app/functions/login"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Event struct {
	SessionID     string
	LastDate      string
	LastTimestamp string
}

type Response events.APIGatewayProxyResponse

type Result struct {
	Job   string
	Err   string
	Chats []chatsess.Chat
}

func handler(ctx context.Context, e Event) (Response, error) {
	sess := session.Must(session.NewSession())
	_, err := login.GetLogin(e.SessionID, sess)
	if err != nil {
		result := Result{
			Job: "Read",
			Err: "Not logged in: " + err.Error(),
		}
		response, _ := EncodeResponse(result, 400)
		return response, nil
	}
	if e.LastDate != "" {
		ltime, err := time.Parse(time.RFC3339, e.LastTimestamp)
		if err != nil {
			result := Result{
				Job: "Read",
				Err: err.Error(),
			}
			response, _ := EncodeResponse(result, 500)
			return response, nil
		}
		chats, err := chatsess.GetChatsAfter(e.LastDate, ltime, sess)
		if err != nil {
			result := Result{
				Job: "Read",
				Err: "Cannot read: " + err.Error(),
			}
			response, _ := EncodeResponse(result, 500)
			return response, nil
		}
		response, _ := EncodeResponse(Result{Job: "Read", Chats: chats}, 200)
		return response, nil
	}
	chats, err := chatsess.GetChats(sess)
	if err != nil {
		result := Result{
			Job: "Read",
			Err: "Cannot read: " + err.Error(),
		}
		response, _ := EncodeResponse(result, 500)
		return response, nil
	}
	response, _ := EncodeResponse(Result{Job: "Read", Chats: chats}, 200)
	return response, nil
}

func main() {
	lambda.Start(handler)
}

func EncodeResponse(res Result, statusCode int) (Response, error) {
	var buf bytes.Buffer
	body, err := json.Marshal(res)
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)
	resp := Response{
		StatusCode:      statusCode,
		Body:            buf.String(),
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	return resp, nil
}
