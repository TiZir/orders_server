SET search_path TO postgres;

create table public.orders (
	order_uid          varchar(256)  PRIMARY KEY, 
	track_number       varchar(256)  UNIQUE, 
	entry              varchar(256), 
	locale             varchar(256), 
	internal_signature varchar(256), 
	customer_id        varchar(256),  
	delivery_service   varchar(256), 
	shardkey           varchar(256),  
	sm_id              int,
	date_created       int,
    oof_shard          varchar(256)
);

create table public.delivery (
    order_uid varchar(256) PRIMARY KEY,
	name    varchar(256),
	phone   varchar(256),
	zip     varchar(256),
	city    varchar(256),
	address varchar(256),
	region  varchar(256),
	email   varchar(256),
    FOREIGN KEY (order_uid) REFERENCES  public.orders (order_uid)
);

create table public.payment (
	transaction   varchar(256) PRIMARY KEY,
    request_id    int,
	currency      varchar(256), 
	provider      varchar(256),
	amount        int,
	payment_dt    int,  
	bank          varchar(256),
	delivery_cost int,
	goods_total   int,
    custom_fee    int,
    FOREIGN KEY (transaction) REFERENCES  public.orders (order_uid)
);

create table public.items (
	chrt_id      int PRIMARY KEY, 
    track_number varchar(256),
	price        int,    
	rid          varchar(256), 
	name         varchar(256), 
	sale         int,    
	size         varchar(256), 
	total_price  int,    
	nm_id        int,    
	brand        varchar(256), 
    status       int,
    FOREIGN KEY (track_number) REFERENCES  public.orders (track_number)
);
