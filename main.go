package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type App struct {
	Config Config
}

type SnsReply struct {
	XMLName   xml.Name `xml:"PublishResponse"`
	Namespace string   `xml:"xmlns,attr"`
	MessageID string   `xml:"PublishResult>MessageId"`
	RequestID string   `xml:"ResponseMetadata>ResponseMetadata"`
}

func (app *App) handleRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Something went wrong reading request body", 500)
		fmt.Printf("Something went wrong reading request body: %s\n", err)

		return
	}

	values, err := url.ParseQuery(string(body))

	if err != nil {
		http.Error(w, "Can't parse request body", 400)
		fmt.Printf("Can't parse request body: %s\n", err)

		return
	}

	if values["Action"][0] != "Publish" {
		http.Error(w, "I can only handle Publish action", 400)
		fmt.Printf("I can only handle Publish action: %s\n", err)

		return
	}

	messageID := pseudoUUID()

	for _, s := range app.Config.Subscriptions {
		if s.Topic == values["TopicArn"][0] {
			s.Publish(messageID, values["Message"][0])
		}
	}

	reply, _ := xml.Marshal(&SnsReply{
		MessageID: messageID,
		RequestID: pseudoUUID(),
		Namespace: "http://sns.amazonaws.com/doc/2010-03-31/",
	})

	fmt.Fprintf(w, "%s", string(reply))
}

func main() {
	config, err := getConfig()

	if err != nil {
		fmt.Printf("Error parsing config.json: %s\n", err)
		return
	}

	app := new(App)
	app.Config = config

	for _, s := range config.Subscriptions {
		fmt.Printf("New subscription: [%s] -> %s\n", s.Topic, s.QueueName)
	}

	http.HandleFunc("/", app.handleRequest)
	http.ListenAndServe(":"+config.Port, nil)
}
