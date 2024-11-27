package responses

import (
	"maple/internal/schema"
	"time"
)

type Product struct {
	ID      uint64 `json:"id"`
	Creator uint64 `json:"creator"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Category string `json:"category"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ProductFromSchema(product *schema.Product) Product {
	var response Product
	response.ID = uint64(product.ID)
	response.Creator = uint64(product.Creator)
	response.Name = product.Name
	response.Description = product.Description
	response.Category = product.Category
	response.CreatedAt = product.CreatedAt
	response.UpdatedAt = product.UpdatedAt

	return response
}

type ProductWithShortUser struct {
	Product
	Creator ShortUser `json:"creator"`
}

func ProductWithShortUserFromSchema(product *schema.ListProductsRow) ProductWithShortUser {
	var response ProductWithShortUser
	response.ID = uint64(product.Product.ID)
	response.Creator = ShortUser{
		ID: uint64(product.User.ID),
	}
	if product.User.Nickname.Valid {
		response.Creator.Nickname = &product.User.Nickname.String
	}
	response.Name = product.Product.Name
	response.Description = product.Product.Description
	response.CreatedAt = product.Product.CreatedAt
	response.UpdatedAt = product.Product.UpdatedAt
	response.Category = product.Product.Category

	return response
}
