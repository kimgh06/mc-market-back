package products

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"maple/internal/api"
	"maple/internal/api/responses"
	"maple/internal/perrors"
	"maple/internal/schema"
	"maple/internal/utilities"
	"math"
	"net/http"
	"regexp"
	"strings"
)

var SpecialCharacterRegexp, _ = regexp.Compile("/[\\{\\}\\[\\]\\/?.,;:|\\)*~`!^\\-_+<>@\\#$%&\\\\\\=\\(\\'\\\"]/g")

type listFilters struct {
	Category        *string `form:"category"`
	Keyword         *string `form:"keyword"`
	PriceRangeStart *int    `form:"price_range_start"`
	PriceRangeEnd   *int    `form:"price_range_end"`
}

func scanProductRows(rows *sql.Rows) ([]*schema.ListProductsRow, error) {
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var items []*schema.ListProductsRow
	for rows.Next() {
		var i schema.ListProductsRow
		if err := rows.Scan(
			&i.Product.ID,
			&i.Product.Creator,
			&i.Product.Category,
			&i.Product.Name,
			&i.Product.Description,
			&i.Product.Usage,
			&i.Product.Price,
			&i.Product.PriceDiscount,
			&i.Product.CreatedAt,
			&i.Product.UpdatedAt,
			&i.User.ID,
			&i.User.Nickname,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func listProducts(ctx *gin.Context) {
	a := api.Get(ctx)
	//user := middlewares.GetUser(ctx)

	offset := utilities.Clamp(api.QueryIntDefault(ctx, "offset", 0), 0, math.MaxInt)
	limit := utilities.Clamp(api.QueryIntDefault(ctx, "limit", 20), 0, 20)
	sort := strings.ToLower(api.QueryStringDefault(ctx, "sort", "desc"))
	filters := listFilters{}

	// Ignore filter parse error
	if ctx.ShouldBindQuery(&filters) != nil {
		filters = listFilters{}
	}

	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	// Build Query
	query := squirrel.
		StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("products.id, products.creator, products.category, products.name, products.description, products.usage, products.price, products.price_discount, products.created_at, products.updated_at, u.id, u.nickname").
		From("products").
		LeftJoin(fmt.Sprintf("public.users u ON u.id = products.creator"))

	whereClause := squirrel.And{squirrel.Gt{"products.id": offset}}

	if filters.Category != nil {
		whereClause = append(whereClause, squirrel.Eq{"products.category": filters.Category})
	}
	if filters.Keyword != nil {
		keywordSanitized := SpecialCharacterRegexp.ReplaceAllString(*filters.Keyword, "")
		keywordSplit := strings.Split(keywordSanitized, " ")
		searchQuery := strings.Join(keywordSplit, " & ")
		whereClause = append(whereClause, squirrel.Expr(fmt.Sprintf("products.ts @@ to_tsquery('%s')", searchQuery)))
	}
	if filters.PriceRangeStart != nil {
		whereClause = append(whereClause, squirrel.GtOrEq{"coalesce(products.price_discount, products.price)": *filters.PriceRangeStart})
	}
	if filters.PriceRangeEnd != nil {
		whereClause = append(whereClause, squirrel.Lt{"coalesce(products.price_discount, products.price)": *filters.PriceRangeEnd})
	}

	query = query.Where(whereClause).
		OrderBy(fmt.Sprintf("products.id %s", sort)).
		Limit(uint64(limit))

	generated, args, _ := query.ToSql()
	logrus.Debugf("Generated Query: %s", generated)
	logrus.Debugf("Generated Args: %+v", args)

	rows, err := a.Conn.Query(generated, args...)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	products, err := scanProductRows(rows)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	fullProducts := utilities.Map(products, func(u *schema.ListProductsRow) responses.ProductWithShortUser {
		return responses.ProductWithShortUserFromSchema(u)
	})

	ctx.JSON(http.StatusOK, fullProducts)
}
