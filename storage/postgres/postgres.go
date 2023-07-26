package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"app/config"
	"app/storage"
)

type store struct {
	db               *pgxpool.Pool
	branch           *branchRepo
	stafftarif       *staffTarifRepo
	staff            *staffRepo
	sales            *salesRepo
	staffTransaction *staffTransactionRepo
}

func NewConnectionPostgres(cfg *config.Config) (storage.StorageI, error) {

	connect, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%d sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresDatabase,
		cfg.PostgresPassword,
		cfg.PostgresPort,
	))

	if err != nil {
		return nil, err
	}
	connect.MaxConns = cfg.PostgresMaxConnection

	pgxpool, err := pgxpool.ConnectConfig(context.Background(), connect)
	if err != nil {
		return nil, err
	}

	return &store{
		db: pgxpool,
	}, nil
}

func (s *store) Close() {
	s.db.Close()
}

func (s *store) Branch() storage.BranchRepoI {

	if s.branch == nil {
		s.branch = NewBranchRepo(s.db)
	}

	return s.branch
}

func (s *store) StaffTarif() storage.StaffTarifRepoI {

	if s.stafftarif == nil {
		s.stafftarif = NewstaffTarifRepo(s.db)
	}

	return s.stafftarif
}

func (s *store) Staff() storage.StaffRepoI {

	if s.staff == nil {
		s.staff = NewStaffRepo(s.db)
	}

	return s.staff
}

func (s *store) Sales() storage.SalesRepoI {

	if s.sales == nil {
		s.sales = NewSalesRepo(s.db)
	}

	return s.sales
}

func (s *store) StaffTransaction() storage.StaffTransactionRepoI {

	if s.staffTransaction == nil {
		s.staffTransaction = NewStaffTransactionRepo(s.db)
	}

	return s.staffTransaction
}
