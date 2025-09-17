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
	return nil
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
