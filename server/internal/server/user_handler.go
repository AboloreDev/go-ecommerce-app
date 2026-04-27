package server

import (
	"strconv"

	"github.com/aboloredev/armory/internal/dto"
	"github.com/aboloredev/armory/internal/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Get user profile
// @Description Get current authenticated user's profile information
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.UserResponse} "Profile retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "User not found"
// @Router /users/profile [get]
func (s *Server) GetUserProfileHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	response, err := s.userService.GetUserProfileService(userID)
	if err != nil {
		utils.UnAuthorized(ctx, "Unauthorized Request", err)
		return
	}

	utils.SuccessResponse(ctx, "User Found Successfully", response)
}

// @Summary Update user profile
// @Description Update current authenticated user's profile information
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} utils.Response{data=dto.UserResponse} "Profile updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /users/profile [patch]
func (s *Server) UpdateUserProfileHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	var req dto.UpdateProfileRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request", err)
		return
	}

	response, err := s.userService.UpdateProfileService(userID, &req)
	if err != nil {
		utils.UnAuthorized(ctx, "Unauthorized Request", err)
		return
	}

	utils.SuccessResponse(ctx, "User Updated Successfully", response)
}

// @Summary Get All user
// @Description Get current authenticated user's profile information
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 "successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Admin access required"
// @Router /users/all [get]
func (s *Server) GetAllUsersHandler(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	response, total, err := s.userService.GetAllUsersService(page, pageSize)
	if err != nil {
		utils.BadRequest(ctx, "Failed to fetch User", err)
	}

	utils.PaginatedSuccessResponse(ctx, "Success", response, utils.PaginatedMeta{
		Page:  page,
		Limit: pageSize,
		Total: total,
	})
}
