package main

import (
	"bytes"
	"context"
	"encoding/json"
	"serverless-chat-app/functions/login"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Event struct {
	SessionID string
	Text      string
}

type Result struct {
	Job string
	Err string
}
type Response events.APIGatewayProxyResponse

func handler(ctx context.Context, e Event) (Response, error) {
	sess := session.Must(session.NewSession())
	login, err := login.GetLogin(e.SessionID, sess)
	if err != nil {
		result := Result{
			Job: "Chat",
			Err: "Not logged in: " + err.Error(),
		}
		response, _ := EncodeResponse(result, 400)
	}
	chat := chatsess.NewChat(login.Username, e.Text)
	err = chat.Put(sess)
	if err != nil {
		result = Result{
			Job: "Chat",
			Err: "Could not create chat: " + err.Error()
		}
		response, _ = EncodeResponse(result, 500)
	}
	response := EncodeResponse(Result{Job: "Chat " + e.Text}, 200)
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
