package services

import (
	"errors"
	"fmt"

	"github.com/aboloredev/armory/internal/dto"
	"github.com/aboloredev/armory/internal/events"
	"github.com/aboloredev/armory/internal/models"
	"github.com/aboloredev/armory/internal/notifications"
	"github.com/aboloredev/armory/internal/utils"
	"gorm.io/gorm"
)

type OrderService struct {
	db             *gorm.DB
	eventPublisher events.Publisher
}

func NewOrderService(db *gorm.DB, eventPublisher events.Publisher) *OrderService {
	return &OrderService{
		db:             db,
		eventPublisher: eventPublisher,
	}
}

func (s *OrderService) CreateOrderService(userID uint) (*dto.OrderResponse, error) {
	var orderResponse *dto.OrderResponse

	err := s.db.Transaction(func(tx *gorm.DB) error {

		var cart models.Cart
		err := tx.Preload("CartItems.Product").Preload("User").Where("user_id = ? ", userID).First(&cart).Error
		if err != nil {
			return errors.New("cart not found")
		}

		if len(cart.CartItems) == 0 {
			return errors.New("cart is empty")
		}

		var totalAmount float64
		var orderItems []models.OrderItem

		for i := range cart.CartItems {
			cartItem := &cart.CartItems[i]

			var product models.Product
			err := tx.Set("gorm:query_option", "FOR UPDATE").
				First(&product, cartItem.ProductID).Error
			if err != nil {
				return ErrProductNotFound
			}

			result := tx.Model(&models.Product{}).
				Where("id = ? AND stock >= ?", product.ID, cartItem.Quantity).
				Update("stock", gorm.Expr("stock - ?", cartItem.Quantity))

			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return fmt.Errorf("insufficient stock for %s", product.Name)
			}

			itemTotal := cartItem.Product.Price * float64(cartItem.Quantity)
			totalAmount = totalAmount + itemTotal

			orderItems = append(orderItems, models.OrderItem{
				ProductID: cartItem.ProductID,
				Quantity:  cartItem.Quantity,
				Price:     cartItem.Product.Price,
			})
		}

		order := models.Order{
			UserID:      userID,
			Status:      models.OrderPending,
			TotalAmount: totalAmount,
			OrderItems:  orderItems,
			
		}

		err = tx.Create(&order).Error
		if err != nil {
			return err
		}

		err = tx.Unscoped().Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error
		if err != nil {
			return err
		}

		response, err := s.GetOrderResponse(tx, order.ID)
		if err != nil {
			return err
		}

		orderResponse = response

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Publish events for order
	err = s.eventPublisher.Publish(
		notifications.OrderCreatedSuccessfully,
		orderResponse,
		map[string]string{"Priority": "Please treat this order as priority"},
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to put messages in a queue %v", err)
	}

	return orderResponse, nil
}

func (s *OrderService) ConvertToOrderResponse(order *models.Order) *dto.OrderResponse {
	orderItems := make([]dto.OrderItemResponse, len(order.OrderItems))

	for i := range order.OrderItems {
		orderItems[i] = dto.OrderItemResponse{
			ID: order.OrderItems[i].ID,
			Product: dto.ProductResponse{
				ID:    order.OrderItems[i].Product.ID,
				Name:  order.OrderItems[i].Product.Name,
				SKU:   order.OrderItems[i].Product.SKU,
				Price: order.OrderItems[i].Product.Price,
			},
			Quantity: order.OrderItems[i].Quantity,
			Price:    order.OrderItems[i].Price,
		}
	}

	return &dto.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		UserEmail: order.User.Email,
		UserFirstName: order.User.FirstName,
		UserLastName: order.User.LastName,
		Status:      string(order.Status),
		TotalAmount: order.TotalAmount,
		OrderItems:  orderItems,
		CreatedAt:   order.CreatedAt,
	}
}

func (s *OrderService) GetOrderResponse(tx *gorm.DB, orderID uint) (*dto.OrderResponse, error) {
	var order models.Order

	err := tx.Preload("OrderItems.Product.Category").Preload("User").Where("id = ? ", orderID).First(&order).Error
	if err != nil {
		return nil, err
	}

	return s.ConvertToOrderResponse(&order), nil
}

func (s *OrderService) GetAllUserOrders(userID uint, page, pageSize int) ([]*dto.OrderResponse, *utils.PaginatedMeta, error) {

	var orders []models.Order
	var total int64

	offset := utils.PaginationHelper(page, pageSize)

	s.db.Model(&models.Order{}).Where("user_id = ? ", userID).Count(&total)

	err := s.db.Preload("OrderItems.Product.Category").
		Where("user_id = ? ", userID).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&orders).Error

	if err != nil {
		return nil, nil, err
	}

	response := make([]*dto.OrderResponse, len(orders))

	for _, order := range orders {
		response = append(response, s.ConvertToOrderResponse(&order))
	}

	totalPages := int(total + int64(pageSize-1)/int64(pageSize))

	meta := &utils.PaginatedMeta{
		Page:       page,
		Limit:      pageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	return response, meta, nil
}

func (s *OrderService) GetOrder(orderID, userID uint) (*dto.OrderResponse, error) {
	var order models.Order

	err := s.db.Preload("OrderItems.Product.Category").
		Where("id = ? AND user_id = ? ", orderID, userID).First(&order).Error
	if err != nil {
		return nil, err
	}

	return s.ConvertToOrderResponse(&order), nil
}
