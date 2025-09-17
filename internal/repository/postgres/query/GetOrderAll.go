package query

const GetOrderAll = `SELECT
    o.order_uid,
    o.track_number,
    o.entry,
    o.locale,
    o.internal_signature,
    o.customer_id,
    o.delivery_service,
    o.shardkey,
    o.sm_id,
    o.date_created,
    o.off_shard,

    -- Delivery
    d.name AS delivery_name,
    d.phone AS delivery_phone,
    d.zip AS delivery_zip,
    d.city AS delivery_city,
    d.address AS delivery_address,
    d.region AS delivery_region,
    d.email AS delivery_email,

    -- Payment
    p.transaction,
    p.request_id,
    p.currency,
    p.provider,
    p.amount,
    p.payment_dt,
    p.bank,
    p.delivery_cost,
    p.goods_total,
    p.custom_fee,

    -- Items (агрегируем без размножения строк заказа)
    COALESCE(
        (
            SELECT json_agg(
                json_build_object(
                    'chrt_id', i.chrt_id,
                    'track_number', i.track_number,
                    'price', i.price,
                    'rid', i.rid,
                    'name', i.name,
                    'sale', i.sale,
                    'size', i.size,
                    'total_price', i.total_price,
                    'nm_id', i.nm_id,
                    'brand', i.brand,
                    'status', i.status
                )
            )
            FROM order_items oi
            JOIN items i ON oi.item_id = i.id
            WHERE oi.order_uid = o.order_uid
        ),
        '[]'::json
    ) AS items

FROM orders o
JOIN delivery d ON o.delivery_id = d.id
JOIN payments p ON o.payment_id = p.id

ORDER BY o.date_created DESC
LIMIT $1;`