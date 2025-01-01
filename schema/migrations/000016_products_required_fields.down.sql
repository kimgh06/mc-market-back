alter table products
    alter column description drop not null,
    alter column usage drop not null;