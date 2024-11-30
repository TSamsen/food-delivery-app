package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"food-delivery-app/internal/adapters"
	"food-delivery-app/internal/core"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func PlaceOrder(client *redis.Client, cont echo.Context, kafkaProducer *adapters.KafkaProducer) (*core.Order, error) {
	var ctx = context.Background()
	var request core.OrderRequest

	bodyBytes, err := io.ReadAll(cont.Request().Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}
	fmt.Printf("Request body: %s\n", string(bodyBytes))

	cont.Request().Body = io.NopCloser(bytes.NewReader(bodyBytes))

	err = cont.Bind(&request)
	if err != nil {
		return nil, fmt.Errorf("request does not match: %w", err)
	}

	order_id := uuid.New().String()
	cacheKey := "orders" + order_id
	orderDetailKey := "order_detail:" + order_id

	order := &core.Order{
		OrderID: order_id,
		Status:  "created",
	}

	order_detail := &core.OrderDetail{
		OrderID:      order_id,
		Items:        request.Items,
		RestaurantID: request.RestaurantID,
		Status:       "created",
	}

	orderDetailJSON, err := json.Marshal(order_detail)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize order detail to JSON: %w", err)
	}
	err = client.Set(ctx, orderDetailKey, orderDetailJSON, 10*time.Minute).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to cache order detail: %w", err)
	}

	orderJSON, err := json.Marshal(order)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize order to JSON: %w", err)
	}
	err = client.Set(ctx, cacheKey, orderJSON, 10*time.Minute).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to cache order status: %w", err)
	}

	fmt.Println("Cached order", order_id)

	menuMessage := ""

	for i := 0; i < len(request.Items); i++ {
		menuMessage += fmt.Sprint("Menu " + request.Items[i].MenuID + " quantity " + string(request.Items[i].Quantity) + "\n")
	}

	kafkaMessage := fmt.Sprintf(`Order %s created to restaurant %s list menu \n %s`, order_id, request.RestaurantID, menuMessage)
	err = kafkaProducer.Publish("order_created", kafkaMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to publish message to Kafka: %w", err)
	}
	fmt.Println("Published order", order_id, "to Kafka topic 'order_created'")

	notification_request := &core.NotificationRequest{
		Recipient: "Customer",
		OrderId:   order_id,
		Message:   "Order " + order_id + " created",
	}

	err = SendNotification(*notification_request)
	if err != nil {
		return nil, fmt.Errorf("failed to send notification: %w", err)
	}

	fmt.Println("Notification sent for order", order_id)

	return order, nil
}
