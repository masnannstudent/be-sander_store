package entities

import "time"

type ProductModels struct {
	ID           uint64                 `gorm:"column:id;primaryKey" json:"id"`
	Name         string                 `gorm:"column:name" json:"name"`
	Price        uint64                 `gorm:"column:price" json:"price"`
	Description  string                 `gorm:"column:description" json:"description"`
	Discount     uint64                 `gorm:"column:discount" json:"discount"`
	Rating       float64                `gorm:"column:rating" json:"rating"`
	TotalReviews uint64                 `gorm:"column:total_reviews" json:"total_reviews"`
	Status       string                 `gorm:"column:status;type:VARCHAR(255)" json:"status"`
	CreatedAt    time.Time              `gorm:"column:created_at;type:timestamp" json:"created_at"`
	UpdatedAt    time.Time              `gorm:"column:updated_at;type:timestamp" json:"updated_at"`
	DeletedAt    *time.Time             `gorm:"column:deleted_at;type:TIMESTAMP NULL;index" json:"deleted_at"`
	Photos       []ProductPhotoModels   `gorm:"foreignKey:ProductID" json:"photos"`
	Categories   []*CategoryModels      `gorm:"many2many:product_categories;" json:"categories"`
	Variants     []ProductVariantModels `gorm:"foreignKey:ProductID" json:"variants"`
}

type ProductVariantModels struct {
	ID        uint64     `gorm:"column:id;primaryKey" json:"id"`
	ProductID uint64     `gorm:"column:product_id" json:"product_id"`
	Size      string     `gorm:"column:size;type:VARCHAR(255)" json:"size"`
	Color     string     `gorm:"column:color;type:VARCHAR(255)" json:"color"`
	Stock     uint64     `gorm:"column:stock" json:"stock"`
	Weight    uint64     `gorm:"column:weight" json:"weight"`
	CreatedAt time.Time  `gorm:"column:created_at;type:timestamp" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:timestamp" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:TIMESTAMP NULL;index" json:"deleted_at"`
}

type ProductPhotoModels struct {
	ID        uint64 `gorm:"column:id;primaryKey" json:"id"`
	ProductID uint64 `gorm:"column:product_id" json:"product_id"`
	URL       string `gorm:"column:url" json:"url"`
}

func (ProductModels) TableName() string {
	return "product"
}

func (ProductPhotoModels) TableName() string {
	return "product_photo"
}

func (ProductVariantModels) TableName() string {
	return "variants"
}
