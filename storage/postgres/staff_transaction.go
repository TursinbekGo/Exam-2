package postgres

import (
	"app/api/models"
	"app/pkg/helper"
	"context"
	"database/sql"
	"fmt"
	"log"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type staffTransactionRepo struct {
	db *pgxpool.Pool
}

func NewStaffTransactionRepo(db *pgxpool.Pool) *staffTransactionRepo {
	return &staffTransactionRepo{
		db: db,
	}
}

func (r *staffTransactionRepo) Create(ctx context.Context, req *models.CreateStaffTransaction) (string, error) {

	trx, err := r.db.Begin(ctx)
	if err != nil {
		return "", nil
	}

	defer func() {
		if err != nil {
			trx.Rollback(ctx)
		} else {
			trx.Commit(ctx)
		}
	}()

	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO staff_transaction(id, sales_id, type,source_type,text,amount,staff_id, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
	`

	_, err = trx.Exec(ctx, query,
		id,
		req.SalesID,
		req.Type,
		req.SourceType,
		req.Text,
		req.Amount,
		req.StaffID,
	)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return id, nil
}

func (r *staffTransactionRepo) GetByID(ctx context.Context, req *models.StaffTransactionPrimaryKey) (*models.StaffTransaction, error) {

	var (
		query string

		id         sql.NullString
		salesID    sql.NullString
		typee      sql.NullString
		sourceType sql.NullString
		text       sql.NullString
		amount     sql.NullInt64
		staffID    sql.NullString
		createdAt  sql.NullString
		updatedAt  sql.NullString
		deleted    bool
		deleted_at sql.NullString
	)

	query = `
		SELECT
			id,
			sales_id,
			type,
			source_type,
			text,
			amount,
			staff_id,
			created_at,
			updated_at,
			deleted,
			deleted_at
		FROM staff_transaction
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&salesID,
		&typee,
		&sourceType,
		&text,
		&amount,
		&staffID,
		&createdAt,
		&updatedAt,
		&deleted,
		&deleted_at,
	)

	if err != nil {
		return nil, err
	}

	return &models.StaffTransaction{
		Id:         id.String,
		SalesID:    salesID.String,
		Type:       typee.String,
		SourceType: sourceType.String,
		Text:       text.String,
		Amount:     amount.Int64,
		StaffID:    staffID.String,
		CreatedAt:  createdAt.String,
		UpdatedAt:  updatedAt.String,
		Deleted:    deleted,
		DeletedAt:  deleted_at.String,
	}, nil
}

func (r *staffTransactionRepo) GetList(ctx context.Context, req *models.StaffTransactionGetListRequest) (*models.StaffTransactionGetListResponse, error) {

	var (
		resp   = &models.StaffTransactionGetListResponse{}
		query  string
		where  = " WHERE deleted = false"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			sales_id,
			type,
			source_type,
			text,
			amount,
			staff_id,
			created_at,
			updated_at,
			deleted,
			deleted_at
		FROM staff_transaction
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.Search != "" {
		where += ` AND name ILIKE '%' || '` + req.Search + `' || '%'`
	}

	query += where + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id         sql.NullString
			salesID    sql.NullString
			typee      sql.NullString
			sourceType sql.NullString
			text       sql.NullString
			amount     sql.NullInt64
			staffID    sql.NullString
			createdAt  sql.NullString
			updatedAt  sql.NullString
			deleted    bool
			deleted_at sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&salesID,
			&typee,
			&sourceType,
			&text,
			&amount,
			&staffID,
			&createdAt,
			&updatedAt,
			&deleted,
			&deleted_at,
		)

		if err != nil {
			return nil, err
		}

		resp.StaffTransactions = append(resp.StaffTransactions, &models.StaffTransaction{
			Id:         id.String,
			SalesID:    salesID.String,
			Type:       typee.String,
			SourceType: sourceType.String,
			Text:       text.String,
			Amount:     amount.Int64,
			StaffID:    staffID.String,
			CreatedAt:  createdAt.String,
			UpdatedAt:  updatedAt.String,
			Deleted:    deleted,
			DeletedAt:  deleted_at.String,
		})
	}

	return resp, nil
}

func (r *staffTransactionRepo) Update(ctx context.Context, req *models.UpdateStaffTransaction) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)
	query = `
		UPDATE
		staff_transaction 
		SET
		sales_id = :sales_id,
		type = :type,
		source_type = :source_type,
		text = :text,
		amount = :amount,
		staff_id = :staff_id,
		updated_at = NOW()
		WHERE id = :id
	`
	params = map[string]interface{}{
		"id":          req.Id,
		"sales_id":    req.SalesID,
		"type":        req.Type,
		"source_type": req.SourceType,
		"text":        req.Text,
		"amount":      req.Amount,
		"staff_id":    req.StaffID,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *staffTransactionRepo) Delete(ctx context.Context, req *models.StaffTransactionPrimaryKey) error {

	_, err := r.db.Exec(ctx, "Update staff_transaction SET deleted = true, deleted_at = NOW() WHERE id = $1", req.Id)
	if err != nil {
		return err
	}
	return nil
}
