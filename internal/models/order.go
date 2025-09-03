package models

//TODO ХЗ ВАЩЕ ПО ПОВОДУ ЭТОЙ ТАБЛИЦЫ
type Order struct {
	Order_uid          string
	Track_number       string
	Entry              string
	Delivery_id        int
	Payment_id         int
	Items              []Item
	Locale             string
	Internal_signature string
	Customer_id        string
	Delivery_service   string
	Shardkey           string
	Sm_id              string
	Date_created       string
	Off_shard          string
}
