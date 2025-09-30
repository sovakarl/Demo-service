
--delivery--
CREATE TABLE delivery(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    phone VARCHAR(100),
    zip VARCHAR(100),
    city VARCHAR(100),
    address VARCHAR(100),
    region VARCHAR(100),
    email VARCHAR(100)
);


--payment--
CREATE TABLE payments(
    id SERIAL PRIMARY KEY,
    transaction VARCHAR(100) NOT NULL,
    request_id VARCHAR(50) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    provider VARCHAR(25) NOT NULL,
    amount INT NOT NULL,
    payment_dt INT NOT NULL,
    bank VARCHAR(50) NOT NULL,
    delivery_cost INT NOT NULL,
    goods_total INT,
    custom_fee INT 
);


--ITEMS--
CREATE TABLE items(
    id SERIAL PRIMARY KEY,
    chrt_id BIGINT  NOT NULL ,
    track_number VARCHAR(100) NOT NULL,
    price INT NOT NULL,
    rid VARCHAR(100)  NOT NULL,
    name VARCHAR(100)  NOT NULL,
    sale INT NOT NULL,
    size VARCHAR(100),
    total_price INT NOT NULL,
    nm_id INT NOT NULL,
    brand VARCHAR(100),
    status INT NOT NULL
);


--orders--
CREATE TABLE orders(
    order_uid VARCHAR(100) PRIMARY KEY,
    track_number VARCHAR(100),
    entry VARCHAR(50),

    delivery_id INT,
    payment_id INT,
    
    locale VARCHAR(10),
    internal_signature TEXT NOT NULL,
    customer_id VARCHAR(50) NOT NULL,
    delivery_service VARCHAR(50),
    shardkey TEXT,
    sm_id INT,
    date_created TIMESTAMP  DEFAULT CURRENT_TIMESTAMP,
    off_shard TEXT,

    FOREIGN KEY(delivery_id) REFERENCES  delivery(id),
    FOREIGN KEY(payment_id) REFERENCES  payments(id)
);

--ORDER_ITEMS--
CREATE TABLE order_items(
    order_uid VARCHAR(100),
    item_id INT,

    PRIMARY KEY(order_uid,item_id),
    FOREIGN KEY(item_id) REFERENCES items(id),
    FOREIGN KEY(order_uid) REFERENCES orders(order_uid) ON DELETE CASCADE
);