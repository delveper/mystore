package psql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/delveper/mystore/app/entities"
	"github.com/pkg/errors"
)

type Product struct{ *sql.DB }

func (r *Product) Insert(ctx context.Context, prod entities.Product) error {
	const SQL = `INSERT INTO products(merchant_id, name, description, price, status, created_at) 
				 VALUES ($1, $2, $3, $4, $5, $6) 
				 RETURNING id;`

	row := r.DB.QueryRowContext(ctx, SQL, prod.MerchantID, prod.Name, prod.Description, prod.Price, prod.Status, prod.CreatedAt)

	err := row.Scan(&prod.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert product")
	}

	return nil
}

func (r *Product) Select(ctx context.Context, prod entities.Product) (*entities.Product, error) {
	const SQL = ` SELECT id, merchant_id, name, description, price, status, created_at, deleted_at 
				  FROM products 
				  WHERE id = $1;`

	row := r.DB.QueryRowContext(ctx, SQL, prod.ID)

	err := row.Scan(&prod.ID, &prod.MerchantID, &prod.Name, &prod.Description, &prod.Price, &prod.Status, &prod.CreatedAt, &prod.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select product")
	}

	return &prod, nil
}

func (r *Product) SelectMany(ctx context.Context) ([]entities.Product, error) {
	const SQL = `SELECT id, merchant_id, name, description, price, status, created_at, deleted_at 
			  	 FROM products;`

	rows, err := r.DB.QueryContext(ctx, SQL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select products")
	}

	defer rows.Close()

	var prods []entities.Product

	for rows.Next() {
		var prod entities.Product

		err = rows.Scan(&prod.ID, &prod.MerchantID, &prod.Name, &prod.Description, &prod.Price, &prod.Status, &prod.CreatedAt, &prod.DeletedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan product")
		}

		prods = append(prods, prod)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select products")
	}

	return prods, nil
}

func (r *Product) Update(ctx context.Context, prod entities.Product) error {
	SQL := `UPDATE products 
			SET merchant_id=$1, name=$2, description=$3, price=$4, status=$5, created_at=$6, deleted_at=$7 
			WHERE id=$8;`

	result, err := r.DB.ExecContext(ctx, SQL, prod.MerchantID, prod.Name, prod.Description, prod.Price, prod.Status, prod.CreatedAt, prod.DeletedAt, prod.ID)
	if err != nil {
		return errors.Wrap(err, "failed to update product")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %d not found", prod.ID)
	}

	return nil
}

func (r *Product) Delete(ctx context.Context, prod entities.Product) error {
	const SQL = `DELETE FROM products 
       		  	 WHERE id=$1;`

	result, err := r.DB.ExecContext(ctx, SQL, prod.ID)
	if err != nil {
		return errors.Wrap(err, "failed to delete product")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %d not found", prod.ID)
	}

	return nil
}
