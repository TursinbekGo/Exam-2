package models

type SalesPrimaryKey struct {
	Id string `json:"id"`
}

type CreateSales struct {
	BranchID        string `json:"branch_id"`
	ShopAssistantID string `json:"shop_assistent_id"`
	CashierID       string `json:"cashier_id"`
	Price           int64  `json:"price"`
	PaymentType     string `json:"payment_type"`
	Status          string `json:"status"`
	ClientName      string `json:"client_name"`
}

type Sales struct {
	Id              string `json:"id"`
	BranchID        string `json:"branch_id"`
	ShopAssistantID string `json:"shop_assistent_id"`
	CashierID       string `json:"cashier_id"`
	Price           int64  `json:"price"`
	PaymentType     string `json:"payment_type"`
	Status          string `json:"status"`
	ClientName      string `json:"client_name"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	Deleted         bool   `json:"deleted"`
	DeletedAt       string `json:"deleted_at"`
}

type UpdateSales struct {
	Id              string `json:"id"`
	BranchID        string `json:"branch_id"`
	ShopAssistantID string `json:"shop_assistent_id"`
	CashierID       string `json:"cashier_id"`
	Price           int64  `json:"price"`
	PaymentType     string `json:"payment_type"`
	Status          string `json:"status"`
	ClientName      string `json:"client_name"`
}

type SalesGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
	From   string `json:"from"`
	To     string `json:"to"`
}

type SalesGetListResponse struct {
	Count int      `json:"count"`
	Sales []*Sales `json:"sales"`
}
