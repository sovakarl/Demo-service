CREATE TABLE IF NOT EXISTS order_items(
    chrt_id BIGINT IS NOT NULL ,
    track_number VARCHAR(100) IS NOT NULL,
    price INT IS not NULL,
    rid VARCHAR(100) IS NOT NULL,
    name VARCHAR(100) IS NOT NULL,
    sale INT IS NOT NULL,
    size VARCHAR(100),
    total_price INT IS NOT NULL,
    nm_id INT IS NOT NULL,
    brand VARCHAR(100),
    status INT IS NOT NULL
);