package model

type DeliveryItem struct {
	ID         int    `json:"id" gorm:"primarykey"`
	DeliveryID int    `json:"delivery_id"`
	ItemName   string `json:"item_name"`
	Quantity   string `json:"quantity"`
	Weight     int    `json:"weight"`
}
