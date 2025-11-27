package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/FollG/kafka-with-go/internal/domain/models"

	_ "github.com/lib/pq"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	query := `
		INSERT INTO products (name, weight, unit, color, type, price, attributes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	attributesJSON, err := json.Marshal(product.Attributes)
	if err != nil {
		return fmt.Errorf("failed to marshal attributes: %w", err)
	}

	err = r.db.QueryRowContext(ctx, query,
		product.Name,
		product.Weight,
		product.Unit,
		product.Color,
		product.Type,
		product.Price,
		attributesJSON,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int) (*models.Product, error) {
	query := `
		SELECT id, name, weight, unit, color, type, price, attributes, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var product models.Product
	var attributesJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Weight,
		&product.Unit,
		&product.Color,
		&product.Type,
		&product.Price,
		&attributesJSON,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if err := json.Unmarshal(attributesJSON, &product.Attributes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal attributes: %w", err)
	}

	return &product, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	query := `
		UPDATE products 
		SET name = $1, weight = $2, unit = $3, color = $4, type = $5, 
			price = $6, attributes = $7, updated_at = NOW()
		WHERE id = $8
		RETURNING updated_at
	`

	attributesJSON, err := json.Marshal(product.Attributes)
	if err != nil {
		return fmt.Errorf("failed to marshal attributes: %w", err)
	}

	err = r.db.QueryRowContext(ctx, query,
		product.Name,
		product.Weight,
		product.Unit,
		product.Color,
		product.Type,
		product.Price,
		attributesJSON,
		product.ID,
	).Scan(&product.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.ErrProductNotFound
		}
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrProductNotFound
	}

	return nil
}

func (r *ProductRepository) List(ctx context.Context, filter models.ProductFilter) ([]*models.Product, error) {
	query := `
		SELECT id, name, weight, unit, color, type, price, attributes, created_at, updated_at
		FROM products
		WHERE 1=1
	`
	args := []interface{}{}
	argCounter := 1

	// Добавляем условия фильтрации
	if filter.MinPrice != nil {
		query += fmt.Sprintf(" AND price >= $%d", argCounter)
		args = append(args, *filter.MinPrice)
		argCounter++
	}

	if filter.MaxPrice != nil {
		query += fmt.Sprintf(" AND price <= $%d", argCounter)
		args = append(args, *filter.MaxPrice)
		argCounter++
	}

	if filter.Color != "" {
		query += fmt.Sprintf(" AND color = $%d", argCounter)
		args = append(args, filter.Color)
		argCounter++
	}

	if len(filter.Types) > 0 {
		placeholders := make([]string, len(filter.Types))
		for i, t := range filter.Types {
			placeholders[i] = fmt.Sprintf("$%d", argCounter)
			args = append(args, string(t))
			argCounter++
		}
		query += fmt.Sprintf(" AND type IN (%s)", strings.Join(placeholders, ","))
	}

	query += " ORDER BY created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
	args = append(args, filter.Limit, filter.Offset)
	argCounter += 2

	// Выполняем запрос
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		var attributesJSON []byte

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Weight,
			&product.Unit,
			&product.Color,
			&product.Type,
			&product.Price,
			&attributesJSON,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		if err := json.Unmarshal(attributesJSON, &product.Attributes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal attributes: %w", err)
		}

		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return products, nil
}
