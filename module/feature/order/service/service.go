package service

import (
	"errors"
	"fmt"
	"math"
	"ruti-store/module/entities"
	address "ruti-store/module/feature/address/domain"
	notification "ruti-store/module/feature/notification/domain"
	"ruti-store/module/feature/order/domain"
	product "ruti-store/module/feature/product/domain"
	users "ruti-store/module/feature/user/domain"
	"ruti-store/utils/generator"
	"time"
)

type OrderService struct {
	repo                domain.OrderRepositoryInterface
	generatorID         generator.GeneratorInterface
	productService      product.ProductServiceInterface
	addressService      address.AddressServiceInterface
	userService         users.UserServiceInterface
	notificationService notification.NotificationServiceInterface
}

func NewOrderService(
	repo domain.OrderRepositoryInterface,
	generatorID generator.GeneratorInterface,
	productService product.ProductServiceInterface,
	addressService address.AddressServiceInterface,
	userService users.UserServiceInterface,
	notificationService notification.NotificationServiceInterface,

) domain.OrderServiceInterface {
	return &OrderService{
		repo:                repo,
		generatorID:         generatorID,
		productService:      productService,
		addressService:      addressService,
		userService:         userService,
		notificationService: notificationService,
	}
}

func (s *OrderService) GetAllOrders(page, pageSize int) ([]*entities.OrderModels, int64, error) {
	result, err := s.repo.GetPaginatedOrders(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.repo.GetTotalItems()
	if err != nil {
		return nil, 0, err
	}

	return result, totalItems, nil
}

func (s *OrderService) GetOrdersPage(currentPage, pageSize, totalItems int) (int, int, int, error) {
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
	nextPage := currentPage + 1
	prevPage := currentPage - 1

	if nextPage > totalPages {
		nextPage = 0
	}

	if prevPage < 1 {
		prevPage = 0
	}

	return totalPages, nextPage, prevPage, nil
}

func (s *OrderService) CreateOrder(userID uint64, request *domain.CreateOrderRequest) (*domain.CreateOrderResponse, error) {
	orderID, err := s.generatorID.GenerateUUID()
	if err != nil {
		return nil, errors.New("failed to generate order ID")
	}

	idOrder, err := s.generatorID.GenerateOrderID()
	if err != nil {
		return nil, errors.New("failed to generate order ID")
	}

	addresses, err := s.addressService.GetAddressByID(request.AddressID)
	if err != nil {
		return nil, errors.New("address not found")
	}

	products, err := s.productService.GetProductByID(request.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	var orderDetails []entities.OrderDetailsModels
	var totalQuantity, totalPrice, totalDiscount uint64

	orderDetail := entities.OrderDetailsModels{
		OrderID:       orderID,
		ProductID:     request.ProductID,
		Size:          request.Size,
		Color:         request.Color,
		Quantity:      request.Quantity,
		IsReviewed:    false,
		TotalPrice:    request.Quantity * (products.Price - products.Discount),
		TotalDiscount: products.Discount * request.Quantity,
	}

	totalQuantity += request.Quantity
	totalPrice += orderDetail.TotalPrice
	totalDiscount += orderDetail.TotalDiscount

	orderDetails = append(orderDetails, orderDetail)

	grandTotalPrice := totalPrice
	totalAmountPaid := grandTotalPrice + 2000

	newData := &entities.OrderModels{
		ID:                 orderID,
		IdOrder:            idOrder,
		AddressID:          addresses.ID,
		UserID:             userID,
		Note:               request.Note,
		GrandTotalQuantity: totalQuantity,
		GrandTotalPrice:    grandTotalPrice,
		ShipmentFee:        0,
		AdminFees:          2000,
		TotalAmountPaid:    totalAmountPaid,
		OrderStatus:        "Menunggu Konfirmasi",
		PaymentStatus:      "Menunggu Konfirmasi",
		CreatedAt:          time.Now(),
		OrderDetails:       orderDetails,
	}

	createdOrder, err := s.repo.CreateOrder(newData)
	if err != nil {
		return nil, err
	}

	user, err := s.userService.GetUserByID(createdOrder.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	snapResult, err := s.repo.CreateSnap(createdOrder.ID, user.Name, user.Email, createdOrder.TotalAmountPaid)
	if err != nil {
		return nil, err
	}

	for _, variant := range products.Variants {
		if variant.Size == request.Size && variant.Color == request.Color {
			if err := s.productService.ReduceStockWhenPurchasing(variant.ID, request.Quantity); err != nil {
				return nil, errors.New("failed reduce stock for variant")
			}
			break
		}
	}

	notificationRequest := domain.CreateNotificationPaymentRequest{
		OrderID:       createdOrder.ID,
		UserID:        createdOrder.UserID,
		PaymentStatus: "Menunggu Konfirmasi",
	}
	_, err = s.SendNotificationPayment(notificationRequest)
	if err != nil {
		return nil, err
	}

	response := &domain.CreateOrderResponse{
		OrderID:         createdOrder.ID,
		IdOrder:         createdOrder.IdOrder,
		RedirectURL:     snapResult.RedirectURL,
		TotalAmountPaid: createdOrder.TotalAmountPaid,
	}
	return response, nil
}

func (s *OrderService) CallBack(req map[string]interface{}) error {
	orderID, exist := req["order_id"].(string)
	if !exist {
		return errors.New("invalid notification payload")
	}

	status, err := s.repo.CheckTransaction(orderID)
	if err != nil {
		return err
	}

	transaction, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return errors.New("transaction data not found")
	}

	if status.PaymentStatus == "Konfirmasi" {
		if err := s.ConfirmPayment(transaction.ID); err != nil {
			return err
		}
	} else if status.PaymentStatus == "Gagal" {
		if err := s.CancelPayment(transaction.ID); err != nil {
			return err
		}
	}

	return nil
}

func (s *OrderService) GetOrderByID(orderID string) (*entities.OrderModels, error) {
	result, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *OrderService) ConfirmPayment(orderID string) error {
	orders, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return errors.New("transaction data not found")
	}

	orders.OrderStatus = "Proses"
	orders.PaymentStatus = "Konfirmasi"

	if err := s.repo.UpdatePayment(orders.ID, orders.OrderStatus, orders.PaymentStatus); err != nil {
		return err
	}
	notificationRequest := domain.CreateNotificationPaymentRequest{
		OrderID:       orders.ID,
		UserID:        orders.UserID,
		PaymentStatus: "Konfirmasi",
	}
	_, err = s.SendNotificationPayment(notificationRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) CancelPayment(orderID string) error {
	orders, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return errors.New("transaction data not found")
	}

	orders.OrderStatus = "Gagal"
	orders.PaymentStatus = "Gagal"

	if err := s.repo.UpdatePayment(orderID, orders.OrderStatus, orders.PaymentStatus); err != nil {
		return err
	}

	for _, orderDetail := range orders.OrderDetails {
		if err := s.productService.IncreaseStock(orderDetail.ProductID, orderDetail.Quantity); err != nil {
			return errors.New("failed to increase stock")
		}
	}
	notificationRequest := domain.CreateNotificationPaymentRequest{
		OrderID:       orders.ID,
		UserID:        orders.UserID,
		PaymentStatus: "Menunggu Konfirmasi",
	}
	_, err = s.SendNotificationPayment(notificationRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) SendNotificationPayment(request domain.CreateNotificationPaymentRequest) (string, error) {
	var notificationMsg string
	var err error

	user, err := s.userService.GetUserByID(request.UserID)
	if err != nil {
		return "", err
	}
	orders, err := s.repo.GetOrderByID(request.OrderID)
	if err != nil {
		return "", err
	}

	switch request.PaymentStatus {
	case "Menunggu Konfirmasi":
		notificationMsg = fmt.Sprintf("Halo, %s! Pesanan dengan ID %s sudah berhasil dibuat. Harap ditunggu!", user.Name, orders.IdOrder)
	case "Konfirmasi":
		notificationMsg = fmt.Sprintf("Terima kasih, %s! Pembayaran untuk pesanan dengan ID %s telah kami terima. Semoga harimu menyenangkan!", user.Name, orders.IdOrder)
	case "Gagal":
		notificationMsg = fmt.Sprintf("Maaf, %s. Pembayaran untuk pesanan dengan ID %s gagal. Beritahu kami jika Anda membutuhkan bantuan!", user.Name, orders.IdOrder)
	default:
		return "", errors.New("Status pesanan tidak valid")
	}
	req := &notification.CreateNotificationRequest{
		UserID:  user.ID,
		OrderID: orders.IdOrder,
		Title:   "Status Pembayaran",
		Message: notificationMsg,
	}
	_, err = s.notificationService.CreateNotification(req)
	if err != nil {
		return "", errors.New("error send message")
	}

	return notificationMsg, nil
}

func (s *OrderService) SendNotificationOrder(request domain.CreateNotificationOrderRequest) (string, error) {
	var notificationMsg string
	var err error

	user, err := s.userService.GetUserByID(request.UserID)
	if err != nil {
		return "", err
	}
	orders, err := s.repo.GetOrderByID(request.OrderID)
	if err != nil {
		return "", err
	}

	switch request.OrderStatus {
	case "Pengiriman":
		notificationMsg = fmt.Sprintf("Halo, %s! Pesanan dengan ID %s sedang dalam proses pengiriman. Harap ditunggu!", user.Name, orders.IdOrder)
	case "Selesai":
		notificationMsg = fmt.Sprintf("Selamat, %s! Pesanan dengan ID %s sudah sampai tujuan. Semoga Anda puas!", user.Name, orders.IdOrder)
	case "Menunggu Konfirmasi":
		notificationMsg = fmt.Sprintf("Halo, %s! Pesanan dengan ID %s sedang menunggu konfirmasi. Harap ditunggu!", user.Name, orders.IdOrder)
	case "Proses":
		notificationMsg = fmt.Sprintf("Halo, %s! Pesanan dengan ID %s sedang dalam proses. Harap ditunggu!", user.Name, orders.IdOrder)
	case "Gagal":
		notificationMsg = fmt.Sprintf("Maaf, %s. Pesanan dengan ID %s gagal. Silakan coba lagi.", user.Name, orders.IdOrder)
	default:
		return "", errors.New("Status pengiriman tidak valid")
	}

	req := &notification.CreateNotificationRequest{
		UserID:  user.ID,
		OrderID: orders.IdOrder,
		Title:   "Status Pesanan",
		Message: notificationMsg,
	}
	_, err = s.notificationService.CreateNotification(req)
	if err != nil {
		return "", errors.New("error send message")
	}

	return notificationMsg, nil
}

func (s *OrderService) CreateCart(userID uint64, req *domain.CreateCartRequest) (*entities.CartModels, error) {

	products, err := s.productService.GetProductByID(req.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	existingCartItem, err := s.repo.GetCartItem(user.ID, products.ID)
	if err == nil && existingCartItem != nil {
		existingCartItem.Quantity += req.Quantity

		err := s.repo.UpdateCartItem(existingCartItem)
		if err != nil {
			return nil, errors.New("gagal mengubah jumlah produk di keranjang")
		}

		return existingCartItem, nil
	}

	newData := &entities.CartModels{
		UserID:    user.ID,
		ProductID: products.ID,
		Size:      req.Size,
		Color:     req.Color,
		Quantity:  req.Quantity,
	}

	result, err := s.repo.CreateCart(newData)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (s *OrderService) DeleteCartItems(cartID uint64) error {
	cart, err := s.repo.GetCartByID(cartID)
	if err != nil {
		return err
	}

	err = s.repo.DeleteCartItem(cart.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) GetCartUser(userID uint64) ([]*entities.CartModels, error) {
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	result, err := s.repo.GetCartByUserID(user.ID)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (s *OrderService) CreateOrderCart(userID uint64, request *domain.CreateOrderCartRequest) (*domain.CreateOrderResponse, error) {
	orderID, err := s.generatorID.GenerateUUID()
	if err != nil {
		return nil, errors.New("failed to generate order ID")
	}

	idOrder, err := s.generatorID.GenerateOrderID()
	if err != nil {
		return nil, errors.New("failed to generate order ID")
	}

	addresses, err := s.addressService.GetAddressByID(request.AddressID)
	if err != nil {
		return nil, errors.New("address not found")
	}

	var orderDetails []entities.OrderDetailsModels
	var totalQuantity, totalPrice, totalDiscount uint64

	for _, cartItemRequest := range request.CartItems {

		cartItem, err := s.repo.GetCartByID(cartItemRequest.ID)
		if err != nil {
			return nil, errors.New("cart item not found")
		}

		products, err := s.productService.GetProductByID(cartItem.ProductID)
		if err != nil {
			return nil, errors.New("product not found")
		}

		orderDetail := entities.OrderDetailsModels{
			OrderID:       orderID,
			ProductID:     products.ID,
			Size:          cartItem.Size,
			Color:         cartItem.Color,
			Quantity:      cartItem.Quantity,
			IsReviewed:    false,
			TotalPrice:    cartItem.Quantity * (products.Price - products.Discount),
			TotalDiscount: products.Discount * cartItem.Quantity,
		}

		totalQuantity += cartItem.Quantity
		totalPrice += orderDetail.TotalPrice
		totalDiscount += orderDetail.TotalDiscount

		orderDetails = append(orderDetails, orderDetail)

		for _, variant := range products.Variants {
			if variant.Size == cartItem.Size && variant.Color == cartItem.Color {
				if err := s.productService.ReduceStockWhenPurchasing(variant.ID, cartItem.Quantity); err != nil {
					return nil, errors.New("failed reduce stock for variant")
				}
				break
			}
		}
	}

	grandTotalPrice := totalPrice
	totalAmountPaid := grandTotalPrice + 2000

	newData := &entities.OrderModels{
		ID:                 orderID,
		IdOrder:            idOrder,
		AddressID:          addresses.ID,
		UserID:             userID,
		Note:               request.Note,
		GrandTotalQuantity: totalQuantity,
		GrandTotalPrice:    grandTotalPrice,
		ShipmentFee:        0,
		AdminFees:          2000,
		TotalAmountPaid:    totalAmountPaid,
		OrderStatus:        "Menunggu Konfirmasi",
		PaymentStatus:      "Menunggu Konfirmasi",
		CreatedAt:          time.Now(),
		OrderDetails:       orderDetails,
	}

	createdOrder, err := s.repo.CreateOrder(newData)
	if err != nil {
		return nil, err
	}

	user, err := s.userService.GetUserByID(createdOrder.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	snapResult, err := s.repo.CreateSnap(createdOrder.ID, user.Name, user.Email, createdOrder.TotalAmountPaid)
	if err != nil {
		return nil, err
	}

	notificationRequest := domain.CreateNotificationPaymentRequest{
		OrderID:       createdOrder.ID,
		UserID:        createdOrder.UserID,
		PaymentStatus: "Menunggu Konfirmasi",
	}
	_, err = s.SendNotificationPayment(notificationRequest)
	if err != nil {
		return nil, err
	}

	for _, cartItemRequest := range request.CartItems {
		if err := s.repo.DeleteCartItem(cartItemRequest.ID); err != nil {
			return nil, errors.New("failed to delete cart item")
		}
	}

	response := &domain.CreateOrderResponse{
		OrderID:         createdOrder.ID,
		IdOrder:         createdOrder.IdOrder,
		RedirectURL:     snapResult.RedirectURL,
		TotalAmountPaid: createdOrder.TotalAmountPaid,
	}
	return response, nil
}

func (s *OrderService) AcceptOrder(orderID string) error {
	orders, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	user, err := s.userService.GetUserByID(orders.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	orders.OrderStatus = "Selesai"

	if err := s.repo.AcceptOrder(orders.ID, orders.OrderStatus); err != nil {
		return err
	}
	notificationRequest := domain.CreateNotificationOrderRequest{
		OrderID:     orders.ID,
		UserID:      user.ID,
		OrderStatus: "Selesai",
	}
	_, err = s.SendNotificationOrder(notificationRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) UpdateOrderStatus(req *domain.UpdateOrderStatus) error {
	orders, err := s.repo.GetOrderByID(req.ID)
	if err != nil {
		return errors.New("order not found")
	}

	if err := s.repo.UpdateOrderStatus(orders.ID, req.OrderStatus); err != nil {
		return err
	}

	user, err := s.userService.GetUserByID(orders.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	notificationRequest := domain.CreateNotificationOrderRequest{
		OrderID:     orders.ID,
		UserID:      user.ID,
		OrderStatus: req.OrderStatus,
	}
	_, err = s.SendNotificationOrder(notificationRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) GetAllOrdersByUserID(userID uint64, page, pageSize int) ([]*entities.OrderModels, int64, error) {
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return nil, 0, errors.New("user not found")
	}

	result, totalItems, err := s.repo.GetAllOrdersByUserID(user.ID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return result, totalItems, nil
}

func (s *OrderService) GetCartById(cartID uint64) (*entities.CartModels, error) {
	result, err := s.repo.GetCartByID(cartID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *OrderService) GetAllOrdersWithFilter(userID uint64, orderStatus string, page, pageSize int) ([]*entities.OrderModels, int64, error) {
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return nil, 0, errors.New("user not found")
	}

	result, totalItems, err := s.repo.GetAllOrdersUserWithFilter(user.ID, orderStatus, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return result, totalItems, nil
}

func (s *OrderService) SearchAndPaginateOrder(page, pageSize int, name string) ([]*entities.OrderModels, int64, error) {
	result, totalItems, err := s.repo.GetAllOrdersSearch(page, pageSize, name)
	if err != nil {
		return nil, 0, err
	}
	return result, totalItems, nil
}

func (s *OrderService) GetReportOrder(starDate, endDate time.Time) ([]*entities.OrderModels, error) {
	result, err := s.repo.GetReportOrder(starDate, endDate)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *OrderService) FilterAndPaginatePayment(page, pageSize int, filter string) ([]*entities.OrderModels, int64, error) {
	result, totalItems, err := s.repo.GetAllPaymentFilter(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	return result, totalItems, nil
}

func (s *OrderService) SearchFilterAndPaginatePayment(page, pageSize int, name, filter string) ([]*entities.OrderModels, int64, error) {
	result, totalItems, err := s.repo.GetAllPaymentFilterAndSearch(page, pageSize, name, filter)
	if err != nil {
		return nil, 0, err
	}
	return result, totalItems, nil
}

func (s *OrderService) FilterAndPaginateOrder(page, pageSize int, filter string) ([]*entities.OrderModels, int64, error) {
	result, totalItems, err := s.repo.GetAllOrderFilter(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	return result, totalItems, nil
}

func (s *OrderService) SearchFilterAndPaginateOrder(page, pageSize int, name, filter string) ([]*entities.OrderModels, int64, error) {
	result, totalItems, err := s.repo.GetAllOrderFilterAndSearch(page, pageSize, name, filter)
	if err != nil {
		return nil, 0, err
	}
	return result, totalItems, nil
}
