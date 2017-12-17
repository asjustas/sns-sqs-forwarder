package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type SqsSubscription struct {
	*Subscription
	QueueName string
	Endpoint  string
	Raw       bool
}

func (s SqsSubscription) Publish(id string, msg string) error {
	messageBody := msg

	if !s.Raw {
		snsMessage := map[string]string{
			"Type":      "Notification",
			"MessageId": id,
			"Message":   msg,
			"Timestamp": time.Now().UTC().Format(time.RFC3339),
			"TopicArn":  s.Topic,
		}

		snsMessageJSON, err := json.Marshal(snsMessage)

		if err != nil {
			return err
		}

		messageBody = string(snsMessageJSON)
	}

	fmt.Printf("Dispatching to: [%s] -> %s\n", s.QueueName, messageBody)

	resp, err := http.PostForm(
		fmt.Sprintf("%s?QueueName=%s", s.Endpoint, s.QueueName),
		url.Values{
			"Action":      {"SendMessage"},
			"Version":     {"2012-11-05"},
			"QueueUrl":    {s.QueueName},
			"MessageBody": {messageBody},
		})

	if err != nil {
		return err
	}

	resp.Body.Close()

	return nil
}