package repository

import (
	"database/sql"
)

type IdempotencyRepo struct {
	db *sql.DB
}

func NewIdempotencyRepo(db *sql.DB) *IdempotencyRepo {
	return &IdempotencyRepo{db: db}
}

func (r *IdempotencyRepo) IsIdempotencyKeyExists(key string) (bool, error) {
	query, err := r.db.Prepare("SELECT EXISTS (SELECT * FROM idempotency WHERE idempotency_key=$1)")
	if err != nil {
		return false, err
	}
	defer query.Close()

	row, err := query.Query(key)
	var isExists bool

	if err != nil {
		return false, err
	}

	for row.Next() {
		err = row.Scan(&isExists)
		if err != nil {
			return false, err
		}
	}

	return isExists, nil
}

func (r *IdempotencyRepo) CreateIdempotencyKey(key string, expression_id int) error {
	query, err := r.db.Prepare("INSERT INTO idempotency (idempotency_key, expression_id) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer query.Close()

	query.Exec(key, expression_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *IdempotencyRepo) GetExpressionId(key string) (int, error) {
	var id int = 0

	query, err := r.db.Prepare("SELECT expression_id from idempotency WHERE idempotency_key=$1")
	if err != nil {
		return id, err
	}
	defer query.Close()

	row, err := query.Query(key)
	if err != nil {
		return id, err
	}

	for row.Next() {
		err = row.Scan(&id)
		if err != nil {
			return id, err
		}
	}
	return id, nil
}
