alter table products
    rename constraint name_length to title_check;
alter table products
    rename constraint price_discount_differ to price_not_equal_price_discount;
alter table products
    drop constraint description_length,
    drop constraint usage_length,
    drop constraint price,
    drop constraint price_discount;