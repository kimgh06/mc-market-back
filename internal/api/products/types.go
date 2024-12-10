package products

type CreateProductBody struct {
	Name        string `json:"name" validate:"required,lte=50"`
	Description string `json:"description" validate:"required,lte=300"`
	Usage       string `json:"usage" validate:"lte=300"`
	Category    string `json:"category" validate:"required"`
	Price       int32  `json:"price" validate:"gte=0"`
	Creator     uint64 `json:"creator"`
}

type UpdateProductBody struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Usage       *string `json:"usage"`
	Category    *string `json:"category"`
	Price       *int32  `json:"price" validate:"gte=0"`
	Creator     *uint64 `json:"creator"`
}
