package services

import (
	"errors"
	"fmt"

	"github.com/aboloredev/armory/internal/dto"
	"github.com/aboloredev/armory/internal/models"
	"github.com/aboloredev/armory/internal/utils"
	"gorm.io/gorm"
)

type ProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{
		db: db,
	}
}

func (s *ProductService) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	var existingCat models.Category

	err := s.db.Where("name = ? ", req.Name).First(&existingCat).Error
	if err == nil {
		return nil, fmt.Errorf("category with this name already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	err = s.db.Create(&category).Error
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		Name:        category.Name,
		ID:          category.ID,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

func (s *ProductService) GetCategory(categoryID uint) (*dto.CategoryResponse, error) {
	var category models.Category

	err := s.db.Where("id = ? AND is_active = ? ", categoryID, true).First(&category).Error
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		Name:        category.Name,
		ID:          category.ID,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

func (s *ProductService) GetAllCategories(page int, pageSize int) ([]*dto.CategoryResponse, *utils.PaginatedMeta, error) {
	var categories []models.Category
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	s.db.Model(&models.Category{}).Count(&total)

	err := s.db.Offset(offset).Limit(pageSize).Find(&categories).Error
	if err != nil {
		return nil, nil, err
	}

	response := make([]*dto.CategoryResponse, 0, len(categories))

	for _, category := range categories {
		response = append(response, &dto.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			IsActive:    category.IsActive,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		})
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

func (s *ProductService) UpdateCategory(categoryID uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	var category models.Category

	err := s.db.Where("id = ? AND is_active = ? ", categoryID, true).First(&category).Error
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Description = req.Description
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	err = s.db.Save(&category).Error
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		Name:        category.Name,
		ID:          category.ID,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

func (s *ProductService) DeleteCategory(categoryID uint) error {
	var category models.Category

	var count int64
	s.db.Model(&models.Product{}).Where("category_id = ?", categoryID).Count(&count)
	if count > 0 {
		return fmt.Errorf("cannot delete category with %d products", count)
	}

	err := s.db.Delete(&category, categoryID).Error
	if err != nil {
		return nil
	}

	return nil
}

// A Helper function that converts models to dtoresponse
func (s *ProductService) ConvertToResponse(product *models.Product) *dto.ProductResponse {
	images := make([]dto.ProductImageResponse, len(product.Images))

	for i := range product.Images {
		images[i] = dto.ProductImageResponse{
			ID:        product.Images[i].ID,
			URL:       product.Images[i].URL,
			AltText:   product.Images[i].AltText,
			IsPrimary: product.Images[i].IsPrimary,
			CreatedAt: product.Images[i].CreatedAt,
		}
	}

	return &dto.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		SKU:         product.SKU,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		IsActive:    product.IsActive,
		Category: dto.CategoryResponse{
			ID:          product.Category.ID,
			Name:        product.Category.Name,
			Description: product.Category.Description,
		},
		Images:    images,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}

func (s *ProductService) CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	var category models.Category

	err := s.db.First(&category, req.CategoryID).Error
	if err != nil {
		return nil, err
	}

	var existingProduct models.Product

	err = s.db.Where("name = ? AND sku = ?", req.Name, req.SKU).
		First(&existingProduct).Error

	if err == nil {
		return nil, fmt.Errorf("product with this name and SKU already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CategoryID:  category.ID,
		Stock:       req.Stock,
		SKU:         req.SKU,
	}

	err = s.db.Create(&product).Error
	if err != nil {
		return nil, err
	}

	return s.ConvertToResponse(&product), nil
}

func (s *ProductService) UpdateProduct(productID uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	var product models.Product

	err := s.db.Where("id = ? AND is_active = ? ", productID, true).First(&product).Error
	if err != nil {
		return nil, err
	}

	product.CategoryID = req.CategoryID
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	err = s.db.Save(&product).Error
	if err != nil {
		return nil, err
	}

	return s.ConvertToResponse(&product), nil
}

func (s *ProductService) GetProduct(productID uint) (*dto.ProductResponse, error) {
	var product models.Product

	err := s.db.Preload("Category").Preload("Images").
		Where("id = ? AND is_active = ? ", productID, true).First(&product).Error
	if err != nil {
		return nil, err
	}

	return s.ConvertToResponse(&product), nil
}

func (s *ProductService) DeleteProduct(productID uint) error {
	var product models.Product

	err := s.db.Delete(&product, productID).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *ProductService) AddProductImage(productID uint, url string, altText string) error {
	var count int64

	s.db.Model(&models.ProductImage{}).Where("product_id = ?", productID).Count(&count)

	image := &models.ProductImage{
		ProductID: productID,
		URL:       url,
		AltText:   altText,
		IsPrimary: count == 0,
	}

	err := s.db.Create(&image).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *ProductService) RemoveProductImage(productID uint) error {
	var image models.ProductImage

	err := s.db.Delete(&image, productID).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *ProductService) AllProducts(page int, pageSize int) ([]*dto.ProductResponse, *utils.PaginatedMeta, error) {
	var products []models.Product
	var total int64

	offset := utils.PaginationHelper(page, pageSize)

	s.db.Model(&models.Product{}).Count(&total)

	err := s.db.Preload("Category").Preload("Images").
		Offset(offset).Limit(pageSize).Find(&products).Error
	if err != nil {
		return nil, nil, err
	}

	response := make([]*dto.ProductResponse, 0, len(products))

	for _, product := range products {
		response = append(response, s.ConvertToResponse(&product))
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
