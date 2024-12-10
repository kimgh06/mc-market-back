alter table purchases
    alter column claimed set not null,
    alter column cost set not null;