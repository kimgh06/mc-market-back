package products

import "database/sql"

func GetProductPrice(discount sql.NullInt32, price int32) int32 {
	if discount.Valid {
		return discount.Int32
	}
	return price
}
