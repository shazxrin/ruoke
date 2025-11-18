package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

var (
	ErrNotifierNotify = errors.New("notifier failed to send a notification")
)

type Notifier interface {
	Notify(title, message string) error
}

type PushoverNotifier struct {
	appToken  string
	userToken string

	httpClient *http.Client
}

type PushoverNotification struct {
	Token   string `json:"token"`
	User    string `json:"user"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

func NewPushoverNotifier(appToken, userToken string) Notifier {
	return &PushoverNotifier{
		appToken:  appToken,
		userToken: userToken,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (n *PushoverNotifier) Notify(title, message string) error {
	pushoverNotif := PushoverNotification{
		Token:   n.appToken,
		User:    n.userToken,
		Title:   title,
		Message: message,
	}

	notificationJson, err := json.Marshal(pushoverNotif)
	if err != nil {
		log.Printf("Error marshalling Pushover message: %v\n", err)
		return ErrNotifierNotify
	}

	resp, err := n.httpClient.Post(
		"https://api.pushover.net/1/messages.json",
		"application/json",
		bytes.NewBuffer(notificationJson),
	)

	if err != nil {
		log.Printf("Error sending Pushover notification: %v\n", err)
		return ErrNotifierNotify
	}

	if resp.StatusCode != 200 {
		log.Printf("Error sending Pushover notification: %v\n", resp.Status)
		return ErrNotifierNotify
	}

	return nil
}
