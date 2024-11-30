package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-redis/redis/v8"
)

func GetMenu(client *redis.Client, resId string) (string, error) {
	var ctx = context.Background()

	//Get from cache
	cacheKey := "menus:" + resId
	menus, err := client.Get(ctx, cacheKey).Result()
	if err == nil {
		return string(menus), nil
	}

	//Get from files
	filePath := fmt.Sprintf("mock_data/menu_res%s.json", resId)
	data, err := ioutil.ReadFile(filePath)
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
	client.Set(ctx, cacheKey, jsonString, 10*time.Minute)

	return string(jsonString), nil
}
