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

type staffTarifRepo struct {
	db *pgxpool.Pool
}

func NewstaffTarifRepo(db *pgxpool.Pool) *staffTarifRepo {
	return &staffTarifRepo{
		db: db,
	}
}

func (r *staffTarifRepo) Create(ctx context.Context, req *models.CreateStaffTarif) (string, error) {

	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO staff_tarif(id, name, type,amountforcash,amountforcard, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.Name,
		req.Type,
		req.AmountForCash,
		req.AmountForCard,
	)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return id, nil
}

func (r *staffTarifRepo) GetByID(ctx context.Context, req *models.StaffTarifPrimaryKey) (*models.StaffTarif, error) {

	var (
		query string

		id            sql.NullString
		name          sql.NullString
		typee         sql.NullString
		amountForCash sql.NullInt64
		amountForCard sql.NullInt64
		createdAt     sql.NullString
		updatedAt     sql.NullString
		deleted       bool
		deleted_at    sql.NullString
	)

	query = `
		SELECT
			id,
			name,
			type,
			amountforcash,
			amountforcard,
			created_at,
			updated_at,
			deleted,
			deleted_at
		FROM staff_tarif
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&name,
		&typee,
		&amountForCash,
		&amountForCard,
		&createdAt,
		&updatedAt,
		&deleted,
		&deleted_at,
	)

	if err != nil {
		return nil, err
	}

	return &models.StaffTarif{
		Id:            id.String,
		Name:          name.String,
		Type:          typee.String,
		AmountForCash: amountForCash.Int64,
		AmountForCard: amountForCard.Int64,
		CreatedAt:     createdAt.String,
		UpdatedAt:     updatedAt.String,
		Deleted:       deleted,
		DeletedAt:     deleted_at.String,
	}, nil
}

func (r *staffTarifRepo) GetList(ctx context.Context, req *models.StaffTarifGetListRequest) (*models.StaffTarifGetListResponse, error) {

	var (
		resp   = &models.StaffTarifGetListResponse{}
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
			type,
			amountforcash,
			amountforcard,
			created_at,
			updated_at,
			deleted,
			deleted_at
		FROM staff_tarif
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.Search != "" {
		where += ` AND name ILIKE '%' || '` + req.Search + `' || '%'  `
	}

	query += where + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id            sql.NullString
			name          sql.NullString
			typee         sql.NullString
			amountForCash sql.NullInt64
			amountForCard sql.NullInt64
			createdAt     sql.NullString
			updatedAt     sql.NullString
			deleted       bool
			deleted_at    sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&name,
			&typee,
			&amountForCash,
			&amountForCard,
			&createdAt,
			&updatedAt,
			&deleted,
			&deleted_at,
		)

		if err != nil {
			return nil, err
		}

		resp.StaffTarifs = append(resp.StaffTarifs, &models.StaffTarif{
			Id:            id.String,
			Name:          name.String,
			Type:          typee.String,
			AmountForCash: amountForCash.Int64,
			AmountForCard: amountForCard.Int64,
			CreatedAt:     createdAt.String,
			UpdatedAt:     updatedAt.String,
			Deleted:       deleted,
			DeletedAt:     deleted_at.String,
		})
	}

	return resp, nil
}

func (r *staffTarifRepo) Update(ctx context.Context, req *models.UpdateStaffTarif) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			staff_tarif
		SET
			name = :name,
			type = :type,
			amountforcash = :amountforcash,
			amountforcard = :amountforcard,
			updated_at = NOW()
		WHERE id = :id
	`
	params = map[string]interface{}{
		"id":            req.Id,
		"name":          req.Name,
		"type":          req.Type,
		"amountforcash": req.AmountForCash,
		"amountforcard": req.AmountForCard,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *staffTarifRepo) Delete(ctx context.Context, req *models.StaffTarifPrimaryKey) error {

	_, err := r.db.Exec(ctx, "Update branch SET deleted = true, deleted_at = NOW() WHERE id = $1", req.Id)
	if err != nil {
		return err
	}
	return nil
}
