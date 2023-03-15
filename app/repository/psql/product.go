package psql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/delveper/mystore/app/entities"
	"github.com/delveper/mystore/app/exceptions"
	"github.com/pkg/errors"
)

type Product struct{ *sql.DB }

func (r *Product) Insert(ctx context.Context, prod entities.Product) (int, error) {
	const SQL = `INSERT INTO products(merchant_id, name, description, price, status) 
				 VALUES ($1, $2, $3, $4, $5) 
				 RETURNING id;`

	row := r.DB.QueryRowContext(ctx, SQL, prod.MerchantID, prod.Name, prod.Description, prod.Price, prod.Status)

	err := row.Scan(&prod.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, exceptions.ErrRecordNotFound
		}

		return 0, fmt.Errorf("inserting product: %w", err)
	}

	return prod.ID, nil
}

func (r *Product) Select(ctx context.Context, prod entities.Product) (*entities.Product, error) {
	const SQL = ` SELECT id, merchant_id, name, description, price, status, created_at, deleted_at 
				  FROM products 
				  WHERE id=$1
				  AND deleted_at IS NULL;`

	row := r.DB.QueryRowContext(ctx, SQL, prod.ID)

	err := row.Scan(&prod.ID, &prod.MerchantID, &prod.Name, &prod.Description, &prod.Price, &prod.Status, &prod.CreatedAt, &prod.DeletedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", exceptions.ErrRecordNotFound, err)
		}

		return nil, fmt.Errorf("scanning product: %w", err)
	}

	return &prod, nil
}

func (r *Product) SelectMany(ctx context.Context) ([]entities.Product, error) {
	const SQL = `SELECT id, merchant_id, name, description, price, status, created_at, deleted_at 
			  	 FROM products` // There is a room for improvement we can add WHERE clause alongside with OFFSET and LIMIT

	rows, err := r.DB.QueryContext(ctx, SQL)
	if err != nil {
		return nil, fmt.Errorf("executing query by selecting products: %w", err)
	}

	defer rows.Close()

	var prods []entities.Product

	for rows.Next() {
		var prod entities.Product

		err = rows.Scan(&prod.ID, &prod.MerchantID, &prod.Name, &prod.Description, &prod.Price, &prod.Status, &prod.CreatedAt, &prod.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("scanning product: %w", err)
		}

		prods = append(prods, prod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("during iteration products: %w", err)
	}

	return prods, nil
}

func (r *Product) Update(ctx context.Context, prod entities.Product) error {
	SQL := `UPDATE products 
			SET merchant_id=$1, name=$2, description=$3, price=$4, status=$5
			WHERE id=$6
			AND deleted_at IS NULL;`

	result, err := r.DB.ExecContext(ctx, SQL, prod.MerchantID, prod.Name, prod.Description, prod.Price, prod.Status, prod.ID)

	if err != nil {
		return fmt.Errorf("updating product: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting affected rows while updating product: %w", err)
	}

	if n == 0 {
		return fmt.Errorf("product with id %v not found: %w", prod.ID, exceptions.ErrRecordNotFound)
	}

	return nil
}

func (r *Product) Delete(ctx context.Context, prod entities.Product) error {
	const SQL = `UPDATE products
				 SET deleted_at = now()
       		  	 WHERE id=$1
				 AND deleted_at IS NULL
				 RETURNING id;`

	row := r.DB.QueryRowContext(ctx, SQL, prod.ID)

	err := row.Scan(&prod.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return exceptions.ErrRecordNotFound
		}

		return fmt.Errorf("deleting product: %w", err)
	}

	return nil
}
