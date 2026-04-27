package server

import (
	"strconv"

	"github.com/aboloredev/armory/internal/dto"
	"github.com/aboloredev/armory/internal/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Create a new category
// @Description Create a new product category (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCategoryRequest true "Category data"
// @Success 201 {object} utils.Response{data=dto.CategoryResponse} "Category created successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Admin access required"
// @Router /categories [post]
func (s *Server) CreateCategoryHandler(ctx *gin.Context) {
	var req dto.CreateCategoryRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	response, err := s.productService.CreateCategory(&req)
	if err != nil {
		utils.BadRequest(ctx, "Category Creation Failed", err)
		return
	}

	utils.CreatedResponse(ctx, "Category Created Successfullty", response)
}

// @Summary Update a category
// @Description Update an existing category (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param request body dto.UpdateCategoryRequest true "Category update data"
// @Success 200 {object} utils.Response{data=dto.CategoryResponse} "Category updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Admin access required"
// @Router /categories/{id} [patch]
func (s *Server) UpdateCategoryHandler(ctx *gin.Context) {
	idString := ctx.Param("id")
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Product ID", err)
	}
	categoryID := uint(id)

	var req dto.UpdateCategoryRequest
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	response, err := s.productService.UpdateCategory(categoryID, &req)

	if err != nil {
		utils.BadRequest(ctx, "Category Update Failed", err)
		return
	}

	utils.SuccessResponse(ctx, "Category Updated Successfullty", response)

}

// @Summary Get a category
// @Description Delete a category (Admin only)
// @Tags Categories
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} utils.Response "Category deleted successfully"
// @Failure 400 {object} utils.Response "Invalid category ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Admin access required"
// @Router /categories/{id} [get]
func (s *Server) GetCategoryHandler(ctx *gin.Context) {
	idString := ctx.Param("id")
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Product ID", err)
		return
	}
	categoryID := uint(id)

	response, err := s.productService.GetCategory(categoryID)

	if err != nil {
		utils.BadRequest(ctx, "Invalid Category ID", err)
		return
	}

	utils.SuccessResponse(ctx, "Success", response)

}

// @Summary Get all categories
// @Description Retrieve all active categories
// @Tags Categories
// @Produce json
// @Success 200 {object} utils.Response{data=[]dto.CategoryResponse} "Categories retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /categories [get]
func (s *Server) GetAllCategoriesHandler(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	response, meta, err := s.productService.GetAllCategories(page, pageSize)

	if err != nil {
		utils.BadRequest(ctx, "Invalid Category ID", err)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Success", response, *meta)
}

// @Summary Delete a category
// @Description Delete a category (Admin only)
// @Tags Categories
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} utils.Response "Category deleted successfully"
// @Failure 400 {object} utils.Response "Invalid category ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Admin access required"
// @Router /categories/{id} [delete]
func (s *Server) DeleteCategoryHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Product ID", err)
		return
	}
	categoryID := uint(id)

	err = s.productService.DeleteCategory(categoryID)

	if err != nil {
		utils.BadRequest(ctx, "Error deleting category", err)
		return
	}

	utils.SuccessResponse(ctx, "Success", nil)
}

// @Summary Get all products
// @Description Retrieve paginated list of active products
// @Tags Products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse{data=[]dto.ProductResponse} "Products retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /products [get]
func (s *Server) AllProductHandler(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	response, meta, err := s.productService.AllProducts(page, pageSize)

	if err != nil {
		utils.BadRequest(ctx, "Product Fetch Failed", err)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Success", response, *meta)
}

// @Summary Create a new product
// @Description Create a new product (Admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateProductRequest true "Product data"
// @Success 201 {object} utils.Response{data=dto.ProductResponse} "Product created successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Admin access required"
// @Router /products [post]
func (s *Server) CreateProductHandler(ctx *gin.Context) {
	var req dto.CreateProductRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	response, err := s.productService.CreateProduct(&req)

	if err != nil {
		utils.BadRequest(ctx, "Product creation Failed", err)
		return
	}

	utils.CreatedResponse(ctx, "Product Created Successdully", response)
}

// @Summary Update a product
// @Description Update an existing product (Admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param request body dto.UpdateProductRequest true "Product update data"
// @Success 200 {object} utils.Response{data=dto.ProductResponse} "Product updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Admin access required"
// @Router /products/{id} [patch]
func (s *Server) UpdateProductHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Product ID", err)
	}
	productID := uint(id)

	var req dto.UpdateProductRequest
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	response, err := s.productService.UpdateProduct(productID, &req)

	if err != nil {
		utils.BadRequest(ctx, "Product Update Failed", err)
		return
	}

	utils.SuccessResponse(ctx, "Product Updated Successdully", response)
}

// @Summary Get a product
// @Description Get a product
// @Tags Products
// @Param id path int true "Product ID"
// @Success 200 {object} utils.Response "Product fetched successfully"
// @Failure 400 {object} utils.Response "Invalid product ID"
// @Router /products/{id} [get]
func (s *Server) GetProductHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Product ID", err)
		return
	}
	productID := uint(id)

	response, err := s.productService.GetProduct(productID)

	if err != nil {
		utils.BadRequest(ctx, "Product Fetch Failed", err)
		return
	}

	utils.SuccessResponse(ctx, "Product Fetched Successdully", response)
}

// @Summary Delete a product
// @Description Delete a product (Admin only)
// @Tags Products
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} utils.Response "Product deleted successfully"
// @Failure 400 {object} utils.Response "Invalid product ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Admin access required"
// @Router /products/{id} [delete]
func (s *Server) DeleteProduct(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Product ID", err)
		return
	}
	productID := uint(id)

	err = s.productService.DeleteProduct(productID)

	if err != nil {
		utils.BadRequest(ctx, "Product Deletion Failed", err)
	}

	utils.SuccessResponse(ctx, "Product Deleted Successdully", nil)
}

// @Summary Upload product image
// @Description Upload an image for a product (Admin only)
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param image formData file true "Image file"
// @Success 200 {object} utils.Response{data=map[string]string} "Image uploaded successfully"
// @Failure 400 {object} utils.Response "Invalid request or file"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Admin access required"
// @Router /products/{id}/images [post]
func (s *Server) UploadProductImageService(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Product ID", err)
		return
	}
	productID := uint(id)

	file, err := ctx.FormFile("image")
	if err != nil {
		utils.BadRequest(ctx, "No file uploaded", err)
		return
	}

	url, err := s.uploadService.UploadProductImage(productID, file)
	if err != nil {
		utils.BadRequest(ctx, "Failed to upload file", err)
		return
	}

	err = s.productService.AddProductImage(productID, url, file.Filename)
	if err != nil {
		utils.BadRequest(ctx, "Error Adding", err)
		return
	}

	utils.SuccessResponse(ctx, "Upload successful", map[string]string{"url": url})
}

// @Summary Delete product image
// @Description Upload an image for a product (Admin only)
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} utils.Response{data=map[string]string} "Image Dleted successfully"
// @Failure 400 {object} utils.Response "Invalid id"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Admin access required"
// @Router /products/{id}/images [delete]
func (s *Server) DeleteProductImageService(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Product ID", err)
		return
	}
	productID := uint(id)

	err = s.uploadService.DeleteProductImage(productID)
	if err != nil {
		utils.BadRequest(ctx, "Failed to delete file", err)
	}

	s.productService.RemoveProductImage(productID)

	utils.SuccessResponse(ctx, "Delete successful", nil)
}
