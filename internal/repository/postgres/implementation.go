package postgres

import (
	"context"
	"database/sql"
	"demo-service/internal/models"
	"demo-service/internal/repository/postgres/query"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var ErrConnIsNil = errors.New("connPool is nil")

func (d *Db) Get(ctx context.Context, orderUID string) (*models.Order, error) {

	defer func(start time.Time) {
		duration := time.Since(start)
		d.log.Debug("GetOrder completed",
			"order_uid", orderUID,
			"duration_ms", duration.Milliseconds(),
		)
	}(time.Now())

	var order models.Order
	var jsonData []byte

	err := d.connPool.QueryRow(ctx, query.GetOrder, orderUID).Scan(&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OffShard,

		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,

		&order.Payment.Transaction,
		&order.Payment.RequestID,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDT,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
		&jsonData)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to scan order: %w", err)
	}
	if len(jsonData) > 0 && string(jsonData) != "[]" {
		if err := json.Unmarshal(jsonData, &order.Items); err != nil {
			return nil, fmt.Errorf("failed to unmarshal items: %w", err)
		}
	}
	return &order, nil
}

func (d *Db) Insert(ctx context.Context, order *models.Order) error {
	defer func(start time.Time) {
		duration := time.Since(start)
		d.log.Debug("InsertOrder completed",
			"duration_ms", duration.Milliseconds(),
		)
	}(time.Now())
	tx, err := d.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // откатим, если не вызовем Commit

	var deliveryID int
	err = tx.QueryRow(ctx, `
		INSERT INTO delivery (name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	).Scan(&deliveryID)
	if err != nil {
		return err
	}

	// 2. Вставляем payment → получаем payment_id
	var paymentID int
	err = tx.QueryRow(ctx, `
		INSERT INTO payments (
			transaction, request_id, currency, provider, amount,
			payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDT,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	).Scan(&paymentID)
	if err != nil {
		return err
	}

	// 3. Вставляем заказ, связывая с delivery и payment
	_, err = tx.Exec(ctx, `
		INSERT INTO orders (
			order_uid, track_number, entry, delivery_id, payment_id,
			locale, internal_signature, customer_id, delivery_service,
			shardkey, sm_id, date_created, off_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (order_uid) DO NOTHING  -- ← идемпотентность
	`,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		deliveryID,
		paymentID,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OffShard,
	)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		// Вставляем item → получаем item_id
		var itemID int
		err = tx.QueryRow(ctx, `
			INSERT INTO items (
				chrt_id, track_number, price, rid, name, sale, size,
				total_price, nm_id, brand, status
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id
		`,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		).Scan(&itemID)
		if err != nil {
			return err
		}

		// Связываем заказ с item
		_, err = tx.Exec(ctx, `
			INSERT INTO order_items (order_uid, item_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`,
			order.OrderUID,
			itemID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (d *Db) GetAll(ctx context.Context, rowsCount uint64) ([]*models.Order, error) {

	defer func(start time.Time) {
		duration := time.Since(start)
		d.log.Debug("GetAll completed",
			"duration_ms", duration.Milliseconds(),
		)
	}(time.Now())

	orders := []*models.Order{}

	rows, err := d.connPool.Query(ctx, query.GetOrderAll, rowsCount)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		order := models.Order{}
		jsonData := []byte{}
		err := rows.Scan(&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SmID,
			&order.DateCreated,
			&order.OffShard,

			&order.Delivery.Name,
			&order.Delivery.Phone,
			&order.Delivery.Zip,
			&order.Delivery.City,
			&order.Delivery.Address,
			&order.Delivery.Region,
			&order.Delivery.Email,

			&order.Payment.Transaction,
			&order.Payment.RequestID,
			&order.Payment.Currency,
			&order.Payment.Provider,
			&order.Payment.Amount,
			&order.Payment.PaymentDT,
			&order.Payment.Bank,
			&order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal,
			&order.Payment.CustomFee,
			&jsonData)
		if err != nil {
			return nil, err
		}

		if len(jsonData) > 0 && string(jsonData) != "[]" {
			if err := json.Unmarshal(jsonData, &order.Items); err != nil {
				return nil, fmt.Errorf("failed to unmarshal items: %w", err)
			}
		}

		orders = append(orders, &order)
	}

	return orders, nil
}

func (d *Db) Close() error {
	if d.connPool != nil {
		d.connPool.Close()
		return nil
	}
	return ErrConnIsNil
}
