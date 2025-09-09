package model

type CreateOrderRequest struct {
	CustomerID    string `json:"customer_id" validate:"required,uuid"`
	CustomerName  string `json:"customer_name" validate:"required"`
	CustomerEmail string `json:"customer_email" validate:"required,email"`
	CustomerPhone string `json:"customer_phone" validate:"required"`

	Items []OrderItemRequest `json:"items" validate:"required,dive"`
	Notes string             `json:"notes,omitempty"`

	PaymentMethod   string `json:"payment_method" validate:"required"`
	ShippingAddress string `json:"shipping_address,omitempty"`
}

type OrderItemRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	Name      string `json:"name" validate:"required"`
	Price     int64  `json:"price" validate:"required,gt=0"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
}
