package model

type DeliveryItem struct {
	ID         int    `json:"id" gorm:"primarykey"`
	DeliveryID int    `json:"delivery_id"`
	ItemName   string `json:"item_name"`
	Quantity   int    `json:"quantity"`
	Unit       string `json:"unit"`
	Weight     int    `json:"weight"`
}
