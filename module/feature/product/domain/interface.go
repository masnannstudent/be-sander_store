package domain

import (
	"github.com/gofiber/fiber/v2"
	"ruti-store/module/entities"
)

type ProductRepositoryInterface interface {
	GetPaginatedProducts(page, pageSize int) ([]*entities.ProductModels, error)
	GetTotalItems() (int64, error)
	GetProductByID(productID uint64) (*entities.ProductModels, error)
	CreateProduct(product *entities.ProductModels, categoryIDs []uint64) (*entities.ProductModels, error)
	UpdateProduct(productID uint64, newData *entities.ProductModels, categoryIDs []uint64) error
	DeleteProduct(productID uint64) error
	UpdateTotalReview(productID uint64) error
	UpdateProductRating(productID uint64, newRating float64) error
	GetProductReviews(page, perPage int) ([]*entities.ProductModels, error)
	AddPhotoProduct(newData *entities.ProductPhotoModels) (*entities.ProductPhotoModels, error)
	UpdateProductPhoto(productID uint64, newPhotoURL string) error
	ReduceStockWhenPurchasing(productID, quantity uint64) error
	IncreaseStock(productID, quantity uint64) error
	GenerateRecommendationProduct() ([]string, error)
	FindAllProductRecommendation(productsFromAI []string) ([]*entities.ProductModels, error)
	SearchAndPaginateProducts(name string, page, pageSize int) ([]*entities.ProductModels, int64, error)
	CreateVariantProduct(newData *entities.ProductVariantModels) (*entities.ProductVariantModels, error)
	UpdateProductStatus(productID uint64, status string) error
}

type ProductServiceInterface interface {
	GetAllProducts(page, pageSize int) ([]*entities.ProductModels, int64, error)
	GetProductsPage(currentPage, pageSize, totalItems int) (int, int, int, error)
	GetProductByID(productID uint64) (*entities.ProductModels, error)
	CreateProduct(req *CreateProductRequest) (*entities.ProductModels, error)
	UpdateProduct(productID uint64, req *UpdateProductRequest) error
	DeleteProduct(productID uint64) error
	UpdateTotalReview(productID uint64) error
	UpdateProductRating(productID uint64, newRating float64) error
	GetProductReviews(page, perPage int) ([]*entities.ProductModels, int64, error)
	AddPhotoProducts(req *AddPhotoProductRequest) (*entities.ProductPhotoModels, error)
	UpdatePhotoProduct(productID uint64, photo string) error
	ReduceStockWhenPurchasing(productID, quantity uint64) error
	IncreaseStock(productID, quantity uint64) error
	GetProductRecommendation() ([]string, error)
	GetAllProductsRecommendation() ([]*entities.ProductModels, error)
	SearchAndPaginateProducts(name string, page, pageSize int) ([]*entities.ProductModels, int64, error)
	CreateVariantProduct(req *CreateVariantRequest) (*entities.ProductVariantModels, error)
	UpdateStatusProduct(req *UpdateStatusRequest) error
}

type ProductHandlerInterface interface {
	GetAllProducts(c *fiber.Ctx) error
	GetProductByID(c *fiber.Ctx) error
	CreateProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
	GetAllProductsReview(c *fiber.Ctx) error
	AddPhotoProduct(c *fiber.Ctx) error
	UpdatePhotoProduct(c *fiber.Ctx) error
	GetProductRecommendation(c *fiber.Ctx) error
	GetAllProductsRecommendation(c *fiber.Ctx) error
	CreateVariantProduct(c *fiber.Ctx) error
	UpdateStatusProduct(c *fiber.Ctx) error
}
