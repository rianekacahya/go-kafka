package entity

type OrderHeaders struct {
	ID            int32      `json:"id,omitempty"`
	CustomerEmail NullString `json:"customer_email,omitempty"`
	PurchaseCode  NullString `json:"purchase_code,omitempty"`
	PurchaseDate  NullTime   `json:"purchase_date,omitempty"`
}

type OrderItems struct {
	ID           int32      `json:"id,omitempty"`
	PurchaseCode NullString `json:"purchase_code,omitempty"`
	SKU          NullString `json:"sku,omitempty"`
	Quantity     NullInt32  `json:"quantity,omitempty"`
}

type Orders struct {
	*OrderHeaders
	OrderItems []OrderItems `json:"order_items,omitempty"`
}
