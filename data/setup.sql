create table users (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  name       varchar(255),
  email      varchar(255) not null unique,
  password   varchar(255) not null,
  created_at timestamp not null   
);

create table sessions (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  email      varchar(255),
  user_id    integer references users(id),
  created_at timestamp not null   
);

create table targets (
  id         serial primary key,
  url        varchar(500) not null unique,
  host       varchar(64) not null,             
  created_at timestamp not null,
  name       varchar(64)
);


create table users_targets (
  id         serial primary key,
  uuid       varchar(64) not null unique;
  user_id	 integer references users(id),
  target_id  integer references targets(id),
  created_at timestamp not null
  deleted_at timestamp   
);

create table scrapers (
  id         serial primary key,
  name       varchar(64) not null,
  version    integer not null,
  target_id  integer references targets(id),
  created_at timestamp not null  
);

create table scraping (
  id         serial primary key,
  scraper_id integer references scrapers(id),
  created_at timestamp
);

create table results (
  id            serial primary key,
  scraper_id    integer references scrapers(id),
  scraping_id   integer references scraping(id),
  title         varchar(1000) not null,
  url           varchar(1000) not null unique,
  created_at    timestamp,
  updated_at    timestamp,
  data          json
);