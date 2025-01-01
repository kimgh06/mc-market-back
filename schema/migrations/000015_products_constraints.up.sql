alter table products
    rename constraint title_check to name_length;
alter table products
    rename constraint price_not_equal_price_discount to price_discount_differ;
alter table products
    add constraint description_length check ( char_length(description) <= 300 ),
    add constraint usage_length check ( char_length(usage) <= 300 ),
    add constraint price check ( price >= 0 ),
    add constraint price_discount check ( price_discount >= 0 );