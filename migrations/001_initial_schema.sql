CREATE table users(
  id                serial primary key,
  name              varchar(100) not null,
  email             varchar(255) not null unique,
  password_hash     text not null,
  role              varchar(20) not null default 'customer',
  created_at        timestamp not null default now()
);


CREATE table categories(
  id                serial primary key,
  name              varchar(100) not null,
  slug              varchar(100) not null unique
);

CREATE table products(
  id                serial primary key,
  name              varchar(255) not null,
  stock             int not null default 0,
  category_id       int references categories(id) on delete set null,
  slug              varchar(255) not null unique,
  description       text,
  price             numeric(10,2) not null,
  image_url         text,
  created_at        timestamp not null default now()
);

CREATE table orders(
  id                serial primary key,
  user_id           int not null references users(id) on delete cascade,
  status            varchar(20) not null default 'pending',
  total_price       numeric(10,2) not null default 0,
  created_at        timestamp not null default now()
);

CREATE table order_items(
  id                serial primary key,
  order_id          int not null references orders(id) on delete cascade,
  product_id        int not null references products(id) on delete restrict,
  quantity          int not null check(quantity > 0),
  unit_price        numeric(10,2) not null
);

CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_slug        ON products(slug);
CREATE INDEX idx_orders_user_id       ON orders(user_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);