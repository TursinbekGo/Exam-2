package models

type StaffTransactionPrimaryKey struct {
	Id string `json:"id"`
}

type CreateStaffTransaction struct {
	SalesID    string `json:"sales_id"`
	Type       string `json:"type"`
	SourceType string `json:"source_type"`
	Text       string `json:"text"`
	Amount     int64  `json:"amount"`
	StaffID    string `json:"staff_id"`
}

type StaffTransaction struct {
	Id         string `json:"id"`
	SalesID    string `json:"sales_id"`
	Type       string `json:"type"`
	SourceType string `json:"source_type"`
	Text       string `json:"text"`
	Amount     int64  `json:"amount"`
	StaffID    string `json:"staff_id"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Deleted    bool   `json:"deleted"`
	DeletedAt  string `json:"deleted_at"`
}

type UpdateStaffTransaction struct {
	Id         string `json:"id"`
	SalesID    string `json:"sales_id"`
	Type       string `json:"type"`
	SourceType string `json:"source_type"`
	Text       string `json:"text"`
	Amount     int64  `json:"amount"`
	StaffID    string `json:"staff_id"`
}

type StaffTransactionGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type StaffTransactionGetListResponse struct {
	Count             int                 `json:"count"`
	StaffTransactions []*StaffTransaction `json:"staffTransactions"`
}
