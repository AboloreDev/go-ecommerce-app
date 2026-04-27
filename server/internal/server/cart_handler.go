package server

import (
	"fmt"
	"strconv"

	"github.com/aboloredev/armory/internal/dto"
	"github.com/aboloredev/armory/internal/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Get user's cart
// @Description Retrieve current user's shopping cart with all items
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.CartResponse} "Cart retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Cart not found"
// @Router /cart [get]

func (s *Server) GetUserCart(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	response, err := s.cartService.GetUserCart(userID)
	if err != nil {
		utils.BadRequest(ctx, "Failed to get cart", err)
		return
	}

	utils.SuccessResponse(ctx, "Cart fetched successfully", response)
}

// @Summary Add item to cart
// @Description Add a product to the user's shopping cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AddToCartRequest true "Item to add to cart"
// @Success 200 {object} utils.Response{data=dto.CartResponse} "Item added to cart successfully"
// @Failure 400 {object} utils.Response "Invalid request data or insufficient stock"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /cart [post]

func (s *Server) AddToCartHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	fmt.Println("userID from JWT:", userID)

	var req dto.AddToCartRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.cartService.AddToCartSevices(userID, &req)
	if err != nil {
		utils.BadRequest(ctx, "Add to cart failed", err)
		return
	}

	utils.SuccessResponse(ctx, "Item Added to Cart successfully", response)
}

// @Summary Update cart item quantity
// @Description Update the quantity of an item in the user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cart Item ID"
// @Param request body dto.UpdateCartItemRequest true "New quantity"
// @Success 200 {object} utils.Response{data=dto.CartResponse} "Cart item updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data or insufficient stock"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /cart/{id} [patch]
func (s *Server) UpdateCartItemHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Item ID", err)
	}
	itemID := uint(id)

	var req dto.UpdateCartItemRequest
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.cartService.UpdateCartItem(userID, itemID, &req)
	if err != nil {
		utils.BadRequest(ctx, "Cart update failed", err)
		return
	}

	utils.SuccessResponse(ctx, "Cart Updated successfully", response)
}

// @Summary Remove item from cart
// @Description Remove an item from the user's shopping cart
// @Tags Cart
// @Security BearerAuth
// @Param id path int true "Cart Item ID"
// @Success 200 {object} utils.Response "Item removed from cart successfully"
// @Failure 400 {object} utils.Response "Invalid cart item ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /cart/{id} [delete]
func (s *Server) RemoveCartItemHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Item ID", err)
	}
	itemID := uint(id)

	err = s.cartService.RemoveCartItem(userID, itemID)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	utils.SuccessResponse(ctx, "Item Removed from Cart successfully", nil)
}

// @Summary Clear user cart
// @Description Remove all item from the user's shopping cart
// @Tags Cart
// @Security BearerAuth
// @Param id path int true "Cart Item ID"
// @Success 200 {object} utils.Response "Item removed from cart successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /cart [delete]
func (s *Server) ClearUserCartHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	err := s.cartService.ClearCart(userID)
	if err != nil {
		utils.BadRequest(ctx, "Invalid user id", err)
		return
	}

	utils.SuccessResponse(ctx, "Cart cleared successfully", nil)
}
