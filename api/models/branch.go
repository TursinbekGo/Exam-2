package models

type BranchPrimaryKey struct {
	Id string `json:"id"`
}

type CreateBranch struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}
type Branch struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Deleted   bool   `json:"deleted"`
	DeletedAt string `json:"deleted_at"`
}

type UpdateBranch struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type BranchGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type BranchGetListResponse struct {
	Count    int       `json:"count"`
	Branches []*Branch `json:"branches"`
}
