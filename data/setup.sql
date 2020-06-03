create table users (
  id            serial primary key,
  username      varchar(255),
  email         varchar(255) not null unique,
  password      varchar(255) not null,
  createdat     timestamp not null,
  deletedat     timestamp,
  firstname     varchar(64), 
  lastname      varchar(64),
  dateofbirth   date,
  country       varchar(64),
  city          varchar(64),
  gender        varchar(10)
);

create table sessions (
  id            serial primary key,
  uuid          varchar(64) not null unique,
  email         varchar(255),
  userid        integer references users(id),
  createdat     timestamp not null   
);

create table targets (
  id            serial primary key,
  url           varchar(500) not null unique,
  host          varchar(64) not null,             
  createdat     timestamp not null,
  name          varchar(64)
);


create table userstargets (
  id            serial primary key,
  uuid          varchar(64) not null unique;
  userid	    integer references users(id),
  targetid      integer references targets(id),
  createdat     timestamp not null
  deletedat     timestamp   
);

create table scrapers (
  id            serial primary key,
  name          varchar(64) not null,
  version       integer not null,
  targetid      integer references targets(id),
  createdat     timestamp not null  
);

create table scrapings (
  id            serial primary key,
  scraperid     integer references scrapers(id),
  createdat     timestamp
);

create table results (
  id            serial primary key,
  scraperid     integer references scrapers(id),
  scrapingid    integer references scraping(id),
  title         varchar(1000) not null,
  url           varchar(1000) not null unique,
  createdat     timestamp,
  updatedat     timestamp,
  data          json
);