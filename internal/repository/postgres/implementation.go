package postgres

import "demo-service/internal/models"

// const queryGetOrder = `SELECT 
// 			track_number,Entry,Delivery_id,Payment_id    
// 	Items              []Item
// 	Locale             string
// 	Internal_signature string
// 	Customer_id        string
// 	Delivery_service   string
// 	Shardkey           string
// 	Sm_id              string
// 	Date_created       string
// 	Off_shard          string`

func (dataBase *Db) Get(id string) (*models.Order, error) {
	return &models.Order{}, nil
}

func (dataBase *Db) Insert(order *models.Order) error {
	return nil
}
