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

func GetRestaurants(client *redis.Client) (string, error) {
	var ctx = context.Background()

	restaurants, err := client.Get(ctx, "restaurants").Result()
	if err == nil {
		fmt.Println("Load restaurants from cache")
		return string(restaurants), nil
	}

	fmt.Println("Load restaurants from files")
	data, err := ioutil.ReadFile("mock_data/restaurants.json")
	if err != nil {
		return "", err
	}

	var strData map[string]interface{}
	err = json.Unmarshal(data, &strData)
	if err != nil {
		return "", err
	}

	jsonString, err := json.Marshal(strData)
	if err != nil {
		return "", err
	}

	client.Set(ctx, "restaurants", jsonString, 10*time.Minute)

	fmt.Println("Keep restaurants to cache")

	return string(jsonString), nil
}

func RestaurantAcceptOrder(ctx echo.Context, kafkaProducer *adapters.KafkaProducer) (*core.AcceptOrderResponse, error) {
	var request core.AcceptOrderRequest
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

	kafkaMessage := fmt.Sprintf(`%s accepted at restaurant %s`, request.OrderId, request.RestaurantId)
	err = kafkaProducer.Publish("order_accepted", kafkaMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to publish message to Kafka: %w", err)
	}
	fmt.Println("Published order", request.OrderId, "to Kafka topic 'order_accepted'")

	response := &core.AcceptOrderResponse{
		Status: "accepted",
	}

	return response, nil
}
