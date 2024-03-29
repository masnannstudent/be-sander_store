package domain

type Status struct {
	PaymentStatus string
	OrderStatus   string
}

type CreatePaymentRequest struct {
	OrderID         string `json:"order_id"`
	TotalAmountPaid uint64 `json:"total_amount_paid"`
	Name            string `json:"name"`
	Email           string `json:"email"`
}

type CreateOrderRequest struct {
	AddressID uint64 `form:"address_id" json:"address_id" validate:"required"`
	Size      string `form:"size" json:"size"`
	Color     string `form:"color" json:"color"`
	Note      string `form:"note" json:"note"`
	ProductID uint64 `json:"product_id" validate:"required"`
	Quantity  uint64 `json:"quantity" validate:"required"`
}

type CreateNotificationOrderRequest struct {
	OrderStatus string `json:"order_status"`
	UserID      uint64 `json:"user_id"`
	OrderID     string `json:"order_id"`
	Title       string `json:"title"`
	Message     string `json:"message"`
}

type CreateNotificationPaymentRequest struct {
	PaymentStatus string `json:"payment_status"`
	UserID        uint64 `json:"user_id"`
	OrderID       string `json:"order_id"`
	Title         string `json:"title"`
	Message       string `json:"message"`
}

type CreateCartRequest struct {
	ProductID uint64 `json:"product_id" validate:"required"`
	Size      string `form:"size" json:"size"`
	Color     string `form:"color" json:"color"`
	Quantity  uint64 `json:"quantity" validate:"required,min=1"`
}

type CreateOrderCartRequest struct {
	AddressID uint64            `form:"address_id" json:"address_id" validate:"required"`
	Note      string            `form:"note" json:"note"`
	CartItems []CartItemRequest `json:"cart_items" validate:"required"`
}

type CartItemRequest struct {
	ID uint64 `json:"id" validate:"required"`
}

type UpdateOrderStatus struct {
	ID          string `json:"id" validate:"required"`
	OrderStatus string `json:"order_status" validate:"required"`
}
