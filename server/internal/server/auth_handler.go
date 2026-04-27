package server

import (
	"github.com/aboloredev/armory/internal/dto"
	"github.com/aboloredev/armory/internal/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Register a new user
// @Description Create a new user account with email,firstname, lastname and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration data"
// @Success 201 {object} utils.Response{data=dto.AuthResponse} "User registered successfully"
// @Failure 400 {object} utils.Response "Invalid request data or user already exists"
// @Router /auth/register [post]
func (s *Server) RegisterHandler(ctx *gin.Context) {
	var req dto.RegisterRequest
	err := ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	response, err := s.authService.UserRegistrationService(&req)
	if err != nil {
		utils.BadRequest(ctx, "Bad Request", err)
		return
	}

	utils.CreatedResponse(ctx, "User Created Successfully", response)
}

// @Summary User login
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User login credentials"
// @Success 200 {object} utils.Response{data=dto.AuthResponse} "Login successful"
// @Failure 401 {object} utils.Response "Invalid credentials"
// @Router /auth/login [post]
func (s *Server) LoginHandler(ctx *gin.Context) {
	var req dto.LoginRequest
	err := ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	response, err := s.authService.UserLoginService(&req)
	if err != nil {
		utils.UnAuthorized(ctx, "Login Failed", err)
		return
	}

	utils.SuccessResponse(ctx, "User Logged-in Successfully", response)
}

// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} utils.Response{data=dto.AuthResponse} "Token refreshed successfully"
// @Failure 401 {object} utils.Response "Invalid refresh token"
// @Router /auth/refresh [post]
func (s *Server) RefreshTokenHandler(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	err := ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	response, err := s.authService.RefreshTokenService(&req)
	if err != nil {
		utils.UnAuthorized(ctx, "Token refresh failed", err)
		return
	}

	utils.SuccessResponse(ctx, "Token Refreshed Successfully", response)
}

// @Summary User logout
// @Description Invalidate refresh token and logout user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token to invalidate"
// @Success 200 {object} utils.Response "Logout successful"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Router /auth/logout [post]
func (s *Server) LogoutHandler(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	err := ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	err = s.authService.LogoutService(&req)
	if err != nil {
		utils.InternalServerError(ctx, "Logout failed", err)
		return
	}

	utils.SuccessResponse(ctx, "Logout Successful", nil)
}
