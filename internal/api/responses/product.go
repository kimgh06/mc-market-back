package responses

import (
	"maple/internal/schema"
	"strconv"
	"time"
)

type Product struct {
	ID      string `json:"id"`
	Creator string `json:"creator"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Category string `json:"category"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ProductFromSchema(product *schema.Product) Product {
	var response Product
	response.ID = strconv.FormatUint(uint64(product.ID), 10)
	response.Creator = strconv.FormatUint(uint64(product.Creator), 10)
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
	response.ID = strconv.FormatUint(uint64(product.Product.ID), 10)
	response.Creator = ShortUser{
		ID: strconv.FormatUint(uint64(product.User.ID), 10),
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
