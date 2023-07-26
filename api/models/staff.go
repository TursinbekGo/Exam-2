package models

type StaffPrimaryKey struct {
	Id string `json:"id"`
}

type CreateStaff struct {
	BranchID string `json:"branch_id"`
	TarifID  string `json:"tarif_id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Balance  int64  `json:"balance"`
}
type Staff struct {
	Id        string `json:"id"`
	BranchID  string `json:"branch_id"`
	TarifID   string `json:"tarif_id"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Balance   int64  `json:"balance"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Deleted   bool   `json:"deleted"`
	DeletedAt string `json:"deleted_at"`
}

type UpdateStaff struct {
	Id       string `json:"id"`
	BranchID string `json:"branch_id"`
	TarifID  string `json:"tarif_id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Balance  int64  `json:"balance"`
}

type StaffGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type StaffGetListResponse struct {
	Count  int      `json:"count"`
	Staffs []*Staff `json:"staffs"`
}
