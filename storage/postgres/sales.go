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

type salesRepo struct {
	db *pgxpool.Pool
}

func NewSalesRepo(db *pgxpool.Pool) *salesRepo {
	return &salesRepo{
		db: db,
	}
}

type Sales struct {
	Id              string `json:"id"`
	BranchID        string `json:"branch_id"`
	ShopAssistantID string `json:"shop_assistent_id"`
	CashierID       string `json:"cashier_id"`
	Price           string `json:"price"`
	PaymentType     int64  `json:"payment_type"`
	Status          int64  `json:"status"`
	ClientName      int64  `json:"client_name"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	Deleted         bool   `json:"deleted"`
	DeletedAt       string `json:"deleted_at"`
}

func (r *salesRepo) Create(ctx context.Context, req *models.CreateSales) (string, error) {

	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO sales(id, branch_id, shop_assistent_id,cashier_id,price,payment_type,status,client_name, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.BranchID,
		req.ShopAssistantID,
		req.CashierID,
		req.Price,
		req.PaymentType,
		req.Status,
		req.ClientName,
	)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return id, nil
}

func (r *salesRepo) GetByID(ctx context.Context, req *models.SalesPrimaryKey) (*models.Sales, error) {

	var (
		query string

		id              sql.NullString
		branch_id       sql.NullString
		shopAssistantID sql.NullString
		cashierID       sql.NullString
		price           sql.NullInt64
		paymentType     sql.NullString
		status          sql.NullString
		clientName      sql.NullString
		createdAt       sql.NullString
		updatedAt       sql.NullString
		deleted         bool
		deleted_at      sql.NullString
	)

	query = `
		SELECT
			id,
			branch_id,
			shop_assistent_id,
			cashier_id,
			price,
			payment_type,
			status,
			client_name,
			created_at,
			updated_at,
			deleted,
			deleted_at
		FROM sales
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&branch_id,
		&shopAssistantID,
		&cashierID,
		&price,
		&paymentType,
		&status,
		&clientName,
		&createdAt,
		&updatedAt,
		&deleted,
		&deleted_at,
	)

	if err != nil {
		return nil, err
	}

	return &models.Sales{
		Id:              id.String,
		BranchID:        branch_id.String,
		ShopAssistantID: shopAssistantID.String,
		CashierID:       cashierID.String,
		Price:           price.Int64,
		PaymentType:     paymentType.String,
		Status:          status.String,
		ClientName:      clientName.String,
		CreatedAt:       createdAt.String,
		UpdatedAt:       updatedAt.String,
		Deleted:         deleted,
		DeletedAt:       deleted_at.String,
	}, nil
}

func (r *salesRepo) GetList(ctx context.Context, req *models.SalesGetListRequest) (*models.SalesGetListResponse, error) {

	var (
		resp   = &models.SalesGetListResponse{}
		query  string
		where  = " WHERE deleted = false"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
		having = " "
		order  = " ORDER BY created_at DESC "
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			branch_id,
			shop_assistent_id,
			cashier_id,
			price,
			payment_type,
			status,
			client_name,
			created_at,
			updated_at,
			deleted,
			deleted_at
		FROM sales
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}
	if req.Search != "" {
		where += ` AND client_name ILIKE '%' || '` + req.Search + `' || '%'`
	}
	if len(req.From) > 0 && len(req.To) > 0 {
		having += ` HAVING balance BETWEEN	'%' || '` + req.From + `' || '%' AND '%'` + req.To + `' || '%'`
	}
	query += where + order + offset + limit + having

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id              sql.NullString
			branch_id       sql.NullString
			shopAssistantID sql.NullString
			cashierID       sql.NullString
			price           sql.NullInt64
			paymentType     sql.NullString
			status          sql.NullString
			clientName      sql.NullString
			createdAt       sql.NullString
			updatedAt       sql.NullString
			deleted         bool
			deleted_at      sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&branch_id,
			&shopAssistantID,
			&cashierID,
			&price,
			&paymentType,
			&status,
			&clientName,
			&createdAt,
			&updatedAt,
			&deleted,
			&deleted_at,
		)

		if err != nil {
			return nil, err
		}

		resp.Sales = append(resp.Sales, &models.Sales{
			Id:              id.String,
			BranchID:        branch_id.String,
			ShopAssistantID: shopAssistantID.String,
			CashierID:       cashierID.String,
			Price:           price.Int64,
			PaymentType:     paymentType.String,
			Status:          status.String,
			ClientName:      clientName.String,
			CreatedAt:       createdAt.String,
			UpdatedAt:       updatedAt.String,
			Deleted:         deleted,
			DeletedAt:       deleted_at.String,
		})
	}

	return resp, nil
}

func (r *salesRepo) Update(ctx context.Context, req *models.UpdateSales) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
		sales 
		SET
		branch_id = :branch_id,
		shop_assistent_id = :shop_assistent_id,
		cashier_id = :cashier_id,
		price = :price,
		payment_type = :payment_type,
		status = :status,
		client_name = :client_name,
		updated_at = NOW()
		WHERE id = :id
	`
	params = map[string]interface{}{
		"id":                req.Id,
		"branch_id":         req.BranchID,
		"shop_assistent_id": req.ShopAssistantID,
		"cashier_id":        req.CashierID,
		"price":             req.Price,
		"payment_type":      req.PaymentType,
		"status":            req.Status,
		"client_name":       req.ClientName,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *salesRepo) Delete(ctx context.Context, req *models.SalesPrimaryKey) error {

	_, err := r.db.Exec(ctx, "Update sales SET deleted = true, deleted_at = NOW() WHERE id = $1", req.Id)
	if err != nil {
		return err
	}
	return nil
}
