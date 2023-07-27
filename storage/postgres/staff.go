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

type staffRepo struct {
	db *pgxpool.Pool
}

func NewStaffRepo(db *pgxpool.Pool) *staffRepo {
	return &staffRepo{
		db: db,
	}
}
func (r *staffRepo) Create(ctx context.Context, req *models.CreateStaff) (string, error) {

	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO staff(id, branch_id, tarif_id,type,name,balance, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.BranchID,
		req.TarifID,
		req.Type,
		req.Name,
		req.Balance,
	)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return id, nil
}

func (r *staffRepo) GetByID(ctx context.Context, req *models.StaffPrimaryKey) (*models.Staff, error) {

	var (
		query string

		id         sql.NullString
		branch_id  sql.NullString
		tarif_id   sql.NullString
		typee      sql.NullString
		name       sql.NullString
		balance    sql.NullInt64
		createdAt  sql.NullString
		updatedAt  sql.NullString
		deleted    bool
		deleted_at sql.NullString
	)

	query = `
		SELECT
			id,
			branch_id,
			tarif_id,
			type,
			name,
			balance,
			created_at,
			updated_at,
			deleted,
			deleted_at
		FROM staff
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&branch_id,
		&tarif_id,
		&typee,
		&name,
		&balance,
		&createdAt,
		&updatedAt,
		&deleted,
		&deleted_at,
	)

	if err != nil {
		return nil, err
	}

	return &models.Staff{
		Id:        id.String,
		BranchID:  branch_id.String,
		TarifID:   tarif_id.String,
		Type:      typee.String,
		Name:      name.String,
		Balance:   balance.Int64,
		CreatedAt: createdAt.String,
		UpdatedAt: updatedAt.String,
		Deleted:   deleted,
		DeletedAt: deleted_at.String,
	}, nil
}

func (r *staffRepo) GetList(ctx context.Context, req *models.StaffGetListRequest) (*models.StaffGetListResponse, error) {

	var (
		resp   = &models.StaffGetListResponse{}
		query  string
		where  = " WHERE deleted = false"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			branch_id,
			tarif_id,
			type,
			name,
			balance,
			created_at,
			updated_at,
			deleted,
			deleted_at
		FROM staff
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.Search != "" {
		where += ` AND name ILIKE '%' || '` + req.Search + `' || '%'  AND branch_id '` + req.BranchId + `'  AND  tarif_id  '` + req.TarifID + `'  AND  type  '` + req.Type + `' `
	}
	if req.From != "" && req.To != "" {
		where += ` AND balance  >  '` + req.From + `' AND  balance  <  '` + req.To + `'`
		fmt.Println("sdfsdfs")
	}

	query += where + offset + limit
	fmt.Println(query)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id         sql.NullString
			branch_id  sql.NullString
			tarif_id   sql.NullString
			typee      sql.NullString
			name       sql.NullString
			balance    sql.NullInt64
			createdAt  sql.NullString
			updatedAt  sql.NullString
			deleted    bool
			deleted_at sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&branch_id,
			&tarif_id,
			&typee,
			&name,
			&balance,
			&createdAt,
			&updatedAt,
			&deleted,
			&deleted_at,
		)

		if err != nil {
			return nil, err
		}

		resp.Staffs = append(resp.Staffs, &models.Staff{
			Id:        id.String,
			BranchID:  branch_id.String,
			TarifID:   tarif_id.String,
			Type:      typee.String,
			Name:      name.String,
			Balance:   balance.Int64,
			CreatedAt: createdAt.String,
			UpdatedAt: updatedAt.String,
			Deleted:   deleted,
			DeletedAt: deleted_at.String,
		})
	}

	return resp, nil
}

func (r *staffRepo) Update(ctx context.Context, req *models.UpdateStaff) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
		staff
		SET
			branch_id = :branch_id,
			tarif_id = :tarif_id,
			type = :type,
			name = :name,
			balance = :balance,
			updated_at = NOW()
		WHERE id = :id
	`
	params = map[string]interface{}{
		"id":        req.Id,
		"branch_id": req.BranchID,
		"tarif_id":  req.TarifID,
		"type":      req.Type,
		"name":      req.Name,
		"balance":   req.Balance,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *staffRepo) Delete(ctx context.Context, req *models.StaffPrimaryKey) error {

	_, err := r.db.Exec(ctx, "Update staff SET deleted = true, deleted_at = NOW() WHERE id = $1", req.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *staffRepo) GetTopStaff(ctx context.Context, req *models.StaffGetListRequest) (*models.StaffGetListResponse, error) {
	var (
		resp  = &models.StaffGetListResponse{}
		query string
		// where  = " WHERE deleted = false"
		// offset = " OFFSET 0"
		// limit  = " LIMIT 10"
		// having = ""
	)
	query = `
		SELECT
			s.name,
			b.name ,
			SUM(sl.price),
			s.type
		FROM 
		staff AS s 
		JOIN sales AS sl ON sl.shop_assistent_id = s.id
		JOIN branch AS b ON b.id = sl.branch_id
		WHERE sl.status = 'success'
		GROUP BY s.name,b.name,s.type
		ORDER BY SUM(sl.price) DESC
		
	`
	query2 := `
	SELECT
		s.name,
		b.name ,
		SUM(sl.price),
		s.type
	FROM 
	staff AS s 
	JOIN sales AS sl ON sl.cashier_id = s.id
	JOIN branch AS b ON b.id = sl.branch_id
	WHERE sl.status = 'success'
	GROUP BY s.name,b.name,s.type
	ORDER BY SUM(sl.price) DESC
	
`
	// if req.Offset > 0 {
	// 	offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	// }o
	// if req.Search != "" {
	// 	where += ` AND name ILIKE '%' || '` + req.Search + `' || '%' `
	// }
	// if req.From >= 0 && req.To > 0 {
	// 	having += `HAVING balance BETWEEN  '%' || '` + cast.ToString(req.From) + `' || '%' AND '%' ` + cast.ToString(req.To) + `' || '%'`
	// }
	// query += where + offset + limit + having
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var (
			name        sql.NullString
			branch_name sql.NullString
			earned_sum  sql.NullInt64
			typee       sql.NullString
		)
		err := rows.Scan(
			&name,
			&branch_name,
			&earned_sum,
			&typee,
		)
		if err != nil {
			return nil, err
		}
		resp.List = append(resp.List, &models.GetTopStaffs{
			Name:       name.String,
			BranchName: branch_name.String,
			EarnedSum:  earned_sum.Int64,
			Type:       typee.String,
		})
	}

	rows, err = r.db.Query(ctx, query2)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var (
			name        sql.NullString
			branch_name sql.NullString
			earned_sum  sql.NullInt64
			typee       sql.NullString
		)
		err := rows.Scan(
			&name,
			&branch_name,
			&earned_sum,
			&typee,
		)
		if err != nil {
			return nil, err
		}
		resp.List = append(resp.List, &models.GetTopStaffs{
			Name:       name.String,
			BranchName: branch_name.String,
			EarnedSum:  earned_sum.Int64,
			Type:       typee.String,
		})
	}
	return resp, nil
}
