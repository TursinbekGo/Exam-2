package models

type StaffTarifPrimaryKey struct {
	Id string `json:"id"`
}

type CreateStaffTarif struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	AmountForCash int64  `json:"amountforcash"`
	AmountForCard int64  `json:"amountforcard"`
}
type StaffTarif struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	AmountForCash int64  `json:"amountforcash"`
	AmountForCard int64  `json:"amountforcard"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	Deleted       bool   `json:"deleted"`
	DeletedAt     string `json:"deleted_at"`
}

type UpdateStaffTarif struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	AmountForCash int64  `json:"amountforcash"`
	AmountForCard int64  `json:"amountforcard"`
}

type StaffTarifGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type StaffTarifGetListResponse struct {
	Count       int           `json:"count"`
	StaffTarifs []*StaffTarif `json:"staff_tarifs"`
}
