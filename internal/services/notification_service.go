package services

import (
	"fmt"

	"food-delivery-app/internal/core"
)

const NotificationSendUrl = "http://localhost:8088/notification/send"

func SendNotification(request core.NotificationRequest) error {
	fmt.Printf("Sending Request /notification/send: %v\n", request)

	notification_response := &core.NotificationResponse{
		Status: "sent",
	}

	fmt.Printf("Receive Response /notification/send: %v\n", notification_response)
	return nil
}
