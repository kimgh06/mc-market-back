package responses

import (
	"maple/internal/nullable"
	"maple/internal/schema"
	"strconv"
	"time"
)

type Product struct {
	ID      string `json:"id"`
	Creator string `json:"creator"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Usage       string `json:"usage"`

	Details string `json:"details"`

	Tags []string `json:"tags"`

	Category string `json:"category"`

	Price         int    `json:"price"`
	PriceDiscount *int32 `json:"price_discount"`

	Purchases int `json:"purchases"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ProductFromSchema(row *schema.GetProductByIdRow) Product {
	product := row.Product

	var response Product
	response.ID = strconv.FormatUint(uint64(product.ID), 10)
	response.Creator = strconv.FormatUint(uint64(product.Creator), 10)
	response.Name = product.Name
	response.Description = product.Description
	response.Usage = product.Usage
	response.Details = product.Details
	response.Tags = product.Tags
	response.Category = product.Category
	response.CreatedAt = product.CreatedAt
	response.UpdatedAt = product.UpdatedAt
	response.Price = int(product.Price)
	response.PriceDiscount = nullable.Int32ToPointer(product.PriceDiscount)
	response.Purchases = int(row.Count)

	return response
}

type ProductWithShortUser struct {
	Product
	Creator ShortUser `json:"creator"`
}

func ProductWithShortUserFromSchema(product *schema.GetProductByIdRow) ProductWithShortUser {
	var response ProductWithShortUser
	response.Product = ProductFromSchema(product)
	response.Creator = ShortUser{
		ID: strconv.FormatUint(uint64(product.User.ID), 10),
	}
	if product.User.Nickname.Valid {
		response.Creator.Nickname = &product.User.Nickname.String
	}

	return response
}
