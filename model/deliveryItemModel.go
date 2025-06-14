package model

type DeliveryItem struct {
	ID         int    `json:"id" gorm:"primarykey"`
	DeliveryID int    `json:"delibery_id"`
	ItemName   string `json:"item_name"`
	Quantity   string `json:"quantity"`
	Weight     string `json:"weight"`
}
