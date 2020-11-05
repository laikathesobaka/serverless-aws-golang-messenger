package main

import (
	"bytes"
	"context"
	"encoding/json"
	"serverless-chat-app/chat_session"
	"serverless-chat-app/user_session"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Event struct {
	Username string
	Password string
}

type Response events.APIGatewayProxyResponse

type Result struct {
	Job       string
	Err       string
	SessionID string
	Username  string
}

func handler(ctx context.Context, e Event) (Response, error) {
	sess := session.Must(session.NewSession())
	user, err := chat_session.GetDBUserPW(e.Username, e.Password, sess)
	if err != nil {
		res := Result{
			Job: "Login",
			Err: err.Error(),
		}
		resp, _ := EncodeResponse(res, 401)
		return resp, nil
	}
	lg := user_session.NewLogin(e.Username)
	err = lg.Put(sess)
	if err != nil {
		res := Result{
			Job: "Login",
			Err: err.Error(),
		}
		resp, _ := EncodeResponse(res, 500)
		return resp, nil
	}
	res := Result{
		Job:       "Login",
		SessionID: lg.SessionID,
		Username:  user.Username,
	}
	resp, _ := EncodeResponse(res, 200)
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
