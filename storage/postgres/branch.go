package postgres

import (
	"app/api/models"
	"app/pkg/helper"
	"context"
	"database/sql"
	"fmt"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type branchRepo struct {
	db *pgxpool.Pool
}

func NewBranchRepo(db *pgxpool.Pool) *branchRepo {
	return &branchRepo{
		db: db,
	}
}

func (r *branchRepo) Create(ctx context.Context, req *models.CreateBranch) (string, error) {

	var (
		id    = uuid.New().String()
		query string
	)
	query = `
		INSERT INTO branch(id, name, address, updated_at)
		VALUES ($1, $2, $3,NOW())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.Name,
		req.Address,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *branchRepo) GetByID(ctx context.Context, req *models.BranchPrimaryKey) (*models.Branch, error) {

	var (
		query string

		id         sql.NullString
		name       sql.NullString
		address    sql.NullString
		createdAt  sql.NullString
		updatedAt  sql.NullString
		deleted    bool
		deleted_at sql.NullString
	)

	query = `
		SELECT
			id,
			name,
			address,
			created_at,
			updated_at,
			deleted,
			deleted_at
		FROM branch
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&name,
		&address,
		&createdAt,
		&updatedAt,
		&deleted,
		&deleted_at,
	)

	if err != nil {
		return nil, err
	}

	return &models.Branch{
		Id:        id.String,
		Name:      name.String,
		Address:   address.String,
		CreatedAt: createdAt.String,
		UpdatedAt: updatedAt.String,
		Deleted:   deleted,
		DeletedAt: deleted_at.String,
	}, nil
}

func (r *branchRepo) GetList(ctx context.Context, req *models.BranchGetListRequest) (*models.BranchGetListResponse, error) {

	var (
		resp   = &models.BranchGetListResponse{}
		query  string
		where  = " WHERE deleted = false"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			name,
			address,
			created_at,
			updated_at,
			deleted,
			deleted_at
		FROM branch
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.Search != "" {
		where += ` AND name ILIKE '%' || '` + req.Search + `' || '%'  OR address ILIKE '%' || '` + req.Search + `' || '%' `
	}

	query += where + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id         sql.NullString
			name       sql.NullString
			address    sql.NullString
			createdAt  sql.NullString
			updatedAt  sql.NullString
			deleted    bool
			deleted_at sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&name,
			&address,
			&createdAt,
			&updatedAt,
			&deleted,
			&deleted_at,
		)

		if err != nil {
			return nil, err
		}

		resp.Branches = append(resp.Branches, &models.Branch{
			Id:        id.String,
			Name:      name.String,
			Address:   address.String,
			CreatedAt: createdAt.String,
			UpdatedAt: updatedAt.String,
			Deleted:   deleted,
			DeletedAt: deleted_at.String,
		})
	}

	return resp, nil
}

func (r *branchRepo) Update(ctx context.Context, req *models.UpdateBranch) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
		branch
		SET
			name = :name,
			address = :address,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":      req.Id,
		"name":    req.Name,
		"address": req.Address,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *branchRepo) Delete(ctx context.Context, req *models.BranchPrimaryKey) error {

	_, err := r.db.Exec(ctx, "Update branch SET deleted = true, deleted_at = NOW() WHERE id = $1", req.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *branchRepo) GetTopBranch(ctx context.Context, req *models.BranchGetListRequest) (*models.BranchGetListResponse, error) {

	var (
		resp  = &models.BranchGetListResponse{}
		query string
		// where  = " WHERE deleted = false"
		// offset = " OFFSET 0"
		// limit  = " LIMIT 10"
	)
	query = `
		SELECT
			date(sl.created_at),
			b.address,
			SUM(sl.price)
		FROM 
		branch AS b
		JOIN sales AS sl ON sl.branch_id = b.id
		GROUP BY date(sl.created_at),b.address
		ORDER BY SUM(sl.price) DESC

	`
	// if req.Offset > 0 {
	// 	offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	// }

	// if req.Limit > 0 {
	// 	limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	// }

	// if req.Search != "" {
	// 	where += ` AND name ILIKE '%' || '` + req.Search + `' || '%'  OR address ILIKE '%' || '` + req.Search + `' || '%' `
	// }

	// query += where + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			day     sql.NullString
			address sql.NullString
			count   sql.NullInt64
		)

		err := rows.Scan(
			&day,
			&address,
			&count,
		)

		if err != nil {
			return nil, err
		}

		resp.List = append(resp.List, &models.TopBranch{
			Day:     day.String,
			Address: address.String,
			Count:   count.Int64,
		})
	}
	return resp, nil
}
