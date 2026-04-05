package order

type Order struct {
	Id       int32  `json:"id,omitempty"`
	UserId   int32  `json:"user_id,omitempty"`
	Product  string `json:"product,omitempty"`
	Quantity int32  `json:"quantity,omitempty"`
	UserName string `json:"user_name,omitempty"`
}

type CreateOrderRequest struct {
	UserId   int32  `json:"user_id,omitempty"`
	Product  string `json:"product,omitempty"`
	Quantity int32  `json:"quantity,omitempty"`
}

type CreateOrderResponse struct {
	Order *Order `json:"order,omitempty"`
}

type GetOrderRequest struct {
	Id int32 `json:"id,omitempty"`
}

type GetOrderResponse struct {
	Order *Order `json:"order,omitempty"`
}
