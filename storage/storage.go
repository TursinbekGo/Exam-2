package storage

import (
	"app/api/models"
	"context"
)

type StorageI interface {
	Close()
	Branch() BranchRepoI
	StaffTarif() StaffTarifRepoI
	Staff() StaffRepoI
	Sales() SalesRepoI
	StaffTransaction() StaffTransactionRepoI
}

type BranchRepoI interface {
	Create(context.Context, *models.CreateBranch) (string, error)
	GetByID(context.Context, *models.BranchPrimaryKey) (*models.Branch, error)
	GetList(context.Context, *models.BranchGetListRequest) (*models.BranchGetListResponse, error)
	Update(context.Context, *models.UpdateBranch) (int64, error)
	Delete(context.Context, *models.BranchPrimaryKey) error
}

type StaffTarifRepoI interface {
	Create(context.Context, *models.CreateStaffTarif) (string, error)
	GetByID(context.Context, *models.StaffTarifPrimaryKey) (*models.StaffTarif, error)
	GetList(context.Context, *models.StaffTarifGetListRequest) (*models.StaffTarifGetListResponse, error)
	Update(context.Context, *models.UpdateStaffTarif) (int64, error)
	Delete(context.Context, *models.StaffTarifPrimaryKey) error
}

type StaffRepoI interface {
	Create(context.Context, *models.CreateStaff) (string, error)
	GetByID(context.Context, *models.StaffPrimaryKey) (*models.Staff, error)
	GetList(context.Context, *models.StaffGetListRequest) (*models.StaffGetListResponse, error)
	Update(context.Context, *models.UpdateStaff) (int64, error)
	Delete(context.Context, *models.StaffPrimaryKey) error
}

type SalesRepoI interface {
	Create(context.Context, *models.CreateSales) (string, error)
	GetByID(context.Context, *models.SalesPrimaryKey) (*models.Sales, error)
	GetList(context.Context, *models.SalesGetListRequest) (*models.SalesGetListResponse, error)
	Update(context.Context, *models.UpdateSales) (int64, error)
	Delete(context.Context, *models.SalesPrimaryKey) error
}

type StaffTransactionRepoI interface {
	Create(context.Context, *models.CreateStaffTransaction) (string, error)
	GetByID(context.Context, *models.StaffTransactionPrimaryKey) (*models.StaffTransaction, error)
	GetList(context.Context, *models.StaffTransactionGetListRequest) (*models.StaffTransactionGetListResponse, error)
	Update(context.Context, *models.UpdateStaffTransaction) (int64, error)
	Delete(context.Context, *models.StaffTransactionPrimaryKey) error
}
