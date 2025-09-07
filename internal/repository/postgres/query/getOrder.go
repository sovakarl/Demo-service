package query

const GetOrder = `SELECT
    o.order_uid,
    MIN(o.track_number) AS track_number,
    MIN(o.entry) AS entry,
    MIN(o.locale) AS locale,
    MIN(o.internal_signature) AS internal_signature,
    MIN(o.customer_id) AS customer_id,
    MIN(o.delivery_service) AS delivery_service,
    MIN(o.shardkey) AS shardkey,
    MIN(o.sm_id) AS sm_id,
    MIN(o.date_created) AS date_created,
    MIN(o.off_shard) AS off_shard,

    -- Delivery
    MIN(d.name) AS delivery_name,
    MIN(d.phone) AS delivery_phone,
    MIN(d.zip) AS delivery_zip,
    MIN(d.city) AS delivery_city,
    MIN(d.address) AS delivery_address,
    MIN(d.region) AS delivery_region,
    MIN(d.email) AS delivery_email,

    -- Payment
    MIN(p.transaction) AS transaction,
    MIN(p.request_id) AS request_id,
    MIN(p.currency) AS currency,
    MIN(p.provider) AS provider,
    MIN(p.amount) AS amount,
    MIN(p.payment_dt) AS payment_dt,
    MIN(p.bank) AS bank,
    MIN(p.delivery_cost) AS delivery_cost,
    MIN(p.goods_total) AS goods_total,
    MIN(p.custom_fee) AS custom_fee,

    -- Items
    COALESCE(
        json_agg(
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
        ) FILTER (WHERE i.id IS NOT NULL),
        '[]'
    ) AS items

FROM orders o
JOIN delivery d ON o.delivery_id = d.id
JOIN payments p ON o.payment_id = p.id
LEFT JOIN order_items oi ON o.order_uid = oi.order_uid
LEFT JOIN items i ON oi.item_id = i.id

WHERE o.order_uid = $1

GROUP BY o.order_uid;`
