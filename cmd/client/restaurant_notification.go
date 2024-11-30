package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"food-delivery-app/internal/core"
)

const NotificationSendUrl = "http://localhost:8088/notification/send"

func SendNotification(request core.NotificationRequest) error {
	fmt.Printf("Sending Request /notification/send: %v", request)
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return err
	}
	response, err := http.Post(NotificationSendUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		message := fmt.Sprintf("Non-OK HTTP status: %d\n", response.StatusCode)
		return errors.New(message)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	var result core.NotificationResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	fmt.Printf("Receive Response /notification/send: %v", result)
	return nil
}
