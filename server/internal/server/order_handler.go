package server

import (
	"strconv"

	"github.com/aboloredev/armory/internal/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Create an order
// @Description Create an order from the current user's cart
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Success 201 {object} utils.Response{data=dto.OrderResponse} "Order created successfully"
// @Failure 400 {object} utils.Response "Cart is empty or insufficient stock"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /orders [post]

func (s *Server) CreateOrderHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	response, err := s.orderService.CreateOrderService(userID)
	if err != nil {
		utils.BadRequest(ctx, "Order creation failed", err)
		return
	}
	utils.CreatedResponse(ctx, "Order created successfully", response)
}

// @Summary Get order by ID
// @Description Retrieve detailed information about a specific order
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} utils.Response{data=dto.OrderResponse} "Order retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid order ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Order not found"
// @Router /orders/{id} [get]

func (s *Server) GetOrderHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid orderID", err)
		return
	}
	orderID := uint(id)

	response, err := s.orderService.GetOrder(orderID, userID)
	if err != nil {
		utils.BadRequest(ctx, "Error fetching order", err)
		return
	}

	utils.SuccessResponse(ctx, "Order fetched successfully", response)
}

// @Summary Get user's orders
// @Description Retrieve paginated list of user's orders
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse{data=[]dto.OrderResponse} "Orders retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /orders [get]

func (s *Server) GetAllOrdersHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	response, meta, err := s.orderService.GetAllUserOrders(userID, page, pageSize)
	if err != nil {
		utils.BadRequest(ctx, "Failed to get user orders", err)
	}

	utils.PaginatedSuccessResponse(ctx, "Orders fetched succesfully", response, *meta)
}
