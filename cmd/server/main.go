package main

import (
	"fmt"
	"log"
	"net/http"

	"food-delivery-app/internal/adapters"
	"food-delivery-app/internal/core"
	"food-delivery-app/internal/services"

	"github.com/go-redis/redis/v8"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	kafkaProducer, err := adapters.NewKafkaProducer([]string{"localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer func() {
		if err := kafkaProducer.Close(); err != nil {
			log.Printf("Error closing Kafka producer: %v", err)
		}
	}()

	e.GET("/menu", func(c echo.Context) error {
		resId := c.QueryParam("restaurant_id")
		if resId == "" {
			fmt.Println("Have restaurant_id is", resId)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "restaurant_id is required"})
		}

		menu, err := services.GetMenu(redisClient, resId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch menu" + err.Error()})

		}
		return c.String(http.StatusOK, menu)
	})

	e.GET("/restaurant", func(c echo.Context) error {
		restaurant, err := services.GetRestaurants(redisClient)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch restaurant" + err.Error()})

		}
		return c.String(http.StatusOK, restaurant)
	})

	e.GET("/rider", func(c echo.Context) error {
		rider, err := services.GetRiders(redisClient)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch rider" + err.Error()})

		}
		return c.String(http.StatusOK, rider)
	})

	e.POST("/order", func(c echo.Context) error {
		order, err := services.PlaceOrder(redisClient, c, kafkaProducer)
		if err != nil {
			fmt.Printf("Failed to place order: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "failed to place ordered"})
		}
		return c.JSON(http.StatusOK, order)
	})

	e.POST("/notification/send", func(c echo.Context) error {
		request := new(core.NotificationRequest)
		if err := c.Bind(request); err != nil {
			fmt.Printf("Failed to bind request: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "failed to parse request"})
		}
		order := services.SendNotification(*request)
		if err != nil {
			fmt.Printf("Failed to send notification: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "failed to send notification"})
		}

		return c.JSON(http.StatusOK, order)
	})

	e.POST("/restaurant/order/accept", func(c echo.Context) error {
		request := new(core.AcceptOrderRequest)
		if err := c.Bind(request); err != nil {
			fmt.Printf("Failed to bind request: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "failed to parse request"})
		}
		order, err := services.RestaurantAcceptOrder(c, kafkaProducer)
		if err != nil {
			fmt.Printf("Failed to accepted order: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "failed to accepted order"})
		}

		return c.JSON(http.StatusOK, order)
	})

	e.Logger.Fatal(e.Start(":8088"))
}
