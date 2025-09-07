package postgres

import (
	"context"
	"database/sql"
	"demo-service/internal/models"
	"demo-service/internal/repository/postgres/query"
	"encoding/json"
	"errors"
	"fmt"
)

func (dataBase *Db) Get(orderUID string) (*models.Order, error) {
	var order models.Order
	var jsonData []byte
	ctx := context.Background()

	err := dataBase.connPool.QueryRow(ctx, query.GetOrder, orderUID).Scan(&order.OrderUID,
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
			return nil, fmt.Errorf("order not found: %s", orderUID)
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

func (dataBase *Db) Insert(order *models.Order) error {
	return nil
}
