package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"food-delivery-app/internal/adapters"
	"food-delivery-app/internal/core"
	"io"
	"io/ioutil"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

func GetRiders(client *redis.Client) (string, error) {
	var ctx = context.Background()
	//Get from cache
	restaurants, err := client.Get(ctx, "riders").Result()
	if err == nil {

		fmt.Println("Load riders from cache")
		return string(restaurants), nil
	}

	//Get from files
	fmt.Println("Load riders from files")
	data, err := ioutil.ReadFile("mock_data/rider.json")
	if err != nil {
		return "", err
	}

	var strData map[string]interface{}
	err = json.Unmarshal(data, &strData)
	if err != nil {
		return "", err
	}

	jsonString, err := json.Marshal(strData) // Converts map back to JSON
	if err != nil {
		return "", err
	}

	//Cache restaurants in Redis
	client.Set(ctx, "riders", jsonString, 10*time.Minute)
	fmt.Println("Keeo riders to cache")

	return string(jsonString), nil
}

func RiderPickupConfimation(ctx echo.Context, kafkaProducer *adapters.KafkaProducer) (*core.PickupConfirmationResponse, error) {
	var request core.PickupConfirmationRequest
	bodyBytes, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}
	fmt.Printf("Request body: %s\n", string(bodyBytes))

	ctx.Request().Body = io.NopCloser(bytes.NewReader(bodyBytes))

	err = ctx.Bind(&request)
	if err != nil {
		return nil, fmt.Errorf("request does not match: %w", err)
	}

	kafkaMessage := fmt.Sprintf(`%s picked up`, request.OrderId)
	err = kafkaProducer.Publish("order_picked_up", kafkaMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to publish message to Kafka: %w", err)
	}
	fmt.Println("Published order", request.OrderId, "to Kafka topic 'order_picked_up'")

	response := &core.PickupConfirmationResponse{
		Status: "picked_up",
	}

	return response, nil
}

func RiderDeliverOrder(ctx echo.Context, kafkaProducer *adapters.KafkaProducer) (*core.DeliveryOrderResponse, error) {
	var request core.DeliveryOrderRequest
	bodyBytes, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}
	fmt.Printf("Request body: %s\n", string(bodyBytes))

	ctx.Request().Body = io.NopCloser(bytes.NewReader(bodyBytes))

	err = ctx.Bind(&request)
	if err != nil {
		return nil, fmt.Errorf("request does not match: %w", err)
	}

	kafkaMessage := fmt.Sprintf(`%s delivered`, request.OrderId)
	err = kafkaProducer.Publish("order_delivered", kafkaMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to publish message to Kafka: %w", err)
	}
	fmt.Println("Published order", request.OrderId, "to Kafka topic 'order_delivered'")

	response := &core.DeliveryOrderResponse{
		Status: "delivered",
	}

	return response, nil
}
