package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

// Connection descriptor for PostgreSQL
type DBRepository struct {
	connection *pgx.Conn
}

// Constructor for PostgreSQL connection pool
func NewDBRepository(pgConnection *pgx.Conn) (*DBRepository, error) {
	sql := `create table if not exists urls_base (
		short_url text not null primary key,
		original_url text,
		user_id      text,
		deleted_at bool default false		
	); 
	create unique index if not exists original_url_constrain on urls_base(original_url);`

	_, err := pgConnection.Exec(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	return &DBRepository{connection: pgConnection}, nil
}

// Insert new short URL in DB
func (r *DBRepository) Add(key, value, userID string) error {
	sql := `insert into urls_base (short_url, original_url, user_id) values ($1, $2, $3)`
	_, err := r.connection.Exec(context.Background(), sql, key, value, userID)
	if err != nil {
		return err
	}
	return nil
}

// Select origin URL by short URL from DB
func (r *DBRepository) Get(key, userID string) (string, error) {
	sql := `select original_url, deleted_at from urls_base where short_url = $1`
	row := r.connection.QueryRow(context.Background(), sql, key)
	var value string
	var deletedAt bool
	err := row.Scan(&value, &deletedAt)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if deletedAt {
		return "", ErrURLDeleted
	}

	return value, nil
}

// Select URLs by user ID from DB
func (r *DBRepository) GetUserURLs(userID string) (map[string]string, error) {
	result := make(map[string]string)
	sql := `select short_url, original_url from urls_base where user_id = $1 and deleted_at = false`
	rows, err := r.connection.Query(context.Background(), sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var shortURL, originalURL string
		err = rows.Scan(&shortURL, &originalURL)
		if err != nil {
			return nil, err
		}
		result[shortURL] = originalURL
	}
	return result, nil
}

// Delete URL for user ID from DB
func (r *DBRepository) Delete(keys []string, userID string) error {
	sql := `update urls_base set deleted_at = true where short_url = any($1) and user_id = $2`
	_, err := r.connection.Exec(context.Background(), sql, keys, userID)
	if err != nil {
		return err
	}
	return nil
}

// Returs stats from DB
func (r *DBRepository) GetStats() (int, int, error) {
	sql := `select urls, users from (select count(*) from urls_base where deleted_at = false), (select count(distinct user_id) from urls_base)`
	row := r.connection.QueryRow(context.Background(), sql)
	var urls int
	var users int
	err := row.Scan(&urls, &users)
	if err != nil {
		return 0, 0, err
	}

	return urls, users, nil
}
