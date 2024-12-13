alter table products
    add constraint price_not_equal_price_discount check ( price != products.price_discount );