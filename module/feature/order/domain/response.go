package domain

import (
	"ruti-store/module/entities"
	"time"
)

type OrderResponse struct {
	ID                 string                `json:"id"`
	IdOrder            string                `json:"id_order"`
	AddressID          uint64                `json:"address_id"`
	UserID             uint64                `json:"user_id"`
	Note               string                `json:"note"`
	GrandTotalQuantity uint64                `json:"grand_total_quantity"`
	GrandTotalPrice    uint64                `json:"grand_total_price"`
	ShipmentFee        uint64                `json:"shipment_fee"`
	AdminFees          uint64                `json:"admin_fees"`
	GrandTotalDiscount uint64                `json:"grand_total_discount"`
	TotalAmountPaid    uint64                `json:"total_amount_paid"`
	OrderStatus        string                `json:"order_status"`
	PaymentStatus      string                `json:"payment_status"`
	CreatedAt          time.Time             `json:"created_at"`
	Address            AddressResponse       `json:"address"`
	User               UserResponse          `json:"user"`
	OrderDetails       []OrderDetailResponse `json:"order_details"`
}

type OrderDetailResponse struct {
	ID            uint64          `json:"id"`
	OrderID       string          `json:"order_id"`
	ProductID     uint64          `json:"product_id"`
	IsReviewed    bool            `json:"is_reviewed"`
	Size          string          `json:"size"`
	Color         string          `json:"color"`
	Quantity      uint64          `json:"quantity"`
	TotalPrice    uint64          `json:"total_price"`
	TotalDiscount uint64          `json:"total_discount"`
	Product       ProductResponse `json:"product,omitempty"`
}

type ProductPhotoResponse struct {
	ID        uint64 `json:"id"`
	ProductID uint64 `json:"product_id"`
	URL       string `json:"url"`
}

type ProductResponse struct {
	ID            uint64                 `json:"id"`
	Name          string                 `json:"name"`
	Price         uint64                 `json:"price"`
	Discount      uint64                 `json:"discount"`
	ProductPhotos []ProductPhotoResponse `json:"product_photos"`
}

type AddressResponse struct {
	ID           uint64 `json:"id"`
	UserID       uint64 `json:"user_id"`
	AcceptedName string `json:"accepted_name" `
	Phone        string `json:"phone"`
	ProvinceName string `json:"province_name"`
	CityName     string `json:"city_name"`
	Address      string `json:"address"`
	IsPrimary    bool   `json:"is_primary"`
}

type UserResponse struct {
	ID           uint64 `json:"id"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Name         string `json:"name"`
	PhotoProfile string `json:"photo_profile"`
}

func FormatOrderDetail(order *entities.OrderModels) OrderResponse {
	orderResponse := OrderResponse{
		ID:                 order.ID,
		IdOrder:            order.IdOrder,
		AddressID:          order.AddressID,
		UserID:             order.UserID,
		Note:               order.Note,
		GrandTotalQuantity: order.GrandTotalQuantity,
		GrandTotalPrice:    order.GrandTotalPrice,
		ShipmentFee:        order.ShipmentFee,
		AdminFees:          order.AdminFees,
		GrandTotalDiscount: order.GrandTotalDiscount,
		TotalAmountPaid:    order.TotalAmountPaid,
		OrderStatus:        order.OrderStatus,
		PaymentStatus:      order.PaymentStatus,
		CreatedAt:          order.CreatedAt,
		Address: AddressResponse{
			ID:           order.Address.ID,
			UserID:       order.Address.UserID,
			AcceptedName: order.Address.AcceptedName,
			Phone:        order.Address.Phone,
			ProvinceName: order.Address.ProvinceName,
			CityName:     order.Address.CityName,
			Address:      order.Address.Address,
			IsPrimary:    order.Address.IsPrimary,
		},
		User: UserResponse{
			ID:           order.User.ID,
			Email:        order.User.Email,
			Phone:        order.User.Phone,
			Name:         order.User.Name,
			PhotoProfile: order.User.PhotoProfile,
		},
	}

	var orderDetails []OrderDetailResponse
	for _, detail := range order.OrderDetails {
		var productPhotos []ProductPhotoResponse
		for _, photo := range detail.Product.Photos {
			productPhotos = append(productPhotos, ProductPhotoResponse{
				ID:        photo.ID,
				ProductID: photo.ProductID,
				URL:       photo.URL,
			})
		}

		orderDetail := OrderDetailResponse{
			ID:            detail.ID,
			OrderID:       detail.OrderID,
			ProductID:     detail.ProductID,
			IsReviewed:    detail.IsReviewed,
			Size:          detail.Size,
			Color:         detail.Color,
			Quantity:      detail.Quantity,
			TotalPrice:    detail.TotalPrice,
			TotalDiscount: detail.TotalDiscount,
			Product: ProductResponse{
				ID:            detail.Product.ID,
				Name:          detail.Product.Name,
				Price:         detail.Product.Price,
				ProductPhotos: productPhotos,
			},
		}
		if len(detail.Product.Photos) > 0 {
			productPhoto := ProductPhotoResponse{
				ID:        detail.Product.Photos[0].ID,
				ProductID: detail.Product.Photos[0].ProductID,
				URL:       detail.Product.Photos[0].URL,
			}
			orderDetail.Product.ProductPhotos = []ProductPhotoResponse{productPhoto}
		}
		orderDetails = append(orderDetails, orderDetail)
	}

	orderResponse.OrderDetails = orderDetails

	return orderResponse
}

// OrderPaginationResponse Pagination Response
type OrderPaginationResponse struct {
	ID              string                      `json:"id"`
	IdOrder         string                      `json:"id_order"`
	UserID          uint64                      `json:"user_id"`
	TotalAmountPaid uint64                      `json:"total_amount_paid"`
	OrderStatus     string                      `json:"order_status"`
	CreatedAt       time.Time                   `json:"created_at"`
	User            UserPaginationOrderResponse `json:"user"`
}

type UserPaginationOrderResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func FormatOrderPagination(order *entities.OrderModels) *OrderPaginationResponse {
	orderResponse := &OrderPaginationResponse{
		ID:              order.ID,
		IdOrder:         order.IdOrder,
		UserID:          order.UserID,
		TotalAmountPaid: order.TotalAmountPaid,
		OrderStatus:     order.OrderStatus,
		CreatedAt:       order.CreatedAt,
		User: UserPaginationOrderResponse{
			ID:   order.User.ID,
			Name: order.User.Name,
		},
	}
	return orderResponse
}

func FormatterOrder(orders []*entities.OrderModels) []*OrderPaginationResponse {
	var orderFormatters []*OrderPaginationResponse

	for _, order := range orders {
		formattedOrder := FormatOrderPagination(order)
		orderFormatters = append(orderFormatters, formattedOrder)
	}

	return orderFormatters
}

// OrderSummaryResponse Get All Order
type OrderSummaryResponse struct {
	ID              string    `json:"id"`
	IDOrder         string    `json:"id_order"`
	Name            string    `json:"name"`
	Date            time.Time `json:"date"`
	TotalAmountPaid uint64    `json:"total_amount_paid"`
	OrderStatus     string    `json:"order_status"`
}

func ResponseArrayOrderSummary(data []*entities.OrderModels) []*OrderSummaryResponse {
	res := make([]*OrderSummaryResponse, 0)

	for _, order := range data {
		orderRes := &OrderSummaryResponse{
			ID:              order.ID,
			IDOrder:         order.IdOrder,
			Name:            order.User.Name,
			Date:            order.CreatedAt,
			TotalAmountPaid: order.TotalAmountPaid,
			OrderStatus:     order.OrderStatus,
		}
		res = append(res, orderRes)
	}

	return res
}

// PaymentSummaryResponse GetAll Payment
type PaymentSummaryResponse struct {
	ID              string    `json:"id"`
	IDOrder         string    `json:"id_order"`
	Name            string    `json:"name"`
	Date            time.Time `json:"date"`
	TotalAmountPaid uint64    `json:"total_amount_paid"`
	PaymentStatus   string    `json:"payment_status"`
}

func ResponseArrayPaymentSummary(data []*entities.OrderModels) []*PaymentSummaryResponse {
	res := make([]*PaymentSummaryResponse, 0)

	for _, order := range data {
		orderRes := &PaymentSummaryResponse{
			ID:              order.ID,
			IDOrder:         order.IdOrder,
			Name:            order.User.Name,
			Date:            order.CreatedAt,
			TotalAmountPaid: order.TotalAmountPaid,
			PaymentStatus:   order.PaymentStatus,
		}
		res = append(res, orderRes)
	}

	return res
}

type CreateOrderResponse struct {
	OrderID         string `json:"order_id"`
	IdOrder         string `json:"id_order"`
	RedirectURL     string `json:"redirect_url"`
	TotalAmountPaid uint64 `json:"total_amount_paid"`
}

type CreateCartResponse struct {
	ID        uint64 `json:"id"`
	UserID    uint64 `json:"user_id"`
	ProductID uint64 `json:"product_id"`
	Size      string `json:"size"`
	Color     string `json:"color"`
	Quantity  uint64 `json:"quantity"`
}

func CreateCartFormatter(cart *entities.CartModels) *CreateCartResponse {
	return &CreateCartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		ProductID: cart.ProductID,
		Size:      cart.Size,
		Color:     cart.Color,
		Quantity:  cart.Quantity,
	}
}

type CartResponse struct {
	ID        uint64           `json:"id"`
	UserID    uint64           `json:"user_id"`
	ProductID uint64           `json:"product_id"`
	Size      string           `json:"size"`
	Color     string           `json:"color"`
	Quantity  uint64           `json:"quantity"`
	Product   *ProductResponse `json:"product"`
}

func buildProductResponse(product *entities.ProductModels) *ProductResponse {
	return &ProductResponse{
		ID:            product.ID,
		Name:          product.Name,
		Price:         product.Price,
		Discount:      product.Discount,
		ProductPhotos: buildProductPhotoResponses(product.Photos),
	}
}

func buildProductPhotoResponses(photos []entities.ProductPhotoModels) []ProductPhotoResponse {
	photoResponses := make([]ProductPhotoResponse, len(photos))
	for i, photo := range photos {
		photoResponses[i] = ProductPhotoResponse{
			ID:        photo.ID,
			ProductID: photo.ProductID,
			URL:       photo.URL,
		}
	}
	return photoResponses
}
func CartFormatter(cart *entities.CartModels) *CartResponse {
	return &CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		ProductID: cart.ProductID,
		Size:      cart.Size,
		Color:     cart.Color,
		Quantity:  cart.Quantity,
		Product:   buildProductResponse(&cart.Product),
	}
}

func ResponseArrayCart(data []*entities.CartModels) []*CartResponse {
	res := make([]*CartResponse, len(data))

	for i, cart := range data {
		res[i] = &CartResponse{
			ID:        cart.ID,
			UserID:    cart.UserID,
			ProductID: cart.ProductID,
			Size:      cart.Size,
			Color:     cart.Color,
			Quantity:  cart.Quantity,
			Product:   buildProductResponse(&cart.Product),
		}
	}

	return res
}

// GetAllOrderUserResponse Respon to Get Order By UserID
type GetAllOrderUserResponse struct {
	ID              string                `json:"id"`
	IdOrder         string                `json:"id_order"`
	UserID          uint64                `json:"user_id"`
	Note            string                `json:"note"`
	TotalAmountPaid uint64                `json:"total_amount_paid"`
	OrderStatus     string                `json:"order_status"`
	PaymentStatus   string                `json:"payment_status"`
	CreatedAt       time.Time             `json:"created_at"`
	OrderDetails    []OrderDetailResponse `json:"order_details"`
}

func FormatGetAllOrderUser(order *entities.OrderModels) *GetAllOrderUserResponse {
	orderResponse := &GetAllOrderUserResponse{
		ID:              order.ID,
		IdOrder:         order.IdOrder,
		UserID:          order.UserID,
		Note:            order.Note,
		TotalAmountPaid: order.TotalAmountPaid,
		OrderStatus:     order.OrderStatus,
		PaymentStatus:   order.PaymentStatus,
		CreatedAt:       order.CreatedAt,
	}

	var orderDetails []OrderDetailResponse
	for _, detail := range order.OrderDetails {
		var productPhotos []ProductPhotoResponse
		for _, photo := range detail.Product.Photos {
			productPhotos = append(productPhotos, ProductPhotoResponse{
				ID:        photo.ID,
				ProductID: photo.ProductID,
				URL:       photo.URL,
			})
		}

		orderDetail := OrderDetailResponse{
			ID:            detail.ID,
			OrderID:       detail.OrderID,
			ProductID:     detail.ProductID,
			IsReviewed:    detail.IsReviewed,
			Size:          detail.Size,
			Color:         detail.Color,
			Quantity:      detail.Quantity,
			TotalPrice:    detail.TotalPrice,
			TotalDiscount: detail.TotalDiscount,
			Product: ProductResponse{
				ID:            detail.Product.ID,
				Name:          detail.Product.Name,
				Price:         detail.Product.Price,
				Discount:      detail.Product.Discount,
				ProductPhotos: productPhotos,
			},
		}
		if len(detail.Product.Photos) > 0 {
			productPhoto := ProductPhotoResponse{
				ID:        detail.Product.Photos[0].ID,
				ProductID: detail.Product.Photos[0].ProductID,
				URL:       detail.Product.Photos[0].URL,
			}
			orderDetail.Product.ProductPhotos = []ProductPhotoResponse{productPhoto}
		}
		orderDetails = append(orderDetails, orderDetail)
	}

	orderResponse.OrderDetails = orderDetails

	return orderResponse
}

func FormatterGetAllOrderUser(orders []*entities.OrderModels) []*GetAllOrderUserResponse {
	var orderFormatters []*GetAllOrderUserResponse

	for _, order := range orders {
		formattedOrder := FormatGetAllOrderUser(order)
		orderFormatters = append(orderFormatters, formattedOrder)
	}

	return orderFormatters
}

func ResponseArrayOrderUser(data []*entities.OrderModels) []*GetAllOrderUserResponse {
	res := make([]*GetAllOrderUserResponse, 0)

	for _, order := range data {
		orderRes := &GetAllOrderUserResponse{
			ID:              order.ID,
			IdOrder:         order.IdOrder,
			UserID:          order.UserID,
			Note:            order.Note,
			TotalAmountPaid: order.TotalAmountPaid,
			OrderStatus:     order.OrderStatus,
			PaymentStatus:   order.PaymentStatus,
			CreatedAt:       order.CreatedAt,
			OrderDetails:    getOrderDetailResponses(order.OrderDetails),
		}
		res = append(res, orderRes)
	}

	return res
}

func getOrderDetailResponses(data []entities.OrderDetailsModels) []OrderDetailResponse {
	res := make([]OrderDetailResponse, 0)

	for _, detail := range data {
		productPhotos := buildProductPhotoResponses(detail.Product.Photos)

		orderDetail := OrderDetailResponse{
			ID:            detail.ID,
			OrderID:       detail.OrderID,
			ProductID:     detail.ProductID,
			Size:          detail.Size,
			Color:         detail.Color,
			Quantity:      detail.Quantity,
			TotalPrice:    detail.TotalPrice,
			TotalDiscount: detail.TotalDiscount,
			Product: ProductResponse{
				ID:            detail.Product.ID,
				Name:          detail.Product.Name,
				Price:         detail.Product.Price,
				Discount:      detail.Product.Discount,
				ProductPhotos: productPhotos,
			},
		}
		if len(detail.Product.Photos) > 0 {
			productPhoto := ProductPhotoResponse{
				ID:  detail.Product.Photos[0].ID,
				URL: detail.Product.Photos[0].URL,
			}
			orderDetail.Product.ProductPhotos = []ProductPhotoResponse{productPhoto}
		}
		res = append(res, orderDetail)
	}

	return res
}

type OrderReportResponse struct {
	IdOrder            string                  `json:"id_order"`
	UserID             uint64                  `json:"user_id"`
	Note               string                  `json:"note"`
	GrandTotalQuantity uint64                  `json:"grand_total_quantity"`
	GrandTotalPrice    uint64                  `json:"grand_total_price"`
	GrandTotalDiscount uint64                  `json:"grand_total_discount"`
	TotalAmountPaid    uint64                  `json:"total_amount_paid"`
	OrderStatus        string                  `json:"order_status"`
	PaymentStatus      string                  `json:"payment_status"`
	CreatedAt          time.Time               `json:"created_at"`
	User               OrderReportUserResponse `json:"user"`
}

type OrderReportUserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func OrderReportFormatter(order *entities.OrderModels) *OrderReportResponse {
	return &OrderReportResponse{
		IdOrder:            order.IdOrder,
		UserID:             order.UserID,
		Note:               order.Note,
		GrandTotalQuantity: order.GrandTotalQuantity,
		GrandTotalPrice:    order.GrandTotalPrice,
		GrandTotalDiscount: order.GrandTotalDiscount,
		TotalAmountPaid:    order.TotalAmountPaid,
		OrderStatus:        order.OrderStatus,
		PaymentStatus:      order.PaymentStatus,
		CreatedAt:          order.CreatedAt,
		User: OrderReportUserResponse{
			Name:  order.User.Name,
			Email: order.User.Email,
		},
	}
}

func ResponseArrayOrderReport(data []*entities.OrderModels) []*OrderReportResponse {
	res := make([]*OrderReportResponse, len(data))

	for i, order := range data {
		res[i] = OrderReportFormatter(order)
	}

	return res
}
