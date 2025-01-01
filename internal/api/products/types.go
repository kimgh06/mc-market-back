package products

type CreateProductBody struct {
	Name          string   `json:"name" validate:"required,lte=50"`
	Description   string   `json:"description" validate:"required,lte=300"`
	Details       string   `json:"details" validate:"required,lte=5000"`
	Usage         string   `json:"usage" validate:"lte=300"`
	Category      string   `json:"category" validate:"required"`
	Price         int32    `json:"price" validate:"gte=0"`
	PriceDiscount *int32   `json:"price_discount"`
	Creator       uint64   `json:"creator"`
	Tags          []string `json:"tags"`
}

type UpdateProductBody struct {
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	Details       string  `json:"details"`
	Usage         *string `json:"usage"`
	Category      *string `json:"category"`
	Price         *int32  `json:"price"`
	PriceDiscount *int32  `json:"price_discount"`
	Creator       *uint64 `json:"creator"`
}
