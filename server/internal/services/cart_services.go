package services

import (
	"errors"

	"github.com/aboloredev/armory/internal/dto"
	"github.com/aboloredev/armory/internal/models"
	"gorm.io/gorm"
)

var ErrProductNotFound = errors.New("Product not found")
var ErrInsufficientStock = errors.New("Insufficient stock")
var ErrCartItemNotFound = errors.New("cart item not found")
var ErrCartNotFound = errors.New("cart not found")

type CartServices struct {
	db *gorm.DB
}

func NewCartServices(db *gorm.DB) *CartServices {
	return &CartServices{
		db: db,
	}
}

func (s *CartServices) ConvertToCartResponse(cart *models.Cart) *dto.CartResponse {
	cartItem := make([]dto.CartItemResponse, len(cart.CartItems))
	var total float64

	for i := range cart.CartItems {
		subtotal := cart.CartItems[i].Product.Price * float64(cart.CartItems[i].Quantity)
		total = total + subtotal
		cartItem[i] = dto.CartItemResponse{
			ID: cart.CartItems[i].ID,
			Product: dto.ProductResponse{
				ID:         cart.CartItems[i].ProductID,
				CategoryID: cart.CartItems[i].Product.CategoryID,
				Name:       cart.CartItems[i].Product.Name,
				Price:      cart.CartItems[i].Product.Price,
				SKU:        cart.CartItems[i].Product.SKU,
				Category: dto.CategoryResponse{
					ID:        cart.CartItems[i].Product.Category.ID,
					Name:      cart.CartItems[i].Product.Category.Name,
					CreatedAt: cart.CartItems[i].Product.Category.CreatedAt,
				},
			},
			Subtotal:  subtotal,
			Quantity:  cart.CartItems[i].Quantity,
			CreatedAt: cart.CartItems[i].CreatedAt,
		}
	}

	return &dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CartItems: cartItem,
		Total:     total,
		CreatedAt: cart.CreatedAt,
	}
}

func (s *CartServices) GetUserCart(userID uint) (*dto.CartResponse, error) {
	var cart models.Cart
	err := s.db.Preload("CartItems.Product.Category").Where("user_id = ? ", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}

	return s.ConvertToCartResponse(&cart), nil
}

func (s *CartServices) AddToCartSevices(userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error) {
	var product models.Product
	err := s.db.Where("id = ? ", req.ProductID).First(&product).Error
	if err != nil {
		return nil, ErrProductNotFound
	}

	if product.Stock < req.Quantity {
		return nil, ErrInsufficientStock
	}

	var cart models.Cart
	err = s.db.Where("user_id = ? ", userID).First(&cart).Error
	if err != nil {
		err = s.db.Create(&cart).Error
		if err != nil {
			return nil, err
		}
	}

	var item models.CartItem
	err = s.db.Where("cart_id = ? AND product_id = ? ", cart.ID, req.ProductID).First(&item).Error
	if err == nil {
		item.Quantity = item.Quantity + req.Quantity
		if item.Quantity > product.Stock {
			return nil, ErrInsufficientStock
		}
		s.db.Save(&item)
	} else {
		item = models.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
	}
	s.db.Create(&item)

	return s.GetUserCart(userID)
}

func (s *CartServices) UpdateCartItem(userID uint, itemID uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	var cartItem models.CartItem
	err := s.db.Joins("JOIN carts ON cart_items.cart_id = carts.id ").
		Where("cart_items.id = ? AND carts.user_id = ? ", itemID, userID).
		First(&cartItem).Error
	if err != nil {
		return nil, ErrCartItemNotFound
	}

	var product models.Product
	err = s.db.First(&product, cartItem.ProductID).Error
	if err != nil {
		return nil, ErrProductNotFound
	}

	if product.Stock < req.Quantity {
		return nil, ErrInsufficientStock
	}

	cartItem.Quantity = req.Quantity
	err = s.db.Save(&cartItem).Error
	if err != nil {
		return nil, err
	}

	return s.GetUserCart(userID)
}

func (s *CartServices) RemoveCartItem(userID uint, itemID uint) error {
	var cartItem models.CartItem
	err := s.db.Where("id = ? AND cart_id IN (?)", itemID,
		s.db.Select("id").Table("carts").
			Where("user_id = ?", userID)).
		Delete(&cartItem).Error
	if err != nil {
		return nil
	}

	return nil
}

func (s *CartServices) ClearCart(cartID uint) error {
	var cart models.Cart
	err := s.db.Where("user_id = ?", cart.UserID).First(&cart).Error
	if err != nil {
		return ErrCartNotFound
	}

	err = s.db.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error
	if err != nil {
		return err
	}

	return nil
}
