package core

type Menu struct {
	ID    string
	Name  string
	Price string
}

type OrderDetail struct {
	OrderID      string `json:"order_id"`
	RestaurantID string `json:"restaurant_id"`
	Items        []OrderItem
	Status       string `json:"status"`
}

type OrderRequest struct {
	RestaurantID string `json:"restaurant_id"`
	Items        []OrderItem
}

type Order struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

type OrderItem struct {
	MenuID   string `json:"menu_id"`
	Quantity int    `json:"quantity"`
}

type NotificationRequest struct {
	Recipient string `json:"recipient"`
	OrderId   string `json:"order_id"`
	Message   string `json:"message"`
}

type NotificationResponse struct {
	Status string `json:"status"`
}

type AcceptOrderRequest struct {
	OrderId      string `json:"order_id"`
	RestaurantId string `json:"restaurant_id"`
}

type AcceptOrderResponse struct {
	Status string `json:"status"`
}

type PickupConfirmationRequest struct {
	OrderId string `json:"order_id"`
	RiderId string `json:"rider_id"`
}

type PickupConfirmationResponse struct {
	Status string `json:"status"`
}

type DeliveryOrderRequest struct {
	OrderId string `json:"order_id"`
	RiderId string `json:"rider_id"`
}

type DeliveryOrderResponse struct {
	Status string `json:"status"`
}

type Restaurant struct {
	ID   string
	Name string
	Menu []Menu
}
