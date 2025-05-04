create table adcard (
  id serial primary key,
  title text not null,
  image_url text not null,
  link_url text not null,
  created_at timestamp not null default now(),
  index_num integer not null default 0
);