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
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
	Search   string `json:"search"`
	From     string `json:"from"`
	To       string `json:"to"`
	BranchId string `json:"branch_id"`
	TarifID  string `json:"tarif_id"`
	Type     string `json:"type"`
}

type StaffGetListResponse struct {
	Count  int             `json:"count"`
	Staffs []*Staff        `json:"staffs"`
	List   []*GetTopStaffs `json:"list"`
}

type GetTopStaffs struct {
	Name       string `json:"name"`
	BranchName string `json:"branch_name"`
	EarnedSum  int64  `json:"earned_sum"`
	Type       string `json:"type"`
}
