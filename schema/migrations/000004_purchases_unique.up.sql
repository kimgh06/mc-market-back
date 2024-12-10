alter table purchases
    add constraint purchases_unique unique (purchaser, product);