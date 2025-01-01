package products

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
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
	Creator         *int    `form:"creator"`
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
			&i.Product.Details,
			&i.Product.Price,
			&i.Product.PriceDiscount,
			&i.Product.CreatedAt,
			&i.Product.UpdatedAt,
			&i.User.ID,
			&i.User.Nickname,
			&i.Count,
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

	offset := utilities.Clamp(api.QueryIntDefault(ctx, "offset", math.MaxInt), 0, math.MaxInt)
	limit := utilities.Clamp(api.QueryIntDefault(ctx, "limit", 20), 0, 20)
	orderBy := strings.ToLower(api.QueryStringDefault(ctx, "order_by", "time"))
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
		Select("products.id, products.creator, products.category, products.name, products.description, products.usage, products.details, products.price, products.price_discount, products.created_at, products.updated_at, u.id, u.nickname, count(pu)").
		From("products").
		LeftJoin(fmt.Sprintf("public.users u ON u.id = products.creator")).
		LeftJoin(fmt.Sprintf("public.purchases pu on pu.product = products.id"))

	whereClause := squirrel.And{squirrel.Lt{"products.id": offset}}

	if filters.Creator != nil {
		whereClause = append(whereClause, squirrel.Eq{"products.creator": filters.Creator})
	}
	if filters.Category != nil {
		whereClause = append(whereClause, squirrel.Eq{"products.category": filters.Category})
	}
	if filters.Keyword != nil {
		keywordSanitized := SpecialCharacterRegexp.ReplaceAllString(*filters.Keyword, "")
		keywordSplit := strings.Split(keywordSanitized, " ")
		searchQuery := strings.Join(keywordSplit, " & ")
		whereClause = append(whereClause, squirrel.Expr(fmt.Sprintf("products.ts @@ to_tsquery(?)"), searchQuery))
	}
	if filters.PriceRangeStart != nil {
		whereClause = append(whereClause, squirrel.GtOrEq{"coalesce(products.price_discount, products.price)": *filters.PriceRangeStart})
	}
	if filters.PriceRangeEnd != nil {
		whereClause = append(whereClause, squirrel.Lt{"coalesce(products.price_discount, products.price)": *filters.PriceRangeEnd})
	}

	query = query.Where(whereClause).GroupBy("products.id, u.id")

	switch orderBy {
	case "purchases":
		query = query.OrderBy(fmt.Sprintf("count(pu) %s", sort))
	case "price":
		query = query.OrderBy(fmt.Sprintf("coalesce(products.price_discount, products.price) %s", sort))
	case "time":
		fallthrough
	default:
		query = query.OrderBy(fmt.Sprintf("products.created_at %s", sort))
	}

	query = query.Limit(uint64(limit))

	rows, err := query.RunWith(a.Conn).Query()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	products, err := scanProductRows(rows)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	creatorIds := utilities.Map(products, func(row *schema.ListProductsRow) uint64 {
		return uint64(row.User.ID)
	})

	creatorIdToUsernameMap, err := a.SurgeAPI.ResolveUsernamesAsMap(creatorIds)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.MakeJSON(err.Error()))
		return
	}

	fullProducts := utilities.Map(products, func(u *schema.ListProductsRow) responses.ProductWithShortUser {
		converted := responses.ProductWithShortUserFromSchema((*schema.GetProductByIdRow)(u))

		if username, found := creatorIdToUsernameMap[uint64(u.User.ID)]; found {
			converted.Creator.Username = &username
		}

		return converted
	})

	ctx.JSON(http.StatusOK, fullProducts)
}
