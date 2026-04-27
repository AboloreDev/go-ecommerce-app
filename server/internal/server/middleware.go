package server

import (
	"strings"

	"github.com/aboloredev/armory/internal/models"
	"github.com/aboloredev/armory/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Authorization Bearer token
		authHeader := ctx.GetHeader("authorization")
		if authHeader == "" {
			utils.UnAuthorized(ctx, "Authorization header required", nil)
			ctx.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.UnAuthorized(ctx, "Invalid Header Format", nil)
			ctx.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenParts[1], s.config.JWT.JWTSecret)
		if err != nil {
			utils.UnAuthorized(ctx, "Invalid Token", err)
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("user_email", claims.Email)
		ctx.Set("user_role", claims.Role)

		ctx.Next()
	}
}

func (s *Server) AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("user_role")

		if !exists {
			utils.Forbidden(ctx, "Forbidden", nil)
			ctx.Abort()
			return
		}

		if role != string(models.UserRoleAdmin) {
			utils.Forbidden(ctx, "Forbidden", nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
