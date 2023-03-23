package psql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/delveper/mystore/app/entities"
	"github.com/delveper/mystore/app/exceptions"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

const constraintMerchantID = "products_merchant_id_fkey"

// Product represents the PostgresSQL implementation
// of the product repository and satisfies interactors.ProductRepo interface.
type Product struct{ *sql.DB }

// NewProduct creates a new instance of Product.
func NewProduct(conn *sql.DB) Product {
	return Product{conn}
}

// Insert adds a new entities.Product to the database.
func (p Product) Insert(ctx context.Context, prod entities.Product) (int, error) {
	const SQL = `INSERT INTO sales.products(merchant_id, name, description, price, status) 
				 VALUES ($1, $2, $3, $4, $5) 
				 RETURNING id;`

	row := p.QueryRowContext(ctx, SQL, prod.MerchantID, prod.Name, prod.Description, prod.Price, prod.Status)

	err := row.Scan(&prod.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, exceptions.ErrNotFound
		}

		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr); pgxErr != nil && pgxErr.ConstraintName == constraintMerchantID {
			return 0, fmt.Errorf("%w: %w", exceptions.ErrMerchantNotFound, err)
		}

		return 0, fmt.Errorf("inserting product: %w", err)
	}

	return prod.ID, nil
}

// Select retrieves an entities.Product from the database by ID.
func (p Product) Select(ctx context.Context, prod entities.Product) (*entities.Product, error) {
	const SQL = ` SELECT id, merchant_id, name, description, price, status, created_at, deleted_at 
				  FROM sales.products 
				  WHERE id=$1
				  AND deleted_at IS NULL;`

	row := p.QueryRowContext(ctx, SQL, prod.ID)

	err := row.Scan(&prod.ID, &prod.MerchantID, &prod.Name, &prod.Description, &prod.Price, &prod.Status, &prod.CreatedAt, &prod.DeletedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", exceptions.ErrNotFound, err)
		}

		return nil, fmt.Errorf("scanning product: %w", err)
	}

	return &prod, nil
}

// SelectMany retrieves all entities.Product from the database.
func (p Product) SelectMany(ctx context.Context) ([]entities.Product, error) {
	const SQL = `SELECT id, merchant_id, name, description, price, status, created_at, deleted_at 
			  	 FROM sales.products
			  	 WHERE deleted_at is NULL;` // There is a room for improvement we can add WHERE clause alongside with OFFSET and LIMIT

	rows, err := p.QueryContext(ctx, SQL)
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

	if len(prods) == 0 {
		return nil, exceptions.ErrNotFound
	}

	return prods, nil
}

// Update modifies an existing entities.Product in the database.
func (p Product) Update(ctx context.Context, prod entities.Product) error {
	SQL := `UPDATE sales.products 
			SET merchant_id=$1, name=$2, description=$3, price=$4, status=$5
			WHERE id=$6
			AND deleted_at IS NULL;`

	result, err := p.ExecContext(ctx, SQL, prod.MerchantID, prod.Name, prod.Description, prod.Price, prod.Status, prod.ID)

	if err != nil {
		return fmt.Errorf("updating product: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting affected rows while updating product: %w", err)
	}

	if n == 0 {
		return fmt.Errorf("product with id %v not found: %w", prod.ID, exceptions.ErrNotFound)
	}

	return nil
}

// Delete softly removes existing entities.Product from database.
func (p Product) Delete(ctx context.Context, prod entities.Product) (int, error) {
	const SQL = `UPDATE sales.products
				 SET deleted_at = now()
       		  	 WHERE id=$1
				 AND deleted_at IS NULL
				 RETURNING id;`

	row := p.QueryRowContext(ctx, SQL, prod.ID)

	err := row.Scan(&prod.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, exceptions.ErrNotFound
		}

		return 0, fmt.Errorf("deleting product: %w", err)
	}

	return prod.ID, nil
}
