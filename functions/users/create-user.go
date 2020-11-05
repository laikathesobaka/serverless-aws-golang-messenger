package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"serverless-chat-app/chat_session"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration

type Event struct {
	Username string
	Password string
}

type Response events.APIGatewayProxyResponse

type Result struct {
	Job string
	Err string
}

func handler(ctx context.Context, e Event) (Response, error) {
	fmt.Println("event --------------", e)
	s := session.Must(session.NewSession())
	_, err := chat_session.GetDBUser(e.Username, s)
	if err == nil {
		result := Result{Job: "Add user", Err: "User already exists"}
		resp, _ := EncodeResponse(result, 500)
		return resp, nil
	}
	user := chat_session.NewUser(e.Username, e.Password)
	err = user.Put(s)
	if err != nil {
		result := Result{
			Job: "Add user",
			Err: "Could not insert user: " + err.Error(),
		}
		resp, _ := EncodeResponse(result, 500)
		return resp, nil
	}
	result := Result{Job: "Add user", Err: ""}
	resp, _ := EncodeResponse(result, 200)
	return resp, nil
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

func main() {
	lambda.Start(handler)
}
